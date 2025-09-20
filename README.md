# NF-example

This repository is a sample NF for a simple HTTP service!
Try adding more services and learn how to collaborate using GitHub.

## Compile & Run

```sh
make
./bin/nf -c config/nfcfg.yaml
```

## Try Service

```sh
> curl -X GET http://127.0.0.163:8000/default/
"Hello free5GC!"

> curl -X GET http://127.0.0.163:8000/spyfamily/
"Hello SPYxFAMILY!"

> curl -X GET http://127.0.0.163:8000/spyfamily/character/Loid
"Character: Loid Forger"
```

## Go Test

```sh
> go test -v ./...
```
