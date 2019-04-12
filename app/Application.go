package app

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Application struct {
	NumProducers      int
	NumConsumers      int
	workQueue         chan int
	proWaitGroup      sync.WaitGroup
	conWaitGroup      sync.WaitGroup
	producerWorkCount int32
	consumerWorkCount int32
}

const (
	CHANNELSIZE           = 100
	PRODUCER_WORK_PAYLOAD = 20
)

func (a *Application) initProducerRoutines() {
	a.proWaitGroup.Add(a.NumProducers)
	a.spawnProducers()
}

func (a *Application) initConsumerRoutines() {
	a.conWaitGroup.Add(a.NumConsumers)
	a.spawnConsumers()
}

func (a *Application) initApplication() {

	log.Println("Initializing work queue...")
	a.workQueue = make(chan int, CHANNELSIZE)

	log.Println("Initializing producer routines...")
	a.initProducerRoutines()

	log.Println("Initializing consumer routines...")
	a.initConsumerRoutines()

}

func (a *Application) spawnProducers() {
	producerRoutine := func(payload int) {
		for i := 0; i < payload; i++ {
			randInt := rand.Intn(100)
			time.Sleep(time.Duration(randInt) * time.Millisecond)
			a.workQueue <- randInt
			atomic.AddInt32(&a.producerWorkCount, 1)
			totalcount := atomic.LoadInt32(&a.producerWorkCount)
			log.Println(fmt.Sprintf("Pushed %d work time to work-queue. Total items produced = %d", randInt, totalcount))
		}
		a.proWaitGroup.Done()
		log.Println("Producer Done!!!")
	}

	log.Println("Spawning " + fmt.Sprintf("%d", a.NumProducers) + " producers...")
	for i := 0; i < a.NumProducers; i++ {
		go producerRoutine(PRODUCER_WORK_PAYLOAD) //arbitrary number of work requests
	}
}

func (a *Application) spawnConsumers() {
	consumerRoutine := func() {
		for {
			select {
			case work := <-a.workQueue:
				atomic.AddInt32(&a.consumerWorkCount, 1)
				totalcount := atomic.LoadInt32(&a.consumerWorkCount)
				log.Println(fmt.Sprintf("Consuming %d work time. Total work items consumed = %d", work, totalcount))
				time.Sleep(time.Duration(work) * time.Millisecond)
			case <-time.After(1 * time.Second):
				a.conWaitGroup.Done()
				log.Println("Consumer Timed Out!!!")
			}
		}
	}

	log.Println("Spawning " + fmt.Sprintf("%d", a.NumConsumers) + " consumers... ")
	for i := 0; i < a.NumConsumers; i++ {
		go consumerRoutine()
	}
}

func (a *Application) waitForFinish() {
	log.Println("Waiting for pros-cons to finish...")
	a.proWaitGroup.Wait()
	a.conWaitGroup.Wait()
}

func (a *Application) cleanUp() {
	a.workQueue = nil
	log.Println("Done :)")
}

func (a *Application) Run() {
	a.initApplication()
	a.waitForFinish()
	a.cleanUp()
}
