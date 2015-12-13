## Setting up a working deployment
Lets go over a typical setup using [docker-compose](https://docs.docker.com/compose/)

You will need [docker](https://www.docker.com/) installed.

### Composition and definition
Create a custom bridge network with docker
```
docker network create --driver bridge isolated_nw
```

Define and setup [ambassador](docker-compose.yml)

### Starting up
Launch the stack with `docker-compose up`.  This will take a while.

Your app will now be complaining it could not reach a Redis node once the stack
is launched.  Lets fix this by configuring `ambassador` to point to a Redis
endpoint.

In another terminal, copy `ambctl` from the container
```
docker cp ambassador:/ambctl ./ambctl
```

Configure a route to Redis node
```
./ambctl create --name redis --src :6379 --dst redis:6379
```

But why go through all this trouble when you could  just as easily specify
`redis:6379` to connect in your app?  Suppose I want to talk to a different
Redis node, but I don't want to stop my app to change code and/or redeploy.
I could configure `ambassador` by
```
./ambctl cancel --name redis
./ambctl create --name redis --src :6379 --dst remote-redis:6379
```

After a little hiccup, the app starts talking to remote-redis endpoint without
any code change, no deployment required.

### Advanced Configuration
`ambd` true value comes from integrating with popular discovery backend
[etcd](https://github.com/coreos/etcd).  Static IP based route assignment is
greate for testing, but it quickly becomes hard to track which host represents
what service as you add more nodes in your infrastructure.  A common strategy
is to tag or group services by name, and in `etcd`, you place nodes
representing a service under a directory (a service key) e.g. /srv/redis/debug, or
/srv/mongodb/staging.

`ambd` allows you to set route to these services by referring to the directory
they are under.  Whenever a node joins or leaves it's service, `ambd` will react
to obtain the latest set of available nodes.  Thereby removing the need to
reconfigure `ambd` to point to a new host.

To demonstrate `ambd` with discovery backend, stop and remove the stack we
initiated by `Ctrl-C` and `docker-compose rm -f`

### Composition and definition
In order for the test to work, please add an entry to your `/etc/hosts` to
include `127.0.0.1 ambassador`.  This is required because we are testing the
configuration on a single machine, and advertising becomes difficult when
launching with containers and not able to know the interface IP ahead of time.

Define and setup [ambassador with discovery backend](docker-compose.etcd.yml)

### Starting up
Launch the stack with `docker-compose -f docker-compose.etcd.yml up`.
This will take a while.

And again your app will now be complaining it could not reach a Redis node.
Lets fix this by configuring `ambassador` to point to a Redis endpoint.
```
./ambctl create --name redis --src :6379 --srv /srv/redis/debug
```

You can reassign which endpoint to route to by `cancel` followed by `create`.

### Using config key for mass deployment
You are now statisfied with the configuration, now its time to go to production.
To command all of the `ambassador` to behave in the same way, you could set
a config key for `ambassador` to obtain and follow.

Suppose our runtime environment is required to be
```yml
[
    {
        "name": "cache",
        "net": "tcp4",
        "src": ":6379",
        "srv": "/srv/redis/cache"
    },
    {
        "name": "db",
        "net": "tcp4",
        "src": ":27017"
        "dst": [
            "mongo-router-01:27017",
            "mongo-router-02:27017",
            "mongo-router-03:27017"
        ]
    },
]
```

You would use `etcdctl` to submit this json document to `/proxy/myapp/v1`.

Launch ambassador to run with config key:
```
docker run -d --name ambassador -m 128M -p 29091:29091 \
    jeffjen/ambd \
        --addr 0.0.0.0:29091 \
        --prefix /debug/docker/ambassador/nodes \
        --advertise HostIP:29091 \
        --proxycfg /proxy/myapp/v1 \
        etcd://discovery-001:2379,discovery-002:2379,discovery-003:2379
```

Whenever you decide that to add, move, or update endpoint route, you resubmit
your configuration using key `/proxy/myapp/v2`.  Then, instruct the nodes that
are affected by
```
ambctl volly --cluster debug --dsc http://discovery-001:2379 config /proxy/config/v2
```
