## Setting up a working deployment
Lets go over a typical setup using [docker-compose](https://docs.docker.com/compose/)

You will need [docker](https://www.docker.com/) installed.

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

Define and setup [Redis](docker-compose.yml)
```yml
redis:
    container_name: redis
    image: redis
    net: isolated_nw
```

An [application](docker-compose.yml) that will probe Redis
```yml
probe:
    build: .
    command: "node test.js"
    net: isolated_nw
    working_dir: /usr/src/app
```

### Starting up
Launch the stack with `docker-compose up`.  This will take a while.

Your app will now be complaining it could not reach a Redis node.  Lets fix
this by configuring `ambassador` to point to a Redis endpoint.
```
./ambctl create --name redis --src :6379 --dst redis:6379
```

But why go through all this trouble when you could  just as easily specify
`redis:6379` to connect in your app?  Suppose I want to talk to a different
Redis node, but I don't want to stop my app to change code and/or redeploy.
I could configure `ambd` by
```
./ambctl cancel --name redis
./ambctl create --name redis --src :6379 --dst remote-redis:6379
```

After a little hiccup, the app starts talking to remote-redis endpoint without
any code change, no deployment required.

### Advanced Configuration

