package control

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const POSTGRES_USER = "postgres"
const POSTGRES_PASSWORD = "postgres"
const POSTGRES_DB = "main"

func getConn(ctx context.Context, port string) (*pgx.Conn, error) {
	for {
		connString := "postgresql://" + POSTGRES_USER + ":" + POSTGRES_PASSWORD + "@dind:" + port + "/" + POSTGRES_DB
		conn, err := pgx.Connect(ctx, connString)
		if err != nil {
			fmt.Println("pinging database failed, retrying...")
			time.Sleep(500 * time.Millisecond)
		} else {
			fmt.Println("database connected!")
			return conn, nil
		}
	}
}

const DDL = "CREATE TABLE IF NOT EXISTS kv ( key text PRIMARY KEY, value text );"

func SetupDB(ctx context.Context, inst Instance) error {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, DDL)
	if err != nil {
		return err
	}
	return nil
}

const PUB = "CREATE PUBLICATION pub FOR TABLE kv;"

func SetupPrimary(ctx context.Context, inst Instance) error {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, PUB)
	if err != nil {
		return err
	}

	fmt.Printf("active setup at %v\n", inst.Name)
	return nil
}

func SetupStandby(ctx context.Context, inst Instance, active Instance) error {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	// replication slot name can only be numbers, alpha, and underscores.
	sanitized_subscription := strings.ReplaceAll("sub_"+inst.Name, "-", "_")

	// TODO:
	// this is vulnerable to sql injection actually
	// but you can't turn create subscription into a prepared statement
	sub := fmt.Sprintf(
		"CREATE SUBSCRIPTION \"%v\" CONNECTION 'host=%v dbname=%v user=%v password=%v' PUBLICATION pub WITH (disable_on_error = true);",
		sanitized_subscription, active.Name, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD)

	_, err = conn.Exec(ctx, sub)
	if err != nil {
		return err
	}

	fmt.Printf("standby setup at %v\n", inst.Name)
	return nil
}

func RestartStandby(ctx context.Context, inst Instance, active Instance) error {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	// replication slot name can only be numbers, alpha, and underscores.
	sanitized_subscription := strings.ReplaceAll("sub_"+inst.Name, "-", "_")

	// TODO:
	// this is vulnerable to sql injection actually
	// but you can't turn create subscription into a prepared statement
	sub := fmt.Sprintf("ALTER SUBSCRIPTION \"%v\" ENABLE", sanitized_subscription)

	_, err = conn.Exec(ctx, sub)
	if err != nil {
		return err
	}

	fmt.Printf("standby restarted at %v\n", inst.Name)
	return nil
}

func Put(ctx context.Context, inst Instance, key string, value string) error {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "INSERT INTO kv (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = $2", key, value)
	if err != nil {
		return err
	}

	return nil
}

type KV struct {
	Key   string
	Value string
}

func Get(ctx context.Context, inst Instance) ([]KV, error) {
	conn, err := getConn(ctx, inst.Port)
	if err != nil {
		return nil, err
	}

	defer conn.Close(ctx)

	val, err := conn.Query(ctx, "SELECT key, value FROM kv ORDER BY key ASC")
	if err != nil {
		return nil, err
	}

	kv, err := pgx.CollectRows(val, pgx.RowToStructByName[KV])
	if err != nil {
		return nil, err
	}

	return kv, nil
}

type ActiveData struct {
	Application_Name string
	State            string
	Sync_State       string
}

type StandbyData struct {
	Subname    string
	Subenabled bool
}

type ReplicationData struct {
	ActiveData  []ActiveData
	StandbyData []StandbyData
}

func GetReplicationData(ctx context.Context, inst Instance) (ReplicationData, error) {
	conn, err := getConn(ctx, inst.Port)
	resp := ReplicationData{}
	if err != nil {
		return resp, err
	}

	defer conn.Close(ctx)

	active_raw, err := conn.Query(ctx, "SELECT application_name, state, sync_state FROM pg_stat_replication;")
	if err != nil {
		return resp, fmt.Errorf("can't run active query: %w", err)
	}

	active_data, err := pgx.CollectRows(active_raw, pgx.RowToStructByName[ActiveData])
	pgx.CollectRows(active_raw, pgx.RowTo[ActiveData])
	if err != nil {
		return resp, fmt.Errorf("can't marshall active query: %w", err)
	}

	standby_raw, err := conn.Query(ctx, "SELECT subname, subenabled FROM pg_subscription;")
	if err != nil {
		return resp, fmt.Errorf("can't run standby query: %w", err)
	}
	standby_data, err := pgx.CollectRows(standby_raw, pgx.RowToStructByName[StandbyData])
	if err != nil {
		return resp, fmt.Errorf("can't marshall standby query: %w", err)
	}

	return ReplicationData{
		ActiveData:  active_data,
		StandbyData: standby_data,
	}, nil
}
