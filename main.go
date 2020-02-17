package main

import (
	"fmt"
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/config"
	"github.com/RedHatInsights/vmaas-go/app/database"
	"github.com/RedHatInsights/vmaas-go/app/webserver"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	config.SQLiteFilePath = os.Args[1]
	database.Configure()
	var err error
	cache.C = cache.LoadCache()
	PrintMemUsage()

	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()

	f, err = os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	webserver.Run()
	pprof.StopCPUProfile()
	f.Close()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nAlloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\nTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\nSys = %v MiB", bToMb(m.Sys))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
