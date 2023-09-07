package utils

import (
	"log"
	"sync"
)

const numThreads = 5

type Work struct {
	f   func(interface{})
	arg interface{}
}

type ThreadPool struct {
	threads  []chan Work
	queue    chan Work
	mutex    sync.Mutex
	notEmpty *sync.Cond
}

var singleInstance *ThreadPool
var once sync.Once

func GetThreadPoolInstance() *ThreadPool {
	if singleInstance == nil {
		once.Do(
			func() {
				log.Println("Creating Single Thread Pool Instance Now")
				singleInstance = &ThreadPool{}
				singleInstance.threads = make([]chan Work, numThreads)
				singleInstance.queue = make(chan Work, 1)
				singleInstance.mutex = sync.Mutex{}
				singleInstance.notEmpty = sync.NewCond(&singleInstance.mutex)

				for i := 0; i < numThreads; i++ {
					singleInstance.threads[i] = make(chan Work)
					go worker(singleInstance)
				}
			})
	} else {
		log.Println("Single Instance already created")
	}
	return singleInstance
}

func worker(tp *ThreadPool) {
	for {
		tp.mutex.Lock()
		for len(tp.queue) == 0 {
			tp.notEmpty.Wait()
		}

		w := <-tp.queue
		tp.mutex.Unlock()

		w.f(w.arg)
	}
}

func (tp *ThreadPool) ThreadPoolQueue(f func(interface{}), arg interface{}) {
	w := Work{
		f:   f,
		arg: arg,
	}
	tp.mutex.Lock()
	tp.queue <- w
	tp.notEmpty.Signal()
	tp.mutex.Unlock()
}
