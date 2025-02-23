Does three things:

1. Provision postgres databases
2. Setup replication for them
3. Force network partitions 😈

Runs docker-in-docker for provisioning and controlling the networking.

Replication uses postgres' logical replication.

Main environment:

```bash
docker compose -f compose.yaml up --build
```

Test environment:

```bash
docker compose -f compose.yaml -f compose.test.yaml up --build
```
