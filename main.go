package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

const POSTGRES_IMAGE = "postgres:16.3-alpine3.20"

var ENVS = [3]string{
	"POSTGRES_USER=postgres",
	"POSTGRES_PASSWORD=postgres",
	"POSTGRES_DB=main",
}

type ControlPlane struct {
	cli    *client.Client
	pulled bool
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

func (c *ControlPlane) AddDB(ctx context.Context, name string) string {
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

	fmt.Printf("Started container %v\n", ctr.ID)
	return ctr.ID
}

func (c *ControlPlane) ListDBs(ctx context.Context) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		fmt.Printf("%s %s\n", ctr.ID, ctr.Image)
	}
}

func (c *ControlPlane) KillDB(ctx context.Context, id string) {
	err := c.cli.ContainerRemove(ctx, id, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Killed container %v\n", id)
}

func main() {
	ctx := context.Background()
	c := MakeControlPlane(ctx)
	c.PullImage(ctx)

	id := c.AddDB(ctx, "db1")
	c.ListDBs(ctx)
	c.KillDB(ctx, id)
}
