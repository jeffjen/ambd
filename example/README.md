## Setting up a working deployment
Lets go over a typical setup using [docker-compose](https://docs.docker.com/compose/)

### Composition and definition
Create a custom bridge network with docker
```
docker network create --driver bridge isolated_nw
```

Define and setup [ambassador](docker-compose.yml)
```yml
ambassador:
    container_name: ambassador
    command:
        - "--addr"
        - "0.0.0.0:29091"
    environment:
        LOG_LEVEL: INFO
    image: jeffjen/ambd
    net: isolated_nw
    ports:
        - "29091:29091"
```

Copy the `ambctl` from the container
```
docker cp ambassador:/ambctl ./ambctl
```

Define and setup a [Redis database](docker-compose.yml)
```yml
redis:
    container_name: redis
    image: redis
    net: isolated_nw
```

An [application](docker-compose.yml) that will probe a Redis database
```yml
probe:
    build: .
    command: "node test.js"
    net: isolated_nw
    working_dir: /usr/src/app
```

### Starting up
First, build and run the stack with `docker-compose up`.  This will take
awhile.

Your app will now be complaining it could not reach a Redis node.  Lets fix
this by configuring `ambassador` to point to a Redis endpoint.
```
./ambctl create --name redis --src :6379 --dst redis:6379
```

Now the app works properly.  But why go through all this trouble when you could
just as easily used `redis:6379` in your app?

Suppose I want to talk to a differnt redis node, but I don't want to stop my
app to change code and/or redeploy.  I could reconfigure `ambassador` to do
```
./ambctl cancel --name redis
./ambctl create --name redis --src :6379 --dst remote-redis:6379
```

After a little hiccup, the app starts talking to remote-redis endpoint without
any code change, no deployment required.

