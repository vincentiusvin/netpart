package main

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
)

const POSTGRES_IMAGE = "postgres:16.3-alpine3.20"
const PREFIX = "netpart-"

var ENVS = [3]string{
	"POSTGRES_USER=postgres",
	"POSTGRES_PASSWORD=postgres",
	"POSTGRES_DB=main",
}

type Instance struct {
	name         string
	container_id string
	network_id   string
}

type ControlPlane struct {
	cli     *client.Client
	pulled  bool
	servers []Instance
}

func MakeControlPlane(ctx context.Context) *ControlPlane {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	for {
		_, err := cli.Ping(ctx)
		if err != nil {
			fmt.Println(err)
			time.Sleep(500 * time.Millisecond)
		} else {
			fmt.Println("Docker daemon connected!")
			break
		}
	}

	return &ControlPlane{
		cli: cli,
	}
}

func (c *ControlPlane) PullImage(ctx context.Context) {
	if c.pulled {
		return
	}

	fmt.Println("Pulling image...")
	res, err := c.cli.ImagePull(ctx, POSTGRES_IMAGE, image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer res.Close()
	io.Copy(io.Discard, res)

	c.pulled = true
	fmt.Println("Image pulled!")
}

func (c *ControlPlane) AddInstance(ctx context.Context, name string) Instance {
	name = PREFIX + name

	ctr, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: POSTGRES_IMAGE,
		Env:   ENVS[:],
	}, nil, nil, nil, name)

	if err != nil {
		panic(err)
	}

	err = c.cli.ContainerStart(ctx, ctr.ID, container.StartOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Started container %v (%v)\n", name, ctr.ID)

	net, err := c.cli.NetworkCreate(ctx, name, network.CreateOptions{})

	if err != nil {
		panic(err)
	}

	err = c.cli.NetworkConnect(ctx, net.ID, ctr.ID, nil)
	if err != nil {
		panic(err)
	}

	inst := Instance{
		container_id: ctr.ID,
		network_id:   net.ID,
		name:         name,
	}

	c.servers = append(c.servers, inst)

	return inst
}

func (c *ControlPlane) ListInstances(ctx context.Context) []Instance {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}
	networks, err := c.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		panic(err)
	}

	lkp := make(map[string]*Instance)

	for _, c := range containers {
		name := c.Names[0][1:]
		if !strings.HasPrefix(name, PREFIX) {
			continue
		}

		lkp[name] = &Instance{
			name:         name,
			container_id: c.ID,
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
		lkp[n.Name].network_id = n.ID
	}

	ret := make([]Instance, 0)
	for _, c := range lkp {
		if c.network_id == "" {
			panic(fmt.Errorf("found unconnected container: %v (%v)", c.name, c.container_id))
		}
		ret = append(ret, *c)
	}

	return ret
}

func (c *ControlPlane) KillInstance(ctx context.Context, inst Instance) {
	err := c.cli.NetworkRemove(ctx, inst.network_id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Killed network %v\n", inst.network_id)

	err = c.cli.ContainerRemove(ctx, inst.container_id, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Killed container %v\n", inst.container_id)
}

func (c *ControlPlane) Cleanup(ctx context.Context) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, cont := range containers {
		name := cont.Names[0][1:]
		if !strings.HasPrefix(name, PREFIX) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{
				Force: true,
			})
		}()
	}

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
			c.cli.NetworkRemove(ctx, net.ID)
		}()
	}

	wg.Wait()
	fmt.Println("Containers and networks cleaned...")
}

func main() {
	ctx := context.Background()
	c := MakeControlPlane(ctx)
	c.PullImage(ctx)
	c.Cleanup(ctx)

	is := c.ListInstances(ctx)
	fmt.Println(is)

	c.AddInstance(ctx, "udin")
	is = c.ListInstances(ctx)
	fmt.Println(is)

	c.AddInstance(ctx, "samsul")
	is = c.ListInstances(ctx)
	fmt.Println(is)

	// id1 := c.AddInstance(ctx, "db1")
	// c.AddInstance(ctx, "db2")
	// c.AddInstance(ctx, "db3")
	// c.AddInstance(ctx, "db4")
	// c.KillInstance(ctx, id1)
}
