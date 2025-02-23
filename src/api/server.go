package api

import (
	"context"
	"fmt"
	"net/http"
	"netpart/control"

	"github.com/docker/docker/client"
)

func Run(ctx context.Context, addr string) {
	c, err := control.MakeControlPlane(ctx, client.FromEnv)
	if err != nil {
		panic(err)
	}

	http.Handle("GET /ping", pingHandler())
	http.Handle("GET /instances", listInstanceHandler(c))
	http.Handle("POST /instances", addInstanceHandler(c))

	fmt.Printf("Listening at %v\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
