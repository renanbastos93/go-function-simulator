package http_test

import (
	"bytes"
	"context"
	"io"
	httpStd "net/http"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/renanbastos93/go-function-simulator/pkg/http"
	"github.com/stretchr/testify/assert"
)

func TestConvertHTTPRequestToAPIGatewayProxyRequestNative(t *testing.T) {
	ctx := context.Background()
	t.Run("should be convert *http.Request to events.APIGatewayProxyRequest", func(t *testing.T) {
		request := &httpStd.Request{}
		request.Method = "GET"
		request.URL = &url.URL{}
		request.URL.Path = "/path"
		request.Header = httpStd.Header{}
		request.Header.Set("key", "value")
		request.Body = httpStd.NoBody

		apiGatewayProxyRequest := http.ConvertHTTPRequestToAPIGatewayProxyRequest(ctx, request)
		assert.NotNil(t, apiGatewayProxyRequest)
		assert.Equal(t, "GET", apiGatewayProxyRequest.HTTPMethod)
		assert.Equal(t, "/path", apiGatewayProxyRequest.Path)
		assert.Equal(t, "", apiGatewayProxyRequest.Body)
		assert.Equal(t, 0, len(apiGatewayProxyRequest.MultiValueQueryStringParameters))
		assert.Equal(t, 0, len(apiGatewayProxyRequest.PathParameters))
		assert.Equal(t, 1, len(apiGatewayProxyRequest.MultiValueHeaders))
		assert.Equal(t, []string(nil), apiGatewayProxyRequest.MultiValueHeaders["key"])
	})

	t.Run("should be convert fiber.Request to events.APIGatewayProxyRequest", func(t *testing.T) {
		app := fiber.New()
		app.Get("/fiber/path/:id", func(c *fiber.Ctx) error {
			apiGatewayProxyRequest := http.ConvertHTTPRequestToAPIGatewayProxyRequest(c.Context(), c)
			assert.NotNil(t, apiGatewayProxyRequest)
			assert.Equal(t, "GET", apiGatewayProxyRequest.HTTPMethod)
			assert.Equal(t, "/fiber/path/123", apiGatewayProxyRequest.Path)
			assert.Equal(t, "123", apiGatewayProxyRequest.PathParameters["id"])
			assert.Equal(t, "", apiGatewayProxyRequest.Body)
			assert.Equal(t, 2, len(apiGatewayProxyRequest.MultiValueQueryStringParameters))
			assert.Equal(t, 1, len(apiGatewayProxyRequest.PathParameters))
			assert.Equal(t, 2, len(apiGatewayProxyRequest.MultiValueHeaders))
			assert.Equal(t, []string{"value", "value2"}, apiGatewayProxyRequest.MultiValueQueryStringParameters["key"])
			return c.SendStatus(httpStd.StatusOK)
		})

		request := &httpStd.Request{}
		request.Method = "GET"
		request.URL = &url.URL{
			Path:     "/fiber/path/123",
			RawQuery: "key=value&key=value2&foo=bar",
		}
		request.Header = httpStd.Header{}
		request.Header.Set("accept", "application/json")
		request.Body = httpStd.NoBody

		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, httpStd.StatusOK, resp.StatusCode)
	})

	t.Run("should be convert *http.Request to events.APIGatewayProxyRequest using method POST", func(t *testing.T) {
		request := &httpStd.Request{}
		request.Method = "POST"
		request.URL = &url.URL{}
		request.URL.Path = "/path"
		request.Header = httpStd.Header{}
		request.Header.Set("key", "value")
		request.Body = io.NopCloser(bytes.NewReader([]byte("{\"foo\":\"bar\"}")))

		apiGatewayProxyRequest := http.ConvertHTTPRequestToAPIGatewayProxyRequest(ctx, request)
		assert.NotNil(t, apiGatewayProxyRequest)
		assert.Equal(t, "POST", apiGatewayProxyRequest.HTTPMethod)
		assert.Equal(t, "/path", apiGatewayProxyRequest.Path)
		assert.Equal(t, "{\"foo\":\"bar\"}", apiGatewayProxyRequest.Body)
		assert.Equal(t, 0, len(apiGatewayProxyRequest.MultiValueQueryStringParameters))
		assert.Equal(t, 0, len(apiGatewayProxyRequest.PathParameters))
		assert.Equal(t, 1, len(apiGatewayProxyRequest.MultiValueHeaders))
		assert.Equal(t, []string(nil), apiGatewayProxyRequest.MultiValueHeaders["key"])
	})

}
