package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
)

type StackProps struct {
	awscdk.StackProps
}

func main() {
	app := awscdk.NewApp(nil)

	createStack(app, "serverless-api-stack", &StackProps{
		awscdk.StackProps{
			Env: nil,
		},
	})

	app.Synth(nil)
}

func createStack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	// Define the Go Lambda function via the archive
	awslambda.NewFunction(stack, aws.String("Lambda"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(aws.String("function.zip"), nil),
		Handler: aws.String("lambda"),
		Runtime: awslambda.Runtime_GO_1_X(),
		Timeout: awscdk.Duration_Seconds(aws.Float64(60)),
	})

	return stack
}
