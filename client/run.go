package main

import (
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

	startFreePlay(debugSet)
}
