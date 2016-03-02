package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/vmogilev/dlog"
)

type AppConfig struct {
	AppRoot  string
	HttpFQDN string
}

var app AppConfig

var (
	bucket    = flag.String("bucket", "support-pub-dev", "aws bucket name")
	awscred   = flag.String("awsCred", "ale-s3app", "aws credentials profile from ~/.aws/credentials")
	awsregion = flag.String("awsRegion", "us-east-1", "aws region")

	rootDir = flag.String("rootDir", "./", "Root Directory [where the ./conf, ./static, ./templates and ./img dirs are]")

	httpHost    = flag.String("httpHost", "http://localhost", "HTTP Host Name")
	httpPort    = flag.String("httpPort", "8080", "HTTP Port")
	httpMount   = flag.String("httpMount", "/uploads", "HTTP Mount Point")
	dlPrefix    = flag.String("dlPrefix", "/release", "files under this URL path will be handed off to CloundFront CDN")
	httpHostExt = flag.String("httpHostExt", "", "Fully Qualified External Path if using Proxy [EX: http://mydomain.com/path]")

	debug = flag.Bool("debug", false, "Debug")
)

func main() {
	flag.Parse()

	// setup log output streams
	if *debug {
		dlog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		dlog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	var fqdn string
	if httpHostExt != "" {
		fqdn = httpHostExt
	} else {
		mp := MountPoint(httpMount)
		if httpPort == "80" {
			fqdn = httpHost + mp
		} else {
			fqdn = httpHost + ":" + httpPort + mp
		}
	}

	app = AppConfig{
		AppRoot:  *rootDir,
		HttpFQDN: fqdn,
	}

	router := NewRouter(httpMount, rootDir, filepath.Join(*rootDir, "img"))
	dlog.Info.Fatal(http.ListenAndServe(":"+httpPort, router))

}
