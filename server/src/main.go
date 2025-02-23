package main

import (
	"context"
	"netpart/api"
)

func main() {
	ctx := context.Background()
	api.Run(ctx, ":7000")
}
