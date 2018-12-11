package main

import (
	"github.com/faiface/pixel/pixelgl"
	"log"
	"os"
	"runtime/pprof"
)

func main() {

	_, debugSet := os.LookupEnv("sq39_debug")
	if debugSet {
		cpuProfile := "cpu.txt"
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)

			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	pixelgl.Run(run)
}
