package http

import (
	"context"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
)

type FiberCtx = *fiber.Ctx
type Request = *http.Request

func ConvertHTTPRequestToAPIGatewayProxyRequest(ctx context.Context, request any) events.APIGatewayProxyRequest {
	if request == nil {
		panic("request is nil")
	}

	apiGatewayProxyRequest := events.APIGatewayProxyRequest{}
	switch request := request.(type) {
	case Request:
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
		apiGatewayProxyRequest.PathParameters = mux.Vars(request)

		return apiGatewayProxyRequest

	case FiberCtx:
		apiGatewayProxyRequest.PathParameters = request.AllParams()
		apiGatewayProxyRequest.Path = request.Path()
		apiGatewayProxyRequest.HTTPMethod = request.Method()
		apiGatewayProxyRequest.MultiValueHeaders = request.GetReqHeaders()
		apiGatewayProxyRequest.MultiValueQueryStringParameters = make(map[string][]string)
		request.Context().QueryArgs().VisitAll(func(key []byte, value []byte) {
			if v, found := apiGatewayProxyRequest.MultiValueQueryStringParameters[string(key)]; found {
				apiGatewayProxyRequest.MultiValueQueryStringParameters[string(key)] = append(v, string(value))
			} else {
				apiGatewayProxyRequest.MultiValueQueryStringParameters[string(key)] = []string{string(value)}
			}
		})

		apiGatewayProxyRequest.Body = string(request.Body())

		return apiGatewayProxyRequest
	}

	return apiGatewayProxyRequest
}
