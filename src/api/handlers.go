package api

import (
	"fmt"
	"net/http"
	"netpart/control"

	"github.com/gorilla/mux"
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
		var resp AddInstanceErrorResponse

		body, err := decode[AddInstanceBody](r)
		if err != nil {
			resp.Message = "cannot decode request"
			encode(w, r, http.StatusBadRequest, resp)
			return
		}

		if body.Name == "" {
			resp.Message = "invalid instance name"
			encode(w, r, http.StatusBadRequest, resp)
			return
		}

		inst, err := c.AddInstance(ctx, body.Name)
		if err != nil {
			resp.Message = "unknown error"
			encode(w, r, http.StatusInternalServerError, resp)
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
		var resp AddInstanceErrorResponse

		insts, err := c.ListInstances(ctx)
		if err != nil {
			resp.Message = "unknown error"
			encode(w, r, http.StatusInternalServerError, resp)
			return
		}
		encode(w, r, http.StatusOK, insts)
	}
	return http.HandlerFunc(handler)
}

type ModifyInstanceBody struct {
	Primary bool

	Standby   bool
	StandbyTo string
}

type ModifyInstanceResponse struct {
	Message string
}

func modifyInstanceHandler(c *control.ControlPlane) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var resp ModifyInstanceResponse

		name := mux.Vars(r)["name"]
		inst, err := c.GetInstance(ctx, name)
		if err != nil {
			resp.Message = fmt.Sprintf("could not find instance %v", name)
			encode(w, r, http.StatusNotFound, resp)
			return
		}

		body, err := decode[ModifyInstanceBody](r)
		if body.Primary && body.Standby {
			resp.Message = "cannot set a node as primary and secondary"
			encode(w, r, http.StatusBadRequest, resp)
			return
		}

		if body.Primary {
			err = control.SetupPrimary(ctx, inst)
		} else if body.Standby {
			var primary control.Instance
			primary, err = c.GetInstance(ctx, body.StandbyTo)
			if err != nil {
				resp.Message = fmt.Sprintf("unable to find primary %v", primary)
				encode(w, r, http.StatusBadRequest, resp)
				return
			}
			err = control.SetupStandby(ctx, inst, primary)
		}

		if err != nil {
			resp.Message = "unable to set node"
			encode(w, r, http.StatusInternalServerError, resp)
			return
		}

		resp.Message = "OK"
		encode(w, r, http.StatusOK, resp)
	}

	return http.HandlerFunc(handler)
}

type KillInstanceResponse struct {
	Message string
}

func killInstanceHandler(c *control.ControlPlane) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		resp := &KillInstanceResponse{}

		name := mux.Vars(r)["name"]
		inst, err := c.GetInstance(ctx, name)
		if err != nil {
			resp.Message = fmt.Sprintf("could not find instance %v", name)
			encode(w, r, http.StatusNotFound, resp)
			return
		}

		err = c.KillInstance(ctx, inst)
		if err != nil {
			resp.Message = "unknown error"
			encode(w, r, http.StatusInternalServerError, resp)
			return
		}

		resp.Message = "OK"
		encode(w, r, http.StatusOK, resp)
	}
	return http.HandlerFunc(handler)
}

type ConnectResponse struct {
	Message string
}

func connectHandler(c *control.ControlPlane) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		resp := &ConnectResponse{}

		name1 := mux.Vars(r)["name1"]
		inst1, err := c.GetInstance(ctx, name1)
		if err != nil {
			resp.Message = fmt.Sprintf("could not find instance %v", name1)
			encode(w, r, http.StatusNotFound, resp)
			return
		}

		name2 := mux.Vars(r)["name2"]
		inst2, err := c.GetInstance(ctx, name2)
		if err != nil {
			resp.Message = fmt.Sprintf("could not find instance %v", name2)
			encode(w, r, http.StatusNotFound, resp)
			return
		}

		err = c.Connect(ctx, inst1, inst2)
		if err != nil {
			resp.Message = "unknown error"
			encode(w, r, http.StatusInternalServerError, resp)
			return
		}

		resp.Message = "OK"
		encode(w, r, http.StatusOK, resp)
	}
	return http.HandlerFunc(handler)
}

type DisconnectResponse struct {
	Message string
}

func disconnectHandler(c *control.ControlPlane) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		resp := &DisconnectResponse{}

		name1 := mux.Vars(r)["name1"]
		inst1, err := c.GetInstance(ctx, name1)
		if err != nil {
			resp.Message = fmt.Sprintf("could not find instance %v", name1)
			encode(w, r, http.StatusNotFound, resp)
			return
		}

		name2 := mux.Vars(r)["name2"]
		inst2, err := c.GetInstance(ctx, name2)
		if err != nil {
			resp.Message = fmt.Sprintf("could not find instance %v", name2)
			encode(w, r, http.StatusNotFound, resp)
			return
		}

		err = c.Disconnect(ctx, inst1, inst2)
		if err != nil {
			resp.Message = "unknown error"
			encode(w, r, http.StatusInternalServerError, resp)
			return
		}

		resp.Message = "OK"
		encode(w, r, http.StatusOK, resp)
	}
	return http.HandlerFunc(handler)
}

type GetKeysSuccessResponse = []control.KV
type GetKeysFailResponse struct {
	Message string
}

func getKeysHandler(c *control.ControlPlane) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var resp GetKeysFailResponse

		name := mux.Vars(r)["name"]
		inst, err := c.GetInstance(ctx, name)
		if err != nil {
			resp.Message = fmt.Sprintf("could not find instance %v", name)
			encode(w, r, http.StatusNotFound, resp)
			return
		}

		vals, err := control.Get(ctx, inst)
		if err != nil {
			resp.Message = "unable to get values"
			encode(w, r, http.StatusInternalServerError, resp)
			return
		}

		encode(w, r, http.StatusOK, vals)
	}

	return http.HandlerFunc(handler)
}

func pingHandler() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		encode(w, r, http.StatusOK, struct {
			Message string
		}{
			Message: "OK",
		})
	}

	return http.HandlerFunc(handler)
}
