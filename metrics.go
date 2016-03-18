package main

import (
	"time"

	"github.com/vmogilev/dlog"
)

func respTime(what string, d time.Duration) {
	c.gauge("app."+what+".response_time", float64(d/time.Millisecond), nil, 1)
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

func (c *appContext) simpleEvent(title, text string) {
	if c.ddEnabled {
		if err := c.ddClient.SimpleEvent(title, text); err != nil {
			dogError(err)
		}
	}
}
