package main

import (
	"errors"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

func (c *appContext) authDo(ctx context.Context, t string, s *stack) (bool, error) {
	me := "authDo"
	defer respTime(me)() // don't forget the extra parentheses
	s.Push(me, "<-")

	x := make(chan error, 1)
	go func() {
		x <- func() error {
			time.Sleep(1000 * time.Millisecond)
			return errors.New("external auth is not yet implemented")
		}()
	}()

	select {
	case <-ctx.Done():
		s.Push(me, "timeout ended: "+strconv.FormatInt(c.authTimeout, 10))
		s.Push(me, "waiting for 1000 ms to end and x chan to close ...")
		<-x // Wait for sleep 1000 to return
		return false, ctx.Err()
	case err := <-x:
		return false, err
	}

	// return false, errors.New("external auth is not yet implemented")
}
