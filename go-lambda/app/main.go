package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var chiLambda *chiadapter.ChiLambda

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	router := chi.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Print(err)
			return
		}
		_ = render.Render(w, r, &apiResponse{
			Status:      http.StatusOK,
			URL:         r.URL.String(),
			RequestBody: "{}",
		})
	})
	// lambda.StartWithContext(context.Background(), handler)
	lambda.StartWithOptions(handler, lambda.WithContext(context.Background()))
}

type apiResponse struct {
	Status      int    `json:"status_code,omitempty"`
	URL         string `json:"url,omitempty"`
	RequestBody string `json:"request_body,omitempty"`
}

func (a apiResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, a.Status)
	return nil
}
