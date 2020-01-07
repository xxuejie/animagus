package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
)

var astFile = flag.String("astFile", "./ast.bin", "AST file to load")
var jsonFile = flag.String("jsonFile", "./ast.json", "JSON file to load")

func main() {
	flag.Parse()

	jsonContent, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var root Root
	err = json.Unmarshal(jsonContent, &root)
	if err != nil {
		log.Fatal(err)
	}

	value, err := root.Value.ToValue()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(proto.MarshalTextString(value))
}
