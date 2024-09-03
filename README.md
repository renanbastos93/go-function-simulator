# go-function-simulator
Run functions locally with Go.

## Overview
It is a tool designed to help you run and test your functions locally with Go. It supports various HTTP request types and frameworks such as Fiber, Gorilla Mux, and Chi. This allows you to simulate and debug AWS Lambda functions or any HTTP-based services on your local machine.

## Features

- **Support for Multiple Request Types:** Handle and simulate HTTP requests using different frameworks.
- **Flexible Conversion:** Convert HTTP requests to AWS API Gateway proxy requests and vice versa.
- **Easy Integration:** Seamlessly integrate with your existing Go projects.

## Installation

To install `go-function-simulator`, you can use Go modules to add it to your project:

```bash
$ go get github.com/yourusername/go-function-simulator
```

## Usage Fiber
```go
package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go-function-simulator/http"
)

func main() {
	app := fiber.New()

	app.Get("/path/:id", func(c *fiber.Ctx) error {
		apiGatewayProxyRequest := http.ConvertHTTPRequestToAPIGatewayProxyRequest(c.Context(), c)
		fmt.Println("API Gateway Proxy Request:", apiGatewayProxyRequest)
		return c.SendStatus(fiber.StatusOK)
	})

	// Simulate a request
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path:     "/path/123",
			RawQuery: "query=value",
		},
		Header: http.Header{
			"Content-Type": {"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(`{"foo":"bar"}`))),
	}

	resp, err := app.Test(req)
	if err != nil {
		fmt.Println("Error testing request:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.StatusCode)
}
```

## Usage Gorilla Mux
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gorilla/mux"
	"github.com/yourusername/go-function-simulator/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/path/{id}", func(w http.ResponseWriter, r *http.Request) {
		apiGatewayProxyRequest := http.ConvertHTTPRequestToAPIGatewayProxyRequest(context.Background(), r)
		fmt.Println("API Gateway Proxy Request:", apiGatewayProxyRequest)
		w.WriteHeader(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/path/123?query=value", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	fmt.Println("Response Status Code:", rr.Code)
}
```