package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
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
	weaponReader := addedWeaponReader(stack)
	addedStateMachine(stack, pageFetcher, contentFetcher, scraper)
	weaponTable := addedDynamoDBWeaponTable(stack)
	weaponTable.GrantFullAccess(scraper)
	weaponTable.GrantFullAccess(weaponReader)
	weaponReaderAPI := addedWeaponReaderAPI(stack, weaponReader)
	weaponReaderAPI.Root().AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(weaponReader, nil), nil)
	// weaponReaderAPI.Root().AddMethod()
	return stack
}

func addedScraper(stack awscdk.Stack) awslambdago.GoFunction {
	timeout := 10.0
	return awslambdago.NewGoFunction(stack, jsii.String("ffxiv-scraper"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("ffxiv-scraper/main.go"),
		FunctionName: jsii.String("ffxiv-scraper"),
		Timeout:      awscdk.Duration_Minutes(&timeout),
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
	timeout := 2.0
	return awslambdago.NewGoFunction(stack, jsii.String("ffxiv-page-fetcher"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("ffxiv-page-fetcher/main.go"),
		FunctionName: jsii.String("ffxiv-page-fetcher"),
		Timeout:      awscdk.Duration_Minutes(&timeout),
	})
}

func addedWeaponReader(stack awscdk.Stack) awslambdago.GoFunction {
	timeout := 2.0
	return awslambdago.NewGoFunction(stack, jsii.String("ffxiv-weapon-reader"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Entry:        jsii.String("ffxiv-weapon-reader/main.go"),
		FunctionName: jsii.String("ffxiv-weapon-reader"),
		Timeout:      awscdk.Duration_Minutes(&timeout),
	})
}

func newBoolPointer(v bool) *bool {
	return &v
}

func addedWeaponReaderAPI(stack awscdk.Stack, lambda awslambda.IFunction) awsapigateway.RestApi {
	return awsapigateway.NewLambdaRestApi(stack, jsii.String("ffxiv-weapon-reader-api"), &awsapigateway.LambdaRestApiProps{
		CloudWatchRole: newBoolPointer(false),
		// Deploy:                      new(bool),
		// DeployOptions:               &awsapigateway.StageOptions{},
		// DisableExecuteApiEndpoint:   new(bool),
		// DomainName:                  &awsapigateway.DomainNameOptions{},
		// EndpointExportName:          new(string),
		// EndpointTypes:               &[]awsapigateway.EndpointType{},
		// FailOnWarnings:              new(bool),
		// Parameters:                  &map[string]*string{},
		// Policy:                      nil,
		// RestApiName:                 new(string),
		// RetainDeployments:           new(bool),
		// DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{},
		// DefaultIntegration:          nil,
		// DefaultMethodOptions:        &awsapigateway.MethodOptions{},
		// ApiKeySourceType:            "",
		// BinaryMediaTypes:            &[]*string{},
		// CloneFrom:                   nil,
		// Description:                 new(string),
		// EndpointConfiguration:       &awsapigateway.EndpointConfiguration{},
		// MinimumCompressionSize:      new(float64),
		Handler: lambda,
		// Options:                     &awsapigateway.RestApiProps{},
		// Proxy:                       new(bool),
	})
}

func addedDynamoDBWeaponTable(stack awscdk.Stack) awsdynamodb.Table {
	return awsdynamodb.NewTable(stack, jsii.String("ffxiv-weapon-table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("name"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		// SortKey: &awsdynamodb.Attribute{
		// 	Name: nil,
		// 	Type: "",
		// },
		BillingMode:                "",
		ContributorInsightsEnabled: nil,
		Encryption:                 "",
		EncryptionKey:              nil,
		PointInTimeRecovery:        nil,
		ReadCapacity:               nil,
		RemovalPolicy:              "",
		ReplicationRegions:         nil,
		ReplicationTimeout:         nil,
		ServerSideEncryption:       nil,
		Stream:                     "",
		TimeToLiveAttribute:        nil,
		WriteCapacity:              nil,
		KinesisStream:              nil,
		TableName:                  jsii.String("FFXIVWeapon"),
	})
}
func addedStateMachine(stack awscdk.Stack, pageFetcher awslambdago.GoFunction, contentFetcher awslambdago.GoFunction, scraper awslambdago.GoFunction) {
	loop := awsstepfunctions.NewMap(stack, jsii.String("loop"), &awsstepfunctions.MapProps{
		Comment:        nil,
		InputPath:      nil,
		ItemsPath:      jsii.String("$.Payload.pages"),
		MaxConcurrency: nil,
		OutputPath:     nil,
		Parameters:     nil,
		ResultPath:     nil,
		ResultSelector: nil,
	})
	cf := awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("content-fetcher"), &awsstepfunctionstasks.LambdaInvokeProps{
		LambdaFunction: contentFetcher,
	})
	sc := awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("scraper"), &awsstepfunctionstasks.LambdaInvokeProps{
		LambdaFunction: scraper,
	})
	loop.Iterator(cf.Next(sc))
	awsstepfunctions.NewStateMachine(stack, jsii.String("ffxiv-content-flow"), &awsstepfunctions.StateMachineProps{
		Definition: awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("page-fetcher"), &awsstepfunctionstasks.LambdaInvokeProps{
			LambdaFunction: pageFetcher,
		}).Next(loop),
		// .Next(awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("content-fetcher"), &awsstepfunctionstasks.LambdaInvokeProps{
		// 	LambdaFunction: contentFetcher,
		// })).Next(awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("scraper"), &awsstepfunctionstasks.LambdaInvokeProps{
		// 	LambdaFunction: scraper,
		// })),
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
