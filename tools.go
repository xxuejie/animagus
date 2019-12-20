// Inspired from https://marcofranssen.nl/manage-go-tools-via-go-modules/

// +build tools

package main

import (
	_ "github.com/awalterschulze/goderive"
	_ "github.com/golang/protobuf/protoc-gen-go"
)
