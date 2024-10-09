package main

import (
	"log"
	"time"

	"github.com/bclswl0827/heligo"
	"github.com/bclswl0827/mseedio"
)

func main() {
	log.Println("loading miniseed file...")
	var mseed mseedio.MiniSeedData
	err := mseed.Read("testdata.mseed")
	if err != nil {
		panic(err)
	}

	log.Println("creating helicorder context...")
	heli, err := heligo.New(&dataProviderImpl{&mseed}, 24*time.Hour, 30*time.Minute)
	if err != nil {
		panic(err)
	}

	log.Println("drawing plot...")
	err = heli.Plot(mseed.StartTime, 125, 5000, 2.2, 0.5)
	if err != nil {
		panic(err)
	}

	log.Println("saving plot...")
	err = heli.Save(1000, "out.svg")
	if err != nil {
		panic(err)
	}
}
