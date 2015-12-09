# ambd
The Ambassador pattern to container services

### TL;DR
- The Ambassador daemon `ambd`
```
ambd --addr 0.0.0.0:29091 --advertise advertise-host-ip:29091 \
    --prefix /debug/docker/ambassador/nodes \
    --proxy '{"name": "mysql", "net": "tcp", "src": ":3306", "dst": ["mysql-master-ip:3306", "mysql-slave-ip:3306"]}' \
    --proxy '{"name": "redis", "net": "tcp", "src": ":6379", "srv": "/srv/redis/debug"}' \
    --proxycfg /proxy/debug/v1 \
    etcd://etcd1-ip:2379,etcd2-ip:2379,etcd3-ip:2379
```

- Runtime configuration client `ambctl`
    - list: `ambctl list`
    - create: `ambctl create --name mgo --src :3376 --dst mgos-ip-1:3376 --dst mgos-ip-2:3376`
    - cancel: `ambctl cancel --name mgo --name redis`
    - info: `ambctl info`
    - config: `ambctl config /proxy/debug/v2`

- Available as docker image [jeffjen/ambd](https://hub.docker.com/r/jeffjen/ambd/)

### What is Ambassador Pattern
[How To Use the Ambassador Pattern to Dynamically Configure Services](https://www.digitalocean.com/community/tutorials/how-to-use-the-ambassador-pattern-to-dynamically-configure-services-on-coreos)
portrays what this strategy could do to enable service discovery,
enhanced connectivity and routing pattern.

### Why ambd
- It is a light weight proxy daemon.
- Configuration is intuitive and shell scriptable.
- Resiliant to network partition and automatic retry.

Together with [docker](https://www.docker.com/) packaging and network facility,
we can deploy with confidence that code running the development environment will
continue to work in production environment

See [setting up a working deployment](example/README.md) for a walkthrough over the system

