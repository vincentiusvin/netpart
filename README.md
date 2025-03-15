Does three things:

1. Provision postgres databases
2. Setup replication for them
3. Force network partitions 😈

Runs docker-in-docker for provisioning and controlling the networking.

Replication uses postgres' logical replication.

# Run

```bash
docker compose up --build
```

Application is accessible using localhost:7000

# Notes

Explored using tc to apply network delays, but it only works for outgoing packets.

https://serverfault.com/questions/1150987/both-incoming-and-outgoing-packets-is-delayed-even-though-i-target-incoming-only

command:

```bash
tc qdisc add dev eth2 ingress netem delay 10000ms
```

I think this justifies making a tcp proxy now.
