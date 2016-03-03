package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/justinas/alice"
	"github.com/vmogilev/dlog"
)

type appContext struct {
	bucket   string
	cred     string
	region   string
	cdn      string
	host     string
	keyID    string
	privKey  *rsa.PrivateKey
	expHours int
}

var c appContext

var (
	awsBucket  = flag.String("awsBucket", "support-pub-dev", "aws bucket name")
	awsCred    = flag.String("awsCred", "ale-s3app", "aws credentials profile from ~/.aws/credentials")
	awsRegion  = flag.String("awsRegion", "us-east-1", "aws region")
	cdnPath    = flag.String("cdnPath", "/cdn/", "URL path prefix to pass to CDN")
	cdnHost    = flag.String("cdnHost", "http://cdn-dev.alcalcs.com/", "CloudFront CDN Hostname and http|https prefix")
	cfKeyID    = flag.String("cfKeyID", "", "CloudFront Signer Key ID")
	cfKeyFile  = flag.String("cfKeyFile", "", "CloudFront Signer Key File Location")
	cfExpHours = flag.Int("cfExpHours", 1, "CloudFront Signed URL Expiration (in hours)")
	httpPort   = flag.String("httpPort", "8080", "HTTP Port")
	debug      = flag.Bool("debug", false, "Debug")
)

func loadKey(f string) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		dlog.Error.Panicf("Failed to read from %s: %s", f, err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		dlog.Error.Panicf("No key found in %s: %s", f)
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

func main() {
	flag.Parse()

	// setup log output streams
	// Trace, Info, Warning, Error
	if *debug {
		dlog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		dlog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	c = appContext{
		bucket:   *awsBucket,
		cred:     *awsCred,
		region:   *awsRegion,
		cdn:      *cdnPath,
		host:     *cdnHost,
		keyID:    *cfKeyID,
		privKey:  loadKey(*cfKeyFile),
		expHours: *cfExpHours,
	}

	middleware := alice.New(logging, recovery)
	http.Handle(*cdnPath, middleware.ThenFunc(c.cdnHandler))
	http.Handle("/", middleware.ThenFunc(c.dispatchHandler))
	dlog.Info.Fatal(http.ListenAndServe(":"+*httpPort, nil))

}
