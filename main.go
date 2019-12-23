package main

import (
	"fmt"
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/config"
	"github.com/RedHatInsights/vmaas-go/app/database"
	"github.com/RedHatInsights/vmaas-go/app/webserver"
	"os"
	"runtime"
)

func main() {
	config.SQLiteFilePath = os.Args[1]
	database.Configure()
	cache.C = cache.LoadCache()
	cache.C.Inspect()
	PrintMemUsage()
	webserver.Run()
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
