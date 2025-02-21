package control

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const POSTGRES_USER = "postgres"
const POSTGRES_PASSWORD = "postgres"
const POSTGRES_DB = "main"

func getConn(ctx context.Context, port string) (*pgx.Conn, error) {
	connString := "postgresql://" + POSTGRES_USER + ":" + POSTGRES_PASSWORD + "@localhost:" + port + "/" + POSTGRES_DB
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

const DDL = "CREATE TABLE IF NOT EXISTS kv ( key text PRIMARY KEY, value text );"
const PUB = "CREATE PUBLICATION pub FOR TABLE kv;"

func (c *ControlPlane) SetupActive(ctx context.Context, inst Instance) error {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, DDL)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, PUB)
	if err != nil {
		return err
	}

	fmt.Printf("active setup at %v\n", inst.Name)
	return nil
}

func (c *ControlPlane) SetupStandby(ctx context.Context, inst Instance, active Instance) error {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, DDL)
	if err != nil {
		return err
	}

	sub := fmt.Sprintf(
		"CREATE SUBSCRIPTION sub CONNECTION 'host=%v dbname=%v user=%v password=%v' PUBLICATION pub;",
		active.Name, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD)

	_, err = conn.Exec(ctx, sub)
	if err != nil {
		return err
	}

	fmt.Printf("standby setup at %v\n", inst.Name)
	return nil
}
