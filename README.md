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

Explored using tc to apply network delays, but it only works for outgoing packets.
https://serverfault.com/questions/1150987/both-incoming-and-outgoing-packets-is-delayed-even-though-i-target-incoming-only

command:

```bash
tc qdisc add dev eth2 ingress netem delay 10000ms
```

I think justify making a tcp proxy now.
