# ambd

[![Join the chat at https://gitter.im/jeffjen/ambd](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/jeffjen/ambd?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/jeffjen/go-libkv/master/LICENSE)

**ambd** makes Ambassador Pattern simple for *Serivce Oriented Architecture*,
by *bridging* micro services through generic point to point connector, or by
meaningful *context*, backed by discovery service.

**ambd** connectivity logic is provided by [go-proxy](https://github.com/jeffjen/go-proxy).

See [setting up a test deployment](example/README.md).

### Quick Start
- The Ambassador daemon `ambd`
```
ambd --addr 0.0.0.0:29091
```

- Runtime configuration client `ambctl`
    - `ambctl info`
    - `ambctl list`
    - `ambctl create --name mgo --src :27017 --dst mgos-ip-1:27017 --dst mgos-ip-2:27017`
    - `ambctl cancel --name mgo`

### What is Ambassador Pattern
[How To Use the Ambassador Pattern to Dynamically Configure Services](http://do.co/1J99qO3)
describes what this strategy could do to enable service discovery, connectivity
and better routing pattern, without extensive network connection logic
implemented in each service node.

### Why ambd
- It is a light weight proxy daemon.
- With runtime configuration client that is intuitive and scriptable.
- Resilient to network partition and retry.
- Docker image available [jeffjen/ambd](https://hub.docker.com/r/jeffjen/ambd/)

Together with [docker](https://www.docker.com/) packaging and network facility,
we can deploy with confidence that code running in development environment will
continue to work in production environment.

See [setting up a test deployment](example/README.md).

### Have questions?
- Open an issue
- Ask on [gitter](https://gitter.im/jeffjen/ambd?utm_source=share-link&utm_medium=link&utm_campaign=share-link)

