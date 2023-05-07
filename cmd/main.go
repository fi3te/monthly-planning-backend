package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fi3te/monthly-planning-backend/pkg/aws"
)

func main() {
	lambda.Start(aws.HandleRequest)
}
