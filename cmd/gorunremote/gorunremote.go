package main

import (
	"sync"

	"github.com/karlockhart/gorunremote/pkg/api"
	"github.com/karlockhart/gorunremote/pkg/config"
)

func main() {
	config.LoadConfig()
	a := api.NewGoRunRemoteApi()
	var wg sync.WaitGroup
	wg.Add(1)
	go a.Start(&wg)
	wg.Wait()
}
