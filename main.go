package main

import (
	"bufio"
	"fmt"
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/config"
	"github.com/RedHatInsights/vmaas-go/app/database"
	"os"
	"runtime"
)

func main() {
	config.SQLiteFilePath = os.Args[1]
	database.Configure()
	c := cache.LoadCache()
	c.Inspect()
	PrintMemUsage()
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("\nAlloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\nTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\nSys = %v MiB", bToMb(m.Sys))
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
