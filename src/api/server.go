package api

import (
	"context"
	"fmt"
	"net/http"
	"netpart/control"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func Run(ctx context.Context, addr string) {
	c, err := control.MakeControlPlane(ctx, client.FromEnv)
	if err != nil {
		panic(err)
	}
	err = c.PullImage(ctx)
	if err != nil {
		panic(err)
	}

	err = c.Cleanup(ctx)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Handle("/ping", pingHandler()).Methods("GET")
	r.Handle("/instances", listInstanceHandler(c)).Methods("GET")
	r.Handle("/instances", addInstanceHandler(c)).Methods("POST")
	r.Handle("/instances/{name}", killInstanceHandler(c)).Methods("DELETE")
	r.Handle("/instances/{name}", modifyInstanceHandler(c)).Methods("PUT")
	r.Handle("/instances/{name1}/connections/{name2}", connectHandler(c)).Methods("PUT")
	r.Handle("/instances/{name1}/connections/{name2}", disconnectHandler(c)).Methods("DELETE")
	r.Handle("/instances/{name}/keys", getKeysHandler(c)).Methods("GET")
	r.Handle("/instances/{name}/keys/{key}", putKeysHandler(c)).Methods("PUT")
	http.Handle("/", r)

	fmt.Printf("Listening at %v\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
