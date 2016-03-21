package main

import (
	"time"

	"github.com/vmogilev/dlog"
)

// respTime is initialized from a defer call at the start of the routine
// that you need response time measured for example:
//
//    func (c *appContext) cdnHandler ...{
//       defer respTime("cdnHandler")() // don't forget the extra parentheses
//       ... lots of work ...
//    }
// it returns a closure function that is called at the end of the routine to calculate
// response time (NOW-START time diff) and then send this to Data Dog
//
// see the following links for examples of this technique:
//    1) http://www.gopl.io/ CHAPTER 5 Page# 146 / Chapter 5.8:
//          https://github.com/adonovan/gopl.io/blob/master/ch5/trace/main.go
//    2) Why add "()" after closure body in Golang?
//          http://stackoverflow.com/questions/16008604/why-add-after-closure-body-in-golang
func respTime(what string) func() {
	start := time.Now()
	return func() {
		c.gauge("ales3front."+what+".response_time_ms", float64(time.Since(start)/time.Millisecond), nil, 1)
	}
	//c.gauge("ales3front."+what+".response_time_ms", float64(d/time.Millisecond), nil, 1)
}

func countThis(what string, cnt int64) {
	c.count("ales3front."+what+".times_per_second", cnt, nil, 1)
}

func (c *appContext) simpleEvent(title, text string) {
	if c.ddEnabled {
		if err := c.ddClient.SimpleEvent(title, text); err != nil {
			dogError(err)
		}
	}
}

func dogError(e error) {
	dlog.Error.Printf("Failed to send metric to Data Dog: %s", e.Error())
}

func (c *appContext) gauge(name string, value float64, tags []string, rate float64) {
	if c.ddEnabled {
		if err := c.ddClient.Gauge(name, value, tags, rate); err != nil {
			dogError(err)
		}
	}

}

func (c *appContext) count(name string, value int64, tags []string, rate float64) {
	if c.ddEnabled {
		if err := c.ddClient.Count(name, value, tags, rate); err != nil {
			dogError(err)
		}
	}
}

func (c *appContext) histogram(name string, value float64, tags []string, rate float64) {
	if c.ddEnabled {
		if err := c.ddClient.Histogram(name, value, tags, rate); err != nil {
			dogError(err)
		}
	}
}

func (c *appContext) set(name string, value string, tags []string, rate float64) {
	if c.ddEnabled {
		if err := c.ddClient.Set(name, value, tags, rate); err != nil {
			dogError(err)
		}
	}

}

func (c *appContext) timeInMilliseconds(name string, value float64, tags []string, rate float64) {
	if c.ddEnabled {
		if err := c.ddClient.TimeInMilliseconds(name, value, tags, rate); err != nil {
			dogError(err)
		}
	}

}
