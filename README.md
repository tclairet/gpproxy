# gpproxy

Provide a proxy to `eth_gasPrice` from geth.

# Config

You must specify the node endpoint with the env variable `NODE_URL`. 

By default, the service will be running on port `8545`, you can change that with the env variable `PORT`.

You can also configure the websocket node endpoint with `NODE_URL_WS`, by default it will take the `NODE_URL` and replace the protocol by ws. You don't have to set this env variable if you use infura, the url modification is handled by default.

# Build and run

To build simply exec 
```
cd cmd/gpproxy
go build
```

And to run
```
./gpproxy
```

# Endpoints

There three we of interacting with this `gpproxy`
1. use the endpoint `eth/gasprice`. If a rpc request is used, it will proxify it directly, otherwise it will create a rpc request.
2. use the endpoint `/`. This endpoint proxify every incoming request directly to the node. It can be used as rpc endpoint for a rpc client.
3. use the endpoint `/ws`. Same as `/` but with websockets.

Additionally, two others endpoints are exposed
1. `/healthz` that respond 200 ok when the service is up and running.
2. `/metrics` for prometheus metrics.

# Tests

Tests by default are run against a fake node, but if the env variable `NODE_URL` is set, it will use it