package api

import (
	"fmt"
	"net/http"
	"netpart/control"
)

type AddInstanceBody struct {
	Name string
}

type AddInstanceErrorResponse struct {
	Message string
}

type AddInstanceSuccessResponse = control.Instance

func addInstanceHandler(c *control.ControlPlane) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := decode[AddInstanceBody](r)
		if err != nil {
			encode(w, r, http.StatusBadRequest, AddInstanceErrorResponse{
				Message: "cannot decode request",
			})
			return
		}

		if body.Name == "" {
			encode(w, r, http.StatusBadRequest, AddInstanceErrorResponse{
				Message: "invalid instance name",
			})
			return
		}

		inst, err := c.AddInstance(ctx, body.Name)
		if err != nil {
			encode(w, r, http.StatusInternalServerError, AddInstanceErrorResponse{
				Message: "unknown error",
			})
			return
		}
		encode(w, r, http.StatusOK, inst)
	}
	return http.HandlerFunc(handler)
}

type ListInstanceResponse = []control.Instance

func listInstanceHandler(c *control.ControlPlane) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		insts, err := c.ListInstances(ctx)
		if err != nil {
			encode(w, r, http.StatusInternalServerError, AddInstanceErrorResponse{
				Message: "unknown error",
			})
			return
		}
		encode(w, r, http.StatusOK, insts)
	}
	return http.HandlerFunc(handler)
}

func pingHandler() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("yo")
		encode(w, r, http.StatusOK, struct {
			Message string
		}{
			Message: "OK",
		})
	}

	return http.HandlerFunc(handler)
}
