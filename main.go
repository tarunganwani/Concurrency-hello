package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tarunganwani/Concurrency-Hello/patterns"
)

func main() {
	// app := app.Application{
	// 	NumConsumers: 10,
	// 	NumProducers: 10,
	// }
	// app.Run()

	fanin_fanout_flag := true
	if fanin_fanout_flag {
		genChannel := patterns.NumberIncrGenerator(1, 200)

		genSplitChannel1 := patterns.NumFanout(genChannel)
		genSplitChannel2 := patterns.NumFanout(genChannel)
		genSquaresChannel := patterns.NumSquareFanIn(genSplitChannel1, genSplitChannel2)

		for squaredNum := range genSquaresChannel {
			log.Println("Number recd on squares channel = ", squaredNum)
		}
	}

	generatorFlag := false // true
	if generatorFlag == true {
		timeRetChan := patterns.QueryWorkerGeneratorFn()
		for {
			log.Println("Waiting on time channel...")
			select {
			case timeRet := <-timeRetChan:
				log.Println(fmt.Sprintf("Query finished in %d milli seconds", timeRet))
			case <-time.After(2 * time.Second):
				log.Println("Timeout!")
			}
		}
	}

	rangeOverChannelExampleFlag := false
	if rangeOverChannelExampleFlag {
		rangeOverChannelExample := func() {

			log.Println("Begin...")
			fibCh := make(chan int, 10)
			first := 0
			second := 1
			fibGen := func() int {
				second, first = first+second, second
				log.Println("producing fib ", second)
				return second
			}

			var wg sync.WaitGroup
			wg.Add(2)
			//Put fibonacci sequence over the channel for 300 milli-seconds
			go func() {

				log.Println("Put fibonacci sequence over the channel for 300 milli-seconds..")
				for i := 0; i < 100; i++ {
					fibCh <- fibGen()
					// select {
					// case <-time.After(1 * time.Nanosecond):
					// 	close(fibCh)
					// 	wg.Done()
					// 	log.Println("** DONE pushing fib")
					// 	return
					// case fibCh <- fibGen():
					// 	log.Println("pushing fib")
					// }
				}
				wg.Done()
				close(fibCh)
			}()

			go func() {
				log.Println("Get fibonacci sequence from the channel..")
				count := 1
				for {
					fib, more := <-fibCh
					if more {
						log.Println("fib ", count, " value ", fib)
						count++
					} else {
						log.Println("Done!", fib)
						wg.Done()
						return
					}
				}
			}()
			wg.Wait()
		}
		rangeOverChannelExample()
	}

}
