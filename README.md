# ambd

[![Join the chat at https://gitter.im/jeffjen/ambd](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/jeffjen/ambd?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/jeffjen/go-libkv/master/LICENSE)

Ambassador Pattern made simple for micro services.

### TL;DR
- The Ambassador daemon `ambd`
```
ambd --addr 0.0.0.0:29091
```

- Runtime configuration client `ambctl`
    - `ambctl info`
    - `ambctl list`
    - `ambctl create --name mgo --src :27017 --dst mgos-ip-1:27017 --dst mgos-ip-2:27017`
    - `ambctl cancel --name mgo --name redis`

### What is Ambassador Pattern
[How To Use the Ambassador Pattern to Dynamically Configure Services](https://www.digitalocean.com/community/tutorials/how-to-use-the-ambassador-pattern-to-dynamically-configure-services-on-coreos)
portrays what this strategy could do to enable service discovery,
enhanced connectivity and routing pattern.

### Why ambd
- It is a light weight proxy daemon.
- Configuration is intuitive and scriptable.
- Resiliant to network partition and automatic retry.
- Available as docker image [jeffjen/ambd](https://hub.docker.com/r/jeffjen/ambd/)

Together with [docker](https://www.docker.com/) packaging and network facility,
we can deploy with confidence that code running in development environment will
continue to work in production environment.

See [setting up a test deployment](example/README.md).

