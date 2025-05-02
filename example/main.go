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
	heli, err := heligo.New(&dataProviderImpl{&mseed}, 24*time.Hour, 15*time.Minute)
	if err != nil {
		panic(err)
	}

	log.Println("drawing plot...")
	err = heli.Plot(mseed.StartTime, 10000, 500, 1, nil) // Set colorScheme to nil to use default color scheme
	if err != nil {
		panic(err)
	}

	log.Println("saving plot...")
	err = heli.Save(1000, "out.svg")
	if err != nil {
		panic(err)
	}

	// Get plot data in bytes
	bytes, err := heli.Bytes(1000, "png")
	if err != nil {
		panic(err)
	}
	log.Println("plot data in bytes:", len(bytes))

}
