package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xxuejie/animagus/pkg/generic"
	"github.com/xxuejie/animagus/pkg/indexer"
	"google.golang.org/grpc"
)

var astFile = flag.String("astFile", "./ast.bin", "AST file to load")
var redisUrl = flag.String("redisUrl", "redis://127.0.0.1:6379", "Redis URL")
var rpcUrl = flag.String("rpcUrl", "http://127.0.0.1:8114", "CKB RPC URL")
var grpcListenAddress = flag.String("grpcListenAddress", ":4000", "GRPC Listen Address")

func main() {
	flag.Parse()

	astContent, err := ioutil.ReadFile(*astFile)
	if err != nil {
		log.Fatal(err)
	}
	redisPool := &redis.Pool{
		MaxIdle:     2,
		IdleTimeout: 60 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.DialURL(*redisUrl) },
	}
	// TODO: multiple call support later
	i, err := indexer.NewIndexer(astContent, redisPool, *rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	genericServer, err := generic.NewServer(astContent, redisPool, *rpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", *grpcListenAddress)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	generic.RegisterGenericServiceServer(grpcServer, genericServer)

	go func() {
		log.Fatal(i.Run())
	}()

	grpcServer.Serve(lis)
}
