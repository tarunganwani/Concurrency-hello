package patterns

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func QueryWorkerGeneratorFn() <-chan int {

	retChan := make(chan int)

	go func() {
		//Do actual work
		for {
			timeToQuery := rand.Intn(3000)
			time.Sleep(time.Duration(timeToQuery) * time.Millisecond)
			retChan <- timeToQuery
			log.Println(fmt.Sprintf("Pushed %d", timeToQuery))
		}
	}()

	return retChan
}
