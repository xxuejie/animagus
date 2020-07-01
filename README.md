animagus
========

Note the code in this repository is still under development and NOT production ready. You should get familiar with CKB and know what you are doing before using it.

Animagus is a data processor for Nervos CKB. With a unified set of operations, it is designed to aid CKB in the following(but not limited to) tasks:

* Indexing common cells
* Aggregate balances
* Signal special transactions, such as Nervos DAO deposits
* Assemble transactions
* Generate on-chain smart contracts

Animagus works by defining an AST(abstract syntax tree) first, the AST can be generated via any [protobuf](https://github.com/protocolbuffers/protobuf) enabled languages. It's also possible to build small special purpose languages that compile directly to animagus' AST. Then with the AST, animagus can automatically scan CKB, index certain cells, and response to user request via [grpc](https://grpc.io/).

The name animagus comes from [Harry Potter series](https://harrypotter.fandom.com/wiki/Animagus). We are hoping with a well designed AST, animagus can help achieve most, if not all requirements needed when building wonderful dapps and a better ecosystem on CKB.

Note this project is still in fast development phase, the AST used is not yet in a stable state. We might modify it to add more capabilities.

We will continue to add docs here explaining how this project works once the AST is in a relatively more stable state, but if you want to get a taste right now, here's the [slide](https://github.com/xxuejie/animagus/blob/develop/docs/A%20new%20dapp%20framework.pdf) for an introductory talk on animagus.

# How to Run

We've packed a small [example](https://github.com/xxuejie/animagus/tree/develop/examples/balance) that can aggregate balances of different accounts in CKB. This can help overcome some of the challenges brought by CKB's unique flexibility.

Animagus requires the following dependencies:

* [CKB](https://github.com/nervosnetwork/ckb)
* [ckb-graphql-server](https://github.com/xxuejie/ckb-graphql-server)
* [Redis](https://redis.io/)

Animagus is compatible with latest lina mainnet version of CKB. Below I'm building and launching CKB from source, but you can also use precompiled binaries:

```
$ mkdir -p /tmp/animagus-demo
$ cd /tmp/animagus-demo
$ git clone https://github.com/nervosnetwork/ckb
$ cd ckb
$ cargo build --release
$ target/release/ckb init -C mainnet -c mainnet
$ target/release/ckb run -C mainnet
```

Now we can run the GraphQL server:

```
$ mkdir -p /tmp/animagus-demo
$ cd /tmp/animagus-demo
$ git clone https://github.com/xxuejie/ckb-graphql-server
$ cd ckb-graphql-server
$ cargo build --release
$ target/release/ckb-graphql-server --db ../ckb/mainnet/data/db --listen 0.0.0.0:3001
```

I'm using docker to quickly start that a temporary Redis server, but you can also using other ways to launch Redis:

```
$ docker run -d --rm -p 6379:6379 --name animagus-redis redis:alpine
```

First thing we need to do, is to generate an AST dump file for animagus, a [sample](https://github.com/xxuejie/animagus/blob/develop/examples/balance/generate_ast.go) has been prepared for this purpose:

```
$ mkdir -p /tmp/animagus-demo
$ cd /tmp/animagus-demo
$ git clone https://github.com/xxuejie/animagus
$ cd animagus/examples/balance
$ go run generate_ast.go
```

You will noticed a new file `balance.bin` has been generated in `animagus/examples/balance` folder. Now we can start animagus:

```
$ cd /tmp/animagus-demo/animagus
$ make install-tools
$ go build ./cmd/animagus
$ ./animagus -astFile=./examples/balance/balance.bin
```

Notice if you use different ports for GraphQL server and Redis, you might need to tweak animagus start flags, see `./animagus --help` for details

You will notice logs since animagus is indexing cells. We have prepared a small [file](https://github.com/xxuejie/animagus/blob/develop/examples/balance/call_balance.rb) that you can use to check balances. Given the `args` part in a lock script, this file queries against animagus for the current balance of that account:

```
$ cd /tmp/animagus-demo/animagus/examples/balance
$ ruby call_balance.rb ba03db27e31d19ebc4fda56b440fb92310d64d0e
```
