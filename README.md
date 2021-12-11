# Envoy Proxy Apache Cassandra filter

This is a playground for my experiments with Envoy filter and Apache Cassandra CQL binary protocol. The filter is implemented on top of proxy-wasm-go-sdk and
the CQL protocol is based on Datastax `go-Cassandra-native-protocol`. I love Envoy cause is modular and its filter can be used for different purposes for examples

* Service discovery
* Dynamic forward proxy
* Rate limiting、
* Circuit breaker
* Multi-tenancy
* SQL audit
* Session statistics
* Queries metrics 
* etc etc

After some experiments I would like to bring filter in C++ native mode and use datastax/cpp-driver as codec, `At the moment of writing it exports only queries metrics in tiny mode in Prometheus format

## Build

The Go SDK for Proxy-Wasm is powered by TinyGo, then before compile project please check it out [tinygo](https://tinygo.org/)

```sh
make build
```

## Run
```sh
docker-compose -f ./dev-env/docker-compose.yaml up -d
envoy -c envoy.yaml --log-level debug
```
populate db with `./dev-env/data/cass-populate-tables.sh` and then follow commands in `./dev-env/data/verify.txt`

```sh
curl http://localhost:8001/stats/prometheus | grep -i statements_select 
```


## License

This project is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.

## Acknowledgments

- Thanks [envoy proxy](https://www.envoyproxy.io/) The great project.
- Thanks [proxy-wasm-go-sdk](https://github.com/tetratelabs/proxy-wasm-go-sdk) for their sdk.
- Thanks [go-cassandra-native-protocol](https://github.com/datastax/go-cassandra-native-protocol). It Saves me a ton of work and avoid me to make a lot of mistakes.