package main

import (
	"github.com/gobuffalo/packr"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func run() {
	rand.Seed(time.Now().UnixNano())
	_, debugSet := os.LookupEnv(DebugEnvVar)
	if debugSet {
		cpuProfile := CpuProfile
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// ensure assets from this module are packed to the binary
	packr.NewBox("./assets")

	startFreePlay(debugSet)
}
