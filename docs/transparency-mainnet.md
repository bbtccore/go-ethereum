## Transparency Mainnet - Run and Operations Guide

This guide explains how to build, run, and operate a `geth` node configured for the Transparency main network embedded in this repository.

### Network parameters

- Chain ID: 210210
- Genesis hash: 0xb1d5c6c4d8b3d0447e7d3250c7be892fbec03c38e5a877a187d247dbc761a45b
- Default genesis: uses `TransparencyChainConfig` and `DefaultTransparencyGenesisBlock()`; running with no `--networkid` or testnet flags selects Transparency by default.
- Bootnodes: `params.TransparencyBootnodes` is currently empty (placeholders). Provide your own `--bootnodes` until official nodes are published.
- Unit aliases: includes `transparency`, `trcy` as 1e18 denominations in the embedded web3 console.

### Build

Prereqs: Go 1.23+ and a C compiler.

```bash
make geth
```

Binary will be at `build/bin/geth`.

### Quick start (single node)

Initialize and run with sane defaults. On first run, the node will use the Transparency genesis automatically.

```bash
./build/bin/geth \
  --datadir /var/lib/geth-transparency \
  --port 30303 \
  --http --http.addr 0.0.0.0 --http.port 8545 \
  --http.api eth,net,web3,engine,admin \
  --authrpc.addr 0.0.0.0 --authrpc.port 8551 \
  --authrpc.vhosts * \
  --authrpc.jwtsecret /var/lib/geth-transparency/jwt.hex \
  --bootnodes <comma-separated-enode-list>
```

Notes:
- If `--cache` is not provided and no testnet flags are set, `geth` bumps the cache to 4096MB on Transparency mainnet.
- Generate a JWT secret for Engine API:

```bash
openssl rand -hex 32 | tr -d '\n' > /var/lib/geth-transparency/jwt.hex
```

### Data directories

- Default datadir (Linux): `~/.ethereum` (Transparency uses the main profile unless a testnet flag is set). Use `--datadir` to isolate.
- IPC path: `<datadir>/geth.ipc` or override via `--ipcpath`.

### P2P networking

- Discovery/port: TCP 30303 by default. Use `--port` and `--discovery.port` to customize.
- Bootnodes: publish official Transparency bootnodes in `params/bootnodes.go`. Until then, pass `--bootnodes`.
- DNS-based node lists: `KnownDNSNetwork` supports `transparency-mainnet` if a DNS tree is published; otherwise leave blank.

### Recommended flags for production

```bash
--syncmode snap \
--cache 4096 \
--maxpeers 100 \
--http --http.addr 0.0.0.0 --http.api eth,net,web3 \
--ws --ws.addr 0.0.0.0 --ws.api eth,net,web3 \
--authrpc.addr 0.0.0.0 --authrpc.port 8551 --authrpc.jwtsecret /path/to/jwt.hex \
--metrics --metrics.addr 0.0.0.0 --metrics.port 6060
```

Security:
- Restrict RPC exposure and set explicit `--http.vhosts` and `--http.corsdomain` as needed.
- Do not unlock production accounts on the node. Use external signers (`clef`, HSM).

### Running in Docker

```bash
docker build -t transparency/geth .
docker run -d --name transparency-geth \
  -v /srv/geth:/data \
  -p 8545:8545 -p 8551:8551 -p 30303:30303/tcp -p 30303:30303/udp \
  transparency/geth \
  --datadir /data \
  --http --http.addr 0.0.0.0 --http.api eth,net,web3,engine \
  --authrpc.addr 0.0.0.0 --authrpc.port 8551 --authrpc.jwtsecret /data/jwt.hex \
  --bootnodes <enode-list>
```

### Beacon (Consensus) client

Transparency mainnet is post-merge (PoS). Pair `geth` with a consensus client (e.g., Prysm, Lighthouse, Teku, Nimbus) using Engine API on port 8551 with the same `jwt.hex`. Example (pseudocode):

```bash
consensus-client \
  --execution-endpoint http://EXECUTION_HOST:8551 \
  --jwt-secret /var/lib/geth-transparency/jwt.hex \
  --checkpoint-sync-url <optional>
```

Ensure both clients agree on the genesis and chain ID (210210).

### Operations

- Logs: use `--log.debug` for verbose logs during incident response.
- Snapshots and pruning: tune with `--snapshot`, `--gcmode`, and history flags as per storage requirements.
- Health checks: poll `eth_syncing`, `net_peerCount`, and `engine_*` endpoints.
- Backups: stop the node and back up `<datadir>`; for live backups, snapshot the underlying volume.

### Upgrades

- Keep Go toolchain at 1.23+ (see `go.mod`).
- Review release notes and fork schedules in `params/config.go` for any scheduled changes.

### Troubleshooting

- Stuck sync: verify bootnodes/peers and that consensus client is healthy; check Engine API auth.
- Genesis mismatch: remove old datadir or ensure it was initialized with Transparency genesis.
- Peer connectivity: open 30303 TCP/UDP. Verify NAT and `--nat` settings.

### References (source)

- `params/config.go` (Chain ID, forks, names)
- `core/genesis.go` (default genesis, Transparency defaults)
- `params/bootnodes.go` (bootnodes)
- `cmd/geth/main.go` (Transparency defaults and cache bump)

