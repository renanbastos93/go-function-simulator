package awslambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type Lambda interface {
	Handler(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error)
}
