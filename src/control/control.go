package control

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const POSTGRES_IMAGE = "postgres:16.3-alpine3.20"
const PREFIX = "netpart-"

var ENVS = [3]string{
	"POSTGRES_USER=" + POSTGRES_USER,
	"POSTGRES_PASSWORD=" + POSTGRES_PASSWORD,
	"POSTGRES_DB=" + POSTGRES_DB,
}

type Instance struct {
	Name        string
	ContainerID string
	NetworkID   string
	Port        string
}

type ControlPlane struct {
	cli     *client.Client
	pulled  bool
	servers []Instance
}

func MakeControlPlane(ctx context.Context, ops ...client.Opt) (*ControlPlane, error) {
	cli, err := client.NewClientWithOpts(ops...)
	if err != nil {
		return nil, err
	}

	for {
		_, err := cli.Ping(ctx)
		if err != nil {
			fmt.Println(err)
			time.Sleep(500 * time.Millisecond)
		} else {
			fmt.Println("docker daemon connected!")
			break
		}
	}

	return &ControlPlane{
		cli: cli,
	}, nil
}

func (c *ControlPlane) PullImage(ctx context.Context) error {
	if c.pulled {
		return nil
	}

	fmt.Println("pulling image...")
	res, err := c.cli.ImagePull(ctx, POSTGRES_IMAGE, image.PullOptions{})
	if err != nil {
		return err
	}
	defer res.Close()
	io.Copy(io.Discard, res)

	c.pulled = true
	fmt.Println("image pulled!")
	return nil
}

func (c *ControlPlane) AddInstance(ctx context.Context, name string) (Instance, error) {
	name = PREFIX + name

	portMap := nat.PortMap{
		"5432/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "0", // get any unused port
			},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portMap,
	}

	ctr, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: POSTGRES_IMAGE,
		Env:   ENVS[:],
		Cmd:   []string{"postgres", "-c", "wal_level=logical"},
	}, hostConfig, nil, nil, name)

	if err != nil {
		return Instance{}, err
	}

	err = c.cli.ContainerStart(ctx, ctr.ID, container.StartOptions{})

	if err != nil {
		return Instance{}, err
	}

	inspect, err := c.cli.ContainerInspect(ctx, ctr.ID)
	if err != nil {
		return Instance{}, err
	}

	portInfo := inspect.NetworkSettings.Ports["5432/tcp"][0].HostPort
	if portInfo == "" {
		return Instance{}, fmt.Errorf("failed to bind instance port for %v", name)
	}

	fmt.Printf("started container %v at port %v\n", name, portInfo)

	net, err := c.cli.NetworkCreate(ctx, name, network.CreateOptions{})

	if err != nil {
		return Instance{}, err
	}

	err = c.cli.NetworkConnect(ctx, net.ID, ctr.ID, nil)
	if err != nil {
		return Instance{}, err
	}

	inst := Instance{
		ContainerID: ctr.ID,
		NetworkID:   net.ID,
		Name:        name,
		Port:        portInfo,
	}

	c.servers = append(c.servers, inst)

	return inst, nil
}

func (c *ControlPlane) ListInstances(ctx context.Context) ([]Instance, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}
	networks, err := c.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	lkp := make(map[string]*Instance)

	for _, c := range containers {
		name := c.Names[0][1:]
		if !strings.HasPrefix(name, PREFIX) {
			continue
		}

		lkp[name] = &Instance{
			Name:        name,
			ContainerID: c.ID,
			Port:        fmt.Sprint(c.Ports[0].PublicPort),
		}
	}

	for _, n := range networks {
		if !strings.HasPrefix(n.Name, PREFIX) {
			continue
		}

		inst := lkp[n.Name]
		if inst == nil {
			panic(fmt.Errorf("found unused network: %v (%v)", n.Name, n.ID))
		}
		lkp[n.Name].NetworkID = n.ID
	}

	ret := make([]Instance, 0)
	for _, c := range lkp {
		if c.NetworkID == "" {
			panic(fmt.Errorf("found unconnected container: %v (%v)", c.Name, c.ContainerID))
		}
		ret = append(ret, *c)
	}

	return ret, nil
}

func (c *ControlPlane) KillInstance(ctx context.Context, inst Instance) error {
	err := c.cli.ContainerRemove(ctx, inst.ContainerID, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	fmt.Printf("killed container %v\n", inst.Name)

	err = c.cli.NetworkRemove(ctx, inst.NetworkID)
	if err != nil {
		return err
	}

	fmt.Printf("killed network %v\n", inst.Name)
	return nil
}

func (c *ControlPlane) Cleanup(ctx context.Context) error {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errs := make(chan error)

	for _, cont := range containers {
		name := cont.Names[0][1:]
		if !strings.HasPrefix(name, PREFIX) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := c.cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{
				Force: true,
			})
			if err != nil {
				errs <- err
			}
		}()
	}

	wg.Wait()

	networks, err := c.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, net := range networks {
		if !strings.HasPrefix(net.Name, PREFIX) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := c.cli.NetworkRemove(ctx, net.ID)
			if err != nil {
				errs <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		return err
	}

	fmt.Println("containers and networks cleaned...")
	return nil
}

func (c *ControlPlane) Connect(ctx context.Context, inst1 Instance, inst2 Instance) error {
	res := c.cli.NetworkConnect(ctx, inst1.NetworkID, inst2.ContainerID, nil)
	if res != nil {
		return res
	}
	fmt.Printf("connected %v to %v\n", inst1.Name, inst2.Name)
	return nil
}
