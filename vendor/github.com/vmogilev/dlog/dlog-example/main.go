package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/vmogilev/dlog"
)

func main() {
	var myVar string
	var debug bool

	flag.StringVar(&myVar, "myVar", "", "Mandatory Var")
	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.Parse()

	if debug {
		dlog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		dlog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	if myVar == "" {
		dlog.Error.Fatalf("myVar is needed! Exiting ...")
	}

	dlog.Info.Printf("thanks - got myVar=%s\n", myVar)
}
