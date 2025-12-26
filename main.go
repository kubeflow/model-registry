package main

import (
	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/cmd"
)

func main() {
	defer glog.Flush()

	// ADA-KUBEFL-12: commented-out to ease opt-in when doing local development, while not shipping it in the container image ("production code")
	// // start pprof server on 6060
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	cmd.Execute()
}
