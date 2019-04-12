package patterns

import "time"

func NumberIncrGenerator(incr int, count int) <-chan int {
	ch := make(chan int)
	seed := int(0)
	go func() {
		for counter := 0; counter < count; counter++ {
			seed += incr
			time.Sleep(100 * time.Millisecond)
			ch <- seed
		}
		close(ch)
	}()
	return ch
}

func NumFanout(inCh <-chan int) <-chan int {
	outCh := make(chan int)
	go func() {
		for inEl := range inCh {
			outCh <- inEl
		}
		close(outCh)
	}()
	return outCh
}

func NumSquareFanIn(inCh1 <-chan int, inCh2 <-chan int) <-chan int {
	outCh := make(chan int)

	go func() {
		for {
			select {
			case inEl1, more := <-inCh1:
				if more {
					outCh <- inEl1 * inEl1
				} else {
					inCh1 = nil
				}
			case inEl2, more := <-inCh2:
				if more {
					outCh <- inEl2 * inEl2
				} else {
					inCh2 = nil
				}
			}
			if inCh1 == nil && inCh2 == nil {
				break
			}
		}
		close(outCh)
	}()
	return outCh
}
