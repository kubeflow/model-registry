package main

import (
	"github.com/golang/glog"
	"github.com/opendatahub-io/model-registry/cmd"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	defer glog.Flush()

	// start pprof server on 6060
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	cmd.Execute()
}
