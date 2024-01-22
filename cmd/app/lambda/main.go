package main

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app"
	"github.com/mauriciozanettisalomao/go-transaction-service/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.InitStructureLogConfig()
	ginLambda = ginadapter.New(app.Endpoints())
}

// Handler is the lambda handler unsing the gin adapter as a wrapper
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
