package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("no ip in http response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("non 200 response found")
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayV2HTTPResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayV2HTTPResponse{}, ErrNoIP
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "text/html"
	return events.APIGatewayV2HTTPResponse{
		Body:       fmt.Sprintf("<h1>Hello, %v</h1>", string(ip)),
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
