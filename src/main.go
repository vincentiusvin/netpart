package main

import (
	"context"
	"netpart/control"

	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	_, err := control.MakeControlPlane(ctx, client.FromEnv)
	if err != nil {
		panic(err)
	}
}
