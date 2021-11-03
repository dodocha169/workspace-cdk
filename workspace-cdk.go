package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/aws-cdk-go/awscdk/awsstepfunctions"
	"github.com/aws/aws-cdk-go/awscdk/awsstepfunctionstasks"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type WorkspaceCdkStackProps struct {
	awscdk.StackProps
}

func NewWorkspaceCdkStack(scope constructs.Construct, id string, props *WorkspaceCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	scraper := addedScraper(stack)
	contentFetcher := addedContentFetcher(stack)
	pageFetcher := addedPageFetcher(stack)
	addedStateMachine(stack, pageFetcher, contentFetcher, scraper)
	return stack
}

func addedScraper(stack awscdk.Stack) awslambdago.GoFunction {
	return awslambdago.NewGoFunction(stack, jsii.String("ffxiv-scraper"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("ffxiv-scraper/main.go"),
		FunctionName: jsii.String("ffxiv-scraper"),
	})
}

func addedContentFetcher(stack awscdk.Stack) awslambdago.GoFunction {
	timeout := 5.0
	return awslambdago.NewGoFunction(stack, jsii.String("ffxiv-content-fetcher"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("ffxiv-content-fetcher/main.go"),
		FunctionName: jsii.String("ffxiv-content-fetcher"),
		Timeout:      awscdk.Duration_Minutes(&timeout),
	})
}

func addedPageFetcher(stack awscdk.Stack) awslambdago.GoFunction {
	return awslambdago.NewGoFunction(stack, jsii.String("ffxiv-page-fetcher"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("ffxiv-page-fetcher/main.go"),
		FunctionName: jsii.String("ffxiv-page-fetcher"),
	})
}

func addedStateMachine(stack awscdk.Stack, pageFetcher awslambdago.GoFunction, contentFetcher awslambdago.GoFunction, scraper awslambdago.GoFunction) {
	awsstepfunctions.NewStateMachine(stack, jsii.String("ffxiv-content-flow"), &awsstepfunctions.StateMachineProps{
		Definition: awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("page-fetcher"), &awsstepfunctionstasks.LambdaInvokeProps{
			LambdaFunction: pageFetcher,
		}).Next(awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("content-fetcher"), &awsstepfunctionstasks.LambdaInvokeProps{
			LambdaFunction: contentFetcher,
		})).Next(awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("scraper"), &awsstepfunctionstasks.LambdaInvokeProps{
			LambdaFunction: scraper,
		})),
		Logs:             nil,
		Role:             nil,
		StateMachineName: jsii.String("ffxiv-content-flow"),
		StateMachineType: awsstepfunctions.StateMachineType_STANDARD,
		Timeout:          nil,
		TracingEnabled:   nil,
	})
}

func main() {
	app := awscdk.NewApp(nil)

	NewWorkspaceCdkStack(app, "WorkspaceCdkStack", &WorkspaceCdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})
	awscdk.Tags_Of(app).Add(jsii.String("system"), jsii.String("ffxiv"), nil)
	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
