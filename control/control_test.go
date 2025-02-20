package control_test

import (
	"context"
	"netpart/control"
	"testing"

	"github.com/docker/docker/client"
)

const TEST_SOCKET = "unix://../sock/docker.sock"

func getTestControlPlane() (*control.ControlPlane, error) {
	ctx := context.Background()
	c, err := control.MakeControlPlane(ctx, client.WithHost(TEST_SOCKET))
	if err != nil {
		return nil, err
	}
	err = c.Cleanup(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func TestCleanup(t *testing.T) {
	ctx := context.Background()
	c, err := getTestControlPlane()
	if err != nil {
		t.Fatal(err)
	}
	c.Cleanup(ctx)
	c.ListInstances(ctx)
}

func TestMain(t *testing.T) {
	ctx := context.Background()
	c, err := getTestControlPlane()
	if err != nil {
		t.Fatal(err)
	}

	in_name := "udin"

	var inst control.Instance

	t.Run("can add instance", func(t *testing.T) {
		inst, err = c.AddInstance(ctx, in_name)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("can list instances", func(t *testing.T) {
		list_inst, err := c.ListInstances(ctx)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, r := range list_inst {
			if r.Name == inst.Name {
				found = true
			}
		}

		if !found {
			t.Errorf("failed to find added instance %v", in_name)
		}
	})

	t.Run("can delete instance", func(t *testing.T) {
		err := c.KillInstance(ctx, inst)
		if err != nil {
			t.Fatal(err)
		}
	})
}
