package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/aws-sdk-go/aws"
)

func TestLambdaResource(t *testing.T) {
	stack := createStack(awscdk.NewApp(nil), "serverless-test-stack", nil)

	template := assertions.Template_FromStack(stack, nil)

	// Assert that the CloudFormation template has a Lambda resource
	template.HasResourceProperties(aws.String("AWS::Lambda::Function"), map[string]interface{}{
		"Handler": "lambda",
		"Runtime": "go1.x",
		"Timeout": 60,
	})
}
