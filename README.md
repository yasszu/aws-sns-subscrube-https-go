# Amazon SNS Middleware for Go

This middleware is used fo subscribing an HTTP/S endpoint to an Amazon SNS topic.  

Features:
- Confirms the subscription
- Verify the signature of SNS messages

## Installation
Install package:
```
go get -u github.com/yasszu/aws-sns-subscrube-https-go
```

## Example
```go
package main

import (
	"fmt"
	"net/http"

	sns "github.com/yasszu/aws-sns-subscrube-https-go"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	topicArn := "arn:aws:sns:us-west-2:123456789012:MyTopic"
	middleware := sns.NewMiddleware()
	http.HandleFunc("/", middleware.Subscribe(topicArn)(handler))

	http.ListenAndServe(":8080", nil)
}
```
