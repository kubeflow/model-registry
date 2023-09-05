//go:build tools
// +build tools

package tools

import (
        _ "google.golang.org/protobuf/cmd/protoc-gen-go"
        _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)
