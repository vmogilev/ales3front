package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/justinas/alice"
	"github.com/vmogilev/dlog"
)

type appContext struct {
	bucket        string
	cred          string
	region        string
	cdn           string
	host          string
	keyID         string
	privKey       *rsa.PrivateKey
	expHours      int
	htmlPath      string
	trace         bool
	root          string
	ddClient      *statsd.Client
	ddEnabled     bool
	maxDCfiles    int64
	authTimeout   int64
	authEndPoint  string
	httpTransport *http.Transport
	httpClient    *http.Client
}

var c appContext

var (
	awsBucket    = flag.String("awsBucket", "support-pub-dev", "aws bucket name")
	awsCred      = flag.String("awsCred", "ale-s3app", "aws credentials profile from ~/.aws/credentials")
	awsRegion    = flag.String("awsRegion", "us-east-1", "aws region")
	htmlPath     = flag.String("htmlPath", "./html", "absolute or relative path to html templates")
	cdnPath      = flag.String("cdnPath", "/cdn/", "URL path prefix to pass to CDN")
	cdnHost      = flag.String("cdnHost", "http://cdn-dev.alcalcs.com/", "CloudFront CDN Hostname and http|https prefix")
	cfKeyID      = flag.String("cfKeyID", "", "CloudFront Signer Key ID")
	cfKeyFile    = flag.String("cfKeyFile", "", "CloudFront Signer Key File Location")
	cfExpHours   = flag.Int("cfExpHours", 1, "CloudFront Signed URL Expiration (in hours)")
	httpPort     = flag.String("httpPort", "8080", "HTTP Port")
	trace        = flag.Bool("trace", false, "Trace")
	rootToken    = flag.String("rootToken", "gTxHrJ", "With this token any download allowed")
	ddAgent      = flag.String("ddAgent", "", "host:port of the Data Dog DogStatsD Agent, if null - no stats are sent")
	ddPrefix     = flag.String("ddPrefix", "ales3front", "Data Dog namespace prefix (no dot) - added to all metrics")
	maxDCfiles   = flag.Int64("maxDCfiles", 5, "Maximum number of diversified download files to search for")
	authTimeout  = flag.Int64("authTimeout", 600, "Auth token validation timeout in Milliseconds")
	authEndPoint = flag.String("authEndPoint", "", "Auth callback end point at Calabasas; EX: https://support.esd.alcatel-lucent.com/pm/cdlv?t=")
)

func loadKey(f string) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		dlog.Error.Panicf("Failed to read from %s: %s", f, err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		dlog.Error.Panicf("No key found in: %s", f)
	}

	var k *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		k, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			dlog.Error.Panicf("Failed to parse key %s: %s", f, err)
		}
	default:
		dlog.Error.Panicf("ssh: unsupported key type %q", block.Type)
	}
	return k

}

// callDog sets up a connection to DataDog agent if agent is provided
// and returns a pointer to it along with true which is then set at
// application level so that each routine can eval if metrics are enabled
func callDog(agent string, prefix string, tag string) (*statsd.Client, bool) {
	if agent == "" {
		return &statsd.Client{}, false
	}
	c, err := statsd.New(agent)
	if err != nil {
		dlog.Error.Panicf("Failed to dial Data Dog Agent on %s, %s:", agent, err)
	}
	c.Namespace = prefix + "."
	c.Tags = append(c.Tags, tag)
	return c, true
}

func main() {
	flag.Parse()

	// setup log output streams
	// Trace, Info, Warning, Error
	if *trace {
		dlog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		dlog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	ddClient, ddEnabled := callDog(*ddAgent, *ddPrefix, *awsRegion)

	//tr := &http.Transport{DisableKeepAlives: true}
	tr := &http.Transport{MaxIdleConnsPerHost: 10}
	defer tr.CloseIdleConnections()
	httpClient := &http.Client{Transport: tr}

	c = appContext{
		bucket:        *awsBucket,
		cred:          *awsCred,
		region:        *awsRegion,
		cdn:           *cdnPath,
		host:          *cdnHost,
		keyID:         *cfKeyID,
		privKey:       loadKey(*cfKeyFile),
		expHours:      *cfExpHours,
		htmlPath:      *htmlPath,
		trace:         *trace,
		root:          *rootToken,
		ddClient:      ddClient,
		ddEnabled:     ddEnabled,
		maxDCfiles:    *maxDCfiles,
		authTimeout:   *authTimeout,
		authEndPoint:  *authEndPoint,
		httpTransport: tr,
		httpClient:    httpClient,
	}

	// make sure we close Data Dog connection on exit
	defer func(pref string, app appContext) {
		if app.ddEnabled {
			app.simpleEvent(*ddPrefix+".ales3front", "stopping")
			app.ddClient.Close()
		}
	}(*ddPrefix, c)

	/*
		oschan := make(chan os.Signal, 1)
		signal.Notify(oschan, os.Interrupt)
		go func(pref string, app appContext) {
			for range oschan {
				defer func() {
					if app.ddEnabled {
						app.simpleEvent(*ddPrefix+".ales3front", "stopping / CTL-C")
						app.ddClient.Close()
					}
					os.Exit(0)
				}()
			}
		}(*ddPrefix, c)
	*/
	c.simpleEvent(*ddPrefix+".ales3front", "starting")
	middleware := alice.New(logging, recovery)
	http.Handle(*cdnPath, middleware.ThenFunc(c.cdnHandler))
	http.Handle("/", middleware.ThenFunc(c.dispatchHandler))
	dlog.Error.Fatal(http.ListenAndServe(":"+*httpPort, nil))

}
