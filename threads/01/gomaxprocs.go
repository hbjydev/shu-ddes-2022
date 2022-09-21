package main

import (
	"log"
	"runtime"
)

func main() {
    max := runtime.GOMAXPROCS(0)
    log.Printf("GOMAXPROCS is set to %v", max)
}
