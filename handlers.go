package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	//"path"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/nu7hatch/gouuid"
	"github.com/vmogilev/dlog"
)

type s3File struct {
	Key           string
	RealKey       string
	ContentLength int64
	ContentType   string
	LastModified  time.Time
}

// Page is used in the template
type Page struct {
	Title     string
	URLPath   string
	RawURL    string
	SignedURL string
	OK        bool
	Error     string
	Meta      *s3File
	Stats     template.HTML
}

func (c *appContext) validateToken(t string) bool {
	defer respTime("validateToken")() // don't forget the extra parentheses

	if t == "" {
		return false
	}

	if t == c.root {
		return true
	}

	return false
}

func (c *appContext) signURL(rawURL string) (url string, mess string, succ bool) {
	ok := true
	message := ""

	signer := sign.NewURLSigner(c.keyID, c.privKey)
	signedURL, err := signer.Sign(rawURL, time.Now().Add(time.Duration(c.expHours)*time.Hour))
	if err != nil {
		dlog.Error.Printf("Failed to sign url, err: %s", err.Error())
		ok = false
		if c.trace {
			message = fmt.Sprintf("Failed to sign url, err: %s", err.Error())
		} else {
			message = "Failed to sign url, please notify site administrator!"
		}
	}
	return signedURL, message, ok
}

// findAndPickOne strips the S3 key's extension and searches S3
// for up to 5 matches using stripped name as a prefix
// this is to support AOS / Diversity Code changes where we have 5 files:
//
//     abc-v1.pdf
//     abc-v2.pdf
//     abc-v3.pdf
//     abc-v4.pdf
//     abc-v5.pdf
//
// and they all belong to the same root filename:
//
//     abc.pdf
//
// on the S3 size we store the root filename under Content-Disposition: filename="abc.pdf"
// for each of the 5 files so that CDN handles the naming properly on the client side
//
// this is WIP -- it's not completed yet, this function should be called from headS3File
func (c *appContext) findAndPickOne(key string, svc *s3.S3) (string, bool) {
	defer respTime("findAndPickOne")() // don't forget the extra parentheses

	root := strings.TrimSuffix(key, filepath.Ext(key))

	params := &s3.ListObjectsInput{
		Bucket:  aws.String(c.bucket), // Required
		MaxKeys: aws.Int64(c.maxDCfiles),
		Prefix:  aws.String(root),
	}
	resp, err := svc.ListObjects(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		dlog.Error.Printf("error getting diversified files for: %s, %s", root, err)
		return "", false
	}

	if len(resp.Contents) == 0 {
		dlog.Error.Printf("found no diversified files for: %s", root)
		return "", false
	}

	rand.Seed(time.Now().UTC().UnixNano())
	k := resp.Contents[rand.Intn(len(resp.Contents))]
	return *k.Key, true

}

func (c *appContext) newS3Svc() *s3.S3 {
	sess := session.New(&aws.Config{
		Region:      aws.String(c.region),
		Credentials: credentials.NewSharedCredentials("", c.cred),
	})
	return s3.New(sess)
}

func (c *appContext) headS3File(key string, rootKey string, svc *s3.S3) (*s3File, string, bool) {
	defer respTime("headS3File")() // don't forget the extra parentheses

	params := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		dKey, ok := c.findAndPickOne(key, svc)
		if ok {
			return c.headS3File(dKey, key, svc)
		}
		dlog.Error.Printf("couldn't get head of file: %s, %s", key, err)
		message := "File is not found!  Please check the URL"
		if c.trace {
			message = fmt.Sprintf("File is not found! Error: %s", err)
		}
		return &s3File{Key: key}, message, false
	}

	var k string
	k = key
	if rootKey != "" {
		k = "*" + rootKey
	}

	f := &s3File{
		Key:           k,
		RealKey:       key,
		ContentLength: *resp.ContentLength,
		ContentType:   *resp.ContentType,
		LastModified:  *resp.LastModified,
	}

	if c.trace {
		dlog.Trace.Println(resp)
		dlog.Trace.Println(f)
	}

	return f, "", true

}

func (c *appContext) dispatchHandler(w http.ResponseWriter, r *http.Request) {
	countThis("dispatchHandler", 1)
	urlpath := r.URL.Path[0:]
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", c.bucket, urlpath)
}

func (c *appContext) cdnHandler(w http.ResponseWriter, r *http.Request) {
	countThis("cdnHandler.hits", 1)
	defer respTime("cdnHandler")() // don't forget the extra parentheses

	t := r.FormValue("t")
	message := ""

	ok := c.validateToken(t)

	if !ok {
		countThis("cdnHandler.badtokens", 1)
		message = fmt.Sprintf("Download Token: %s is invalid", t)
	}

	urlpath := r.URL.Path[len(c.cdn):]
	//rawURL := path.Join(c.host, urlpath) // this strips :// to :/
	rawURL := c.host + urlpath
	signedURL := "#"
	meta := &s3File{}

	if ok {
		meta, message, ok = c.headS3File(urlpath, "", c.newS3Svc())
	}

	if ok {
		rawURL = c.host + meta.RealKey
		signedURL, message, ok = c.signURL(rawURL)
	}

	tokens := []string{
		"<!-- ",
		"SignedURL: " + signedURL,
		"URLPath: " + urlpath,
		"RawURL: " + rawURL,
		"-->",
	}

	stats := strings.Join(tokens, "\n")

	p := &Page{
		Title:     "Download",
		URLPath:   urlpath,
		RawURL:    rawURL,
		SignedURL: signedURL,
		OK:        ok,
		Error:     message,
		Meta:      meta,
		Stats:     template.HTML(stats),
	}
	c.renderTemplate(w, "download", p)
	//fmt.Fprintf(w, "<h1>%s</h1><div>urlpath: %s</div><div>rawURL: %s</div><div>signedURL: %s</div>", c.cdn, urlpath, rawURL, signedURL)

}
