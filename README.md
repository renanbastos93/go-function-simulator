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
	"github.com/gofiber/fiber/v2"
	"github.com/renanbastos93/go-function-simulator/pkg/http"
)

func main() {
	app := fiber.New()

	app.Get("/path/:id", func(c *fiber.Ctx) error {
		apiGatewayProxyRequest := http.ConvertHTTPRequestToAPIGatewayProxyRequest(c.Context(), c)
        fmt.Println("API Gateway Proxy Request:", apiGatewayProxyRequest)
		// call your lambda function
		// like this: usecase.Lambda(ctx, apiGatewayProxyRequest)
		return c.SendStatus(fiber.StatusOK)
	})

	_ = app.Listen(":3000")
}

```

## Usage Gorilla Mux
```go
package main

import (
	"context"
	"fmt"

	stdHttp "net/http"

	"github.com/gorilla/mux"
	"github.com/renanbastos93/go-function-simulator/pkg/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/path/{id}", func(w stdHttp.ResponseWriter, r *stdHttp.Request) {
		apiGatewayProxyRequest := http.ConvertHTTPRequestToAPIGatewayProxyRequest(context.Background(), r)
		fmt.Println("API Gateway Proxy Request:", apiGatewayProxyRequest)
		/// call your lambda function
		// like this: usecase.Lambda(ctx, apiGatewayProxyRequest)
		w.WriteHeader(stdHttp.StatusOK)
	})

	_ = stdHttp.ListenAndServe(":3000", r)
}
```