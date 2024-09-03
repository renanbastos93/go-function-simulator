package http

import (
	"context"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type FiberRequest interface {
	Request() any
	AllParams() map[string]string
	Path(override ...string) string
	GetReqHeaders() map[string][]string
	Queries() map[string][]string
	Body() []byte
}

func ConvertHTTPRequestToAPIGatewayProxyRequest(ctx context.Context, request any) events.APIGatewayProxyRequest {
	apiGatewayProxyRequest := events.APIGatewayProxyRequest{}

	switch request := request.(type) {
	case *http.Request:
		if request.Body != nil {
			bodyBytes, err := io.ReadAll(request.Body)
			if err != nil {
				return events.APIGatewayProxyRequest{}
			}
			apiGatewayProxyRequest.Body = string(bodyBytes)
		}
		apiGatewayProxyRequest.Path = request.URL.Path
		apiGatewayProxyRequest.HTTPMethod = request.Method
		apiGatewayProxyRequest.MultiValueHeaders = request.Header
		apiGatewayProxyRequest.MultiValueQueryStringParameters = request.URL.Query()
	case FiberRequest:
		apiGatewayProxyRequest.Body = string(request.Body())
		apiGatewayProxyRequest.Path = request.Path()
		apiGatewayProxyRequest.MultiValueHeaders = request.GetReqHeaders()
		apiGatewayProxyRequest.MultiValueQueryStringParameters = request.Queries()
		apiGatewayProxyRequest.PathParameters = request.AllParams()
	default:
		return events.APIGatewayProxyRequest{}
	}

	return apiGatewayProxyRequest
}
