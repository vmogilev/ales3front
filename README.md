# Alcatel AWS Web Microservice

Signs AWS CloundFront URLs to serve files from S3

    Usage of ales3front:
      -awsBucket string
            aws bucket name (default "support-pub-dev")
      -awsCred string
            aws credentials profile from ~/.aws/credentials (default "ale-s3app")
      -awsRegion string
            aws region (default "us-east-1")
      -cdnHost string
            CloudFront CDN Hostname and http|https prefix (default "http://cdn-dev.alcalcs.com/")
      -cdnPath string
            URL path prefix to pass to CDN (default "/cdn/")
      -cfExpHours int
            CloudFront Signed URL Expiration (in hours) (default 1)
      -cfKeyFile string
            CloudFront Signer Key File Location
      -cfKeyID string
            CloudFront Signer Key ID
      -debug
            Debug
      -htmlPath string
            absolute or relative path to html templates (default "./html")
      -httpPort string
            HTTP Port (default "8080")
      -rootToken string
            With this token any download is allowed (default "gTxHrJ")

REV:  # Wed: March 9 2016
