package main

import (
	"context"
	"fmt"
	"netpart/control"
)

func main() {
	ctx := context.Background()
	c := control.MakeControlPlane(ctx)
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
}
