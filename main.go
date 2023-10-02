package main

import (
	"github.com/golang/glog"
	"github.com/opendatahub-io/model-registry/cmd"
)

func main() {
	defer glog.Flush()

	cmd.Execute()
}
