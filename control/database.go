package control

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
)

const SQL = "CREATE TABLE IF NOT EXISTS kv ( key text PRIMARY KEY, value text );"

func (c *ControlPlane) runCommand(ctx context.Context, inst Instance, cmd []string) error {
	conf, err := c.cli.ContainerExecCreate(ctx, inst.ContainerID, container.ExecOptions{
		User: POSTGRES_USER,
		Cmd:  cmd,
	})
	if err != nil {
		return err
	}

	err = c.cli.ContainerExecStart(ctx, conf.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		for {
			ins, err := c.cli.ContainerExecInspect(ctx, conf.ID)
			if ins.Running {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil {
				done <- err
			} else if ins.ExitCode != 0 {
				done <- fmt.Errorf("failed command for %v. pid: %v. exit code: %v", inst.Name, ins.Pid, ins.ExitCode)
			} else {
				close(done)
			}
			break
		}
	}()

	return <-done
}

func (c *ControlPlane) SetupMaster(ctx context.Context, inst Instance) error {
	for {
		err := c.runCommand(ctx, inst, []string{"psql", "-l"})
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	err := c.runCommand(ctx, inst, []string{"psql", "-d", POSTGRES_DB, "-c", SQL})
	if err != nil {
		return err
	}
	fmt.Printf("master setup at %v (%v)\n", inst.Name, inst.ContainerID)
	return nil
}
