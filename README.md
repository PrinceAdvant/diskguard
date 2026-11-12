# diskguard

A tiny Go CLI and library that checks free disk space and free inodes for the filesystems behind one or more paths, and exits non-zero when either drops below a threshold. Standard library only, Linux, macOS and FreeBSD.

## Why

Docker hosts rarely fail loudly when the disk fills up. Builds fail with cryptic errors, Postgres stops accepting writes and containers restart in loops. Old images, build cache and unrotated container logs under `/var/lib/docker` are the usual cause. A check that runs every few minutes and alerts at 15 percent free gives you time to clean up first.

## Install

```
go install github.com/PrinceAdvant/diskguard/cmd/diskguard@latest
```

## Usage

```
diskguard                       # checks / and /var/lib/docker
diskguard -min 10 / /srv/data
```

Each path prints `ok` or `LOW` with free space and free inodes. The exit code is 0 when everything is above the threshold and 1 otherwise, so it fits in cron, a systemd timer or a monitoring heartbeat script:

```
*/5 * * * * diskguard -min 15 || curl -fsS https://monitor.example.com/alert/disk
```

Inodes are checked too, because millions of tiny files in build caches can exhaust inodes while plenty of bytes remain.

## Cleaning up

When it fires, `docker system df` shows where space went, `docker builder prune` and `docker image prune -a` reclaim most of it, and log rotation in `daemon.json` stops it coming back. A full guide on what to do when [Docker runs out of disk space](https://peon.sh/blogs/docker-out-of-disk-space) covers each step. Built by a contributor to [Peon](https://peon.sh/), an open-source self-hosted deployment platform.

MIT licensed.
