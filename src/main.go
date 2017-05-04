package main

/**
 *
 * Nicholas Sardo
 * (c)2017
 *
 * Experiment with naive caching
 */

import (
	"fmt"
	"math/rand"
	"time"
)


/* MAGIC NUMBERS */
const (
	WEIGHT    int     = 2
	MA        float32 = 2.0
)

var LENGTH int = 0

/* DATA STRUCTURE FOR CALLS */
type cacheCall struct {
	id              string
	weight          int
	duration_sample []int
	moving_avg      float32
}

var cache = make(map[string]*cacheCall)

/**
 * Sum Array or Slice
 */
func sum(s []int) float32 {
	var sm float32 = 0
	l := len(s)

	for i := 0; i < l; i++ {
		sm += float32(s[i])
	}
	return sm
}

/*
 *	SIMULATE A CALL WITH VARYING DEGREE'S OF LATENCY
 */
func (c *cacheCall) randomLatencyCall() {
	//SEED RANDOM NUMBER GEN
	s1 := rand.NewSource(time.Now().UnixNano())

	//CREATE
	r1 := rand.New(s1)

	//GRAB ONE BETWEEN 1 AND 10
	irand := r1.Intn(10) + 1

	fmt.Printf("delay is %d seconds\n", irand)

	//BUMP FREQUENCY OF THIS CALL
	c.weight++

	//STORE TIME IT TOOK TO RESOLVE THIS FUNCTION CALL
	c.duration_sample = append(c.duration_sample, irand)

	//CREATE MOVING AVERAGE OF CALL TIME'S FOR THIS CALL
	moving_average := sum(c.duration_sample) / float32(len(c.duration_sample))

	c.moving_avg = moving_average

	fmt.Printf("ma = %f\n", moving_average)

	if  LENGTH >= 4 {
		pruneCache()
	}

	if _, ok := cache[c.id]; ok { //SEND IF CACHED
		cache[c.id] = c
		fmt.Printf("Call id: %s is cached, weight is %d, durations: %v, ma: %f\n",
								c.id, cache[c.id].weight, cache[c.id].duration_sample, cache[c.id].moving_avg )
		return
	} else {                      //OR CACHE
		cache[c.id] = c
		LENGTH++
	}

	//DELAY FOR irand SECONDS TO SIMULATE LATENCY
	//timer1 := time.NewTimer(time.Second * time.Duration(irand))
	//<-timer1.C

	//fmt.Println("Timer 1 expired\n")
}


func pruneCache() {
	for k, _ := range cache {
		//IF WEIGHT AND MOVING AVERAGE ARE AT PRESCRIBED LEVEL, KEEP IT
		if cache[k].weight >= WEIGHT && cache[k].moving_avg >= MA {
			//KEEP
		} else {
			delete( cache, k )
			LENGTH--
		}
	}
}


func main() {
	//SEED RANDOM NUMBER GEN
	s1 := rand.NewSource(time.Now().UnixNano())
	//CREATE
	r1 := rand.New(s1)

	//GET START TIME
	now := time.Now()

	//CREATE INSTANCES OF SOME CALLS
	call1 := &cacheCall{"one", 0,make([]int,0), 0 }
	call2 := &cacheCall{"two", 0,make([]int,0), 0 }
	call3 := &cacheCall{"three", 0,make([]int,0), 0 }
	call4 := &cacheCall{"four", 0,make([]int,0), 0 }
	call5 := &cacheCall{"five", 0,make([]int,0), 0 }

	//RUN THE SIMULATION
	for i := 1; i < 20; i++ {
		irand := r1.Intn(5) + 1
		switch  irand {
		case 1:
			call1.randomLatencyCall()
			break
		case 2:
			call2.randomLatencyCall()
			break
		case 3:
			call3.randomLatencyCall()
			break
		case 4:
			call4.randomLatencyCall()
			break
		case 5:
			call5.randomLatencyCall()
			break
		}
	}
	fmt.Printf("Total Elapsed Time: %v", time.Now().Sub(now))

}
