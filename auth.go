package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/vmogilev/dlog"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

func (c *appContext) authDo(ctx context.Context, t string, s *Stack) (bool, error) {
	me := "authDo"
	defer respTime(me)() // don't forget the extra parentheses
	s.Push(me, "<-")

	// tr := &http.Transport{}
	// defer tr.CloseIdleConnections()
	// hc := &http.Client{Transport: tr}
	// defer c.httpTransport.CloseIdleConnections() // this causes even more idle [TIME_WAIT] connections

	et := url.QueryEscape(t)
	s.Push(me, "token raw: "+t)
	s.Push(me, "token enc: "+et)
	dlog.Info.Printf("token raw: %s token enc: %s", t, et)

	req, hcErr := http.NewRequest("GET", c.authEndPoint+et, nil)
	if hcErr != nil {
		s.Push(me, "failed to create http req")
		dlog.Error.Printf("failed to create http req: %s", hcErr)
		return false, hcErr
	}

	// res, doErr := ctxhttp.Do(ctx, hc, req)
	res, doErr := ctxhttp.Do(ctx, c.httpClient, req)
	if doErr != nil {
		s.Push(me, "failed calling ctxhttp.Do")
		dlog.Error.Printf("failed calling ctxhttp.Do: %s", doErr)
		return false, doErr
	}

	defer res.Body.Close()
	b, raErr := ioutil.ReadAll(res.Body)
	if raErr != nil {
		s.Push(me, "failed reading res.Body")
		dlog.Error.Printf("failed reading res.Body: %s", raErr)
		return false, raErr
	}

	ar := string(b)
	if res.StatusCode != http.StatusOK {
		s.Push(me, "HTTP Error: "+strconv.Itoa(res.StatusCode))
		s.Push(me, ar[:strings.Index(ar, "<!DOCTYPE html")])
		s.Push(me, "Err ->")
		return false, fmt.Errorf("Token validation failed due to HTTP error: %d", res.StatusCode)
	}

	if ar == "YES" {
		s.Push(me, "Ok ->")
		return true, nil
	}
	s.Push(me, "No ->")
	dlog.Error.Printf("Token validation failed: %s", ar)
	return false, fmt.Errorf("Token validation failed: %s", ar)

}
