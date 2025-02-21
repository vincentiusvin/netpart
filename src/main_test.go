package main

import (
	"context"
	"fmt"
	"netpart/control"
	"testing"

	"github.com/docker/docker/client"
)

func getTestControlPlane() (*control.ControlPlane, error) {
	ctx := context.Background()
	c, err := control.MakeControlPlane(ctx, client.FromEnv)
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

func TestProvision(t *testing.T) {
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

func TestConnection(t *testing.T) {
	ctx := context.Background()
	c, err := getTestControlPlane()
	if err != nil {
		t.Fatal(err)
	}

	names := []string{
		"db1", "db2", "db3",
	}

	instances := make([]control.Instance, len(names))

	for i, name := range names {
		inst, err := c.AddInstance(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		instances[i] = inst
	}

	err = c.Connect(ctx, instances[0], instances[1])
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabase(t *testing.T) {
	ctx := context.Background()
	c, err := getTestControlPlane()
	if err != nil {
		t.Fatal(err)
	}

	active, err := c.AddInstance(ctx, "db1")
	if err != nil {
		t.Fatal(err)
	}

	passive, err := c.AddInstance(ctx, "db2")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Connect(ctx, active, passive)
	if err != nil {
		t.Fatal(err)
	}

	err = control.SetupPrimary(ctx, active)
	if err != nil {
		t.Fatal(err)
	}

	err = control.SetupStandby(ctx, passive, active)
	if err != nil {
		t.Fatal(err)
	}

	in_key := "test"
	in_value := "val"

	err = control.Put(ctx, active, in_key, in_value)
	if err != nil {
		t.Fatal(err)
	}

	err = findVal(active, in_key, in_value)
	if err != nil {
		t.Fatal("failed to find data on primary")
	}

	err = findVal(active, in_key, in_value)
	if err != nil {
		t.Fatal("failed to find data on standby")
	}

}

func findVal(inst control.Instance, key string, value string) error {
	ctx := context.Background()
	val, err := control.Get(ctx, inst)
	if err != nil {
		return err
	}

	inserted := false
	for _, e := range val {
		found := e.Key == key && e.Value == value
		if found {
			inserted = true
			break
		}
	}

	if !inserted {
		return fmt.Errorf("failed to find inserted data")
	}

	return nil
}
