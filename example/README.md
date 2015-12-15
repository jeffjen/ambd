# Setting up a test deployment
Lets go over a typical setup using [docker](https://www.docker.com/) and
[docker-compose](https://docs.docker.com/compose/).

### Have questions?
- Open an issue
- Ask on [gitter](https://gitter.im/jeffjen/ambd?utm_source=share-link&utm_medium=link&utm_campaign=share-link)

### Composition and definition
Create a custom bridge network with docker
```
docker network create --driver bridge isolated_nw
```

Service stack is defined in [docker-compose.yml](docker-compose.yml).

### Starting up
Launch the stack by executing `docker-compose up`.  This will take a while to
download necessary images and building our test app.

Your app after staring will be complaining it could not reach a Redis node.
Lets fix this by configuring `ambassador` to point to a Redis endpoint.

In another terminal, copy `ambctl` from the container
```
docker cp ambassador:/ambctl ./ambctl
```

Configure a route to Redis node
```
./ambctl create --name redis --src :6379 --dst redis:6379
```

Your app should quickly start reporting the role of the Redis it connected to,
but why go through all this trouble when you could specify `redis:6379` to
connect in your app?

The reason is configurability and moving between environment.  When I ship my
service, it needs to talk to our production Redis node, but I don't want to
stop my service to change code and/or redeploy.  I could configure `ambd` by
```
./ambctl cancel --name redis
./ambctl create --name redis --src :6379 --dst production-redis:6379
```

After a little hiccup, the app starts talking to `production-redis:6379`
without code change nor deployment.

### Advanced Configuration
`ambd` true triumphs comes from integrating with popular discovery backend
[etcd](https://github.com/coreos/etcd).  Static IP based route assignment is
greate for testing, but it quickly becomes hard to track which host represents
what service as you add more nodes in your infrastructure.  A common strategy
is to tag or group services by name, and in `etcd`, you place nodes
representing a service under a directory (a service key) e.g. /srv/redis/debug, or
/srv/mongodb/staging.

`ambd` allows you to set route to these services by referring to the directory
they are in.  Whenever a node joins or leaves it's service, `ambd` will get
notification to obtain the latest set of available nodes.  Thereby removing the
need to reconfigure `ambd`.

To demonstrate `ambd` with discovery backend, stop and remove the stack we
initiated by `Ctrl-C`, followed by `docker-compose rm -f`.

### Composition and definition
In order for the test to work, please add `127.0.0.1 ambassador` to your
`/etc/hosts`.  This is required because we are testing on a single machine, and
advertising becomes difficult with containers when not knowing interface IP
ahead of time.

Service stack is defined in [docker-compose.etcd.yml](docker-compose.etcd.yml).

### Starting up
Launch the stack with `docker-compose -f docker-compose.etcd.yml up`.

Again upon launch your app will start complaining it could not reach a Redis
node.  Fix this by configuring `ambd` to point to a Redis endpoint.
```
./ambctl create --name redis --src :6379 --srv /srv/redis/debug
```

Notice that we have a little help with
[docker-monitor](https://github.com/jeffjen/docker-monitor), which provides
regitration to discovery backend if service declare themselves by label.  If
you disable docker-montior, you need to register Redis node to
`/srv/redis/debug` manually through `etcdctl` command.
```
docker cp discovery:/etcdctl ./etcdctl
etcdctl set /srv/redis/debug/127.0.0.1:16379 127.0.0.1:16379
```

### Using config key for consistent deployment
You are statisfied with the configuration, now its time to go big.  Set a
config key for `ambd` to obtain and follow.

Suppose our runtime environment is required to be:
```
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

Store this document to `cfg.json`

**Tip**: use jq to manipulate json documents on the command line:
```
curl -sSL https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 -o jq
chmod +x ./jq
```

Send config to discovery backend:
```
./etcdctl set /proxy/myapp/v1 "$(jq -c . <cfg.json)"
```

Launch ambassador to run with config key:
```
docker run -d --name ambassador -m 128M -p 29091:29091 \
    jeffjen/ambd \
        --addr 0.0.0.0:29091 \
        --cluster debug \
        --advertise HostIP:29091 \
        --proxycfg /proxy/myapp/v1 \
        etcd://discovery-001:2379,discovery-002:2379,discovery-003:2379
```

If chanages are made to config key `/proxy/maypp/v1`, all `ambd` following that
key will react to obtain the latest setting.  If you prefer strict
versioning, submit new config with key `/proxy/myapp/v2`, and execute the
following command:
```
ambctl volly --cluster debug --dsc http://discovery-001:2379 config /proxy/config/v2
```

All `ambd` nodes registered under cluster `debug` will obtain and follow new
configuration from `/proxy/config/v2`.

