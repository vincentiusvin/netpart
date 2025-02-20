package main

import (
	"context"
	"netpart/control"
)

func main() {
	ctx := context.Background()
	_, err := control.MakeControlPlane(ctx)
	if err != nil {
		panic(err)
	}
}
