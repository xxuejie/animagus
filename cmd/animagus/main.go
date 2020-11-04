package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"

	"github.com/xxuejie/animagus/pkg/generic"
	"github.com/xxuejie/animagus/pkg/indexer"
	"github.com/xxuejie/animagus/pkg/store"
	"google.golang.org/grpc"
)

var astFile = flag.String("astFile", "./ast.bin", "AST file to load")

var dataDir = flag.String("dataDir", "./badgerdb", "DB dir")
var rpcUrl = flag.String("rpcUrl", "http://127.0.0.1:8114", "CKB RPC URL")
var grpcListenAddress = flag.String("grpcListenAddress", ":4000", "GRPC Listen Address")

func main() {
	flag.Parse()

	astContent, err := ioutil.ReadFile(*astFile)
	if err != nil {
		log.Fatal(err)
	}

	storeClient := store.NewClient(*dataDir)
	err = storeClient.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer storeClient.Close()

	// TODO: multiple call support later
	i, err := indexer.NewIndexer(astContent, storeClient, *rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	genericServer, err := generic.NewServer(astContent, storeClient, *rpcUrl)
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
