package terraformtoresources

import (
	"testing"

	"github.com/diagram-code-generator/resources/pkg/resources"
	hcl "github.com/joselitofilho/hcl-parser-go/pkg/parser/hcl"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"

	"github.com/stretchr/testify/require"
)

func TestTransformer_Transform(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *hcl.Config
	}

	lambdaResource := resources.NewGenericResource("1", "exampleReceiver", awsresources.LambdaType.String())

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "empty",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &hcl.Config{},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "API Gateway route",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_apigatewayv2_route",
							Name:   "apigw_route_example_api_receiver",
							Labels: []string{"aws_apigatewayv2_route", "apigw_route_example_api_receiver"},
							Attributes: map[string]any{
								"api_id":    "aws_apigatewayv2_api.mystack_api.id",
								"route_key": "POST /v1/examples",
								"target":    "integrations/${aws_apigatewayv2_integration.example_api_receiver.id}",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "POST /v1/examples", awsresources.APIGatewayType.String())},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "API Gateway integration",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_apigatewayv2_integration",
							Name:   "example_api_receiver",
							Labels: []string{"aws_apigatewayv2_integration", "example_api_receiver"},
							Attributes: map[string]any{
								"api_id":          "aws_apigatewayv2_api.mystack_api.id",
								"integration_uri": "aws_lambda_function.example_api_receiver_lambda.invoke_arn",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "cloudwatch event rune",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_cloudwatch_event_rule",
							Name:   "example_receiver_cron_rule",
							Labels: []string{"aws_cloudwatch_event_rule", "example_receiver_cron_rule"},
							Attributes: map[string]any{
								"schedule_expression": "cron()",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{resources.NewGenericResource("1", "cron()",
					awsresources.CronType.String())},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "cron",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_cloudwatch_event_target",
							Name:   "example_receiver_cron_target",
							Labels: []string{"aws_cloudwatch_event_target", "example_receiver_cron_target"},
							Attributes: map[string]any{
								"rule": "aws_cloudwatch_event_rule.example_receiver_cron.name",
								"arn":  "aws_lambda_function.example_receiver_lambda.arn",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "endpoint",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_apigatewayv2_domain_name",
							Name:   "my_restful_api",
							Labels: []string{"aws_apigatewayv2_domain_name", "my_restful_api"},
							Attributes: map[string]any{
								"domain_name": "local.api_domain",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "local.api_domain", awsresources.EndpointType.String())},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_kinesis_stream",
							Name:   "my_stream",
							Labels: []string{"aws_kinesis_stream", "my_stream"},
							Attributes: map[string]any{
								"name": "MyStream",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "MyStream", awsresources.KinesisType.String())},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "lambda event source mapping with kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_event_source_mapping",
							Name:   "example_checker_lambda_sqs_trigger",
							Labels: []string{"aws_lambda_event_source_mapping", "example_checker_lambda_sqs_trigger"},
							Attributes: map[string]any{
								"event_source_arn": "aws_kinesis_stream.my_stream_kinesis.arn",
								"function_name":    "aws_lambda_function.example_checker_lambda.arn",
							},
						},
						{
							Type:   "aws_kinesis_stream",
							Name:   "my_stream",
							Labels: []string{"aws_kinesis_stream", "my_stream"},
							Attributes: map[string]any{
								"name": "MyStream",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "MyStream", awsresources.KinesisType.String())},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "example_receiver_lambda",
							Labels: []string{"aws_lambda_function", "example_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "exampleReceiver",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "lambda as module",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{{
						Labels: []string{"example_receiver_lambda"},
						Attributes: map[string]any{
							"lambda_function_name": "exampleReceiver",
						},
					}},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "s3 bucket",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_s3_bucket",
							Name:   "my_bucket",
							Labels: []string{"aws_s3_bucket", "my_bucket"},
							Attributes: map[string]any{
								"bucket": "var.client-var.environment-my-bucket",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "var.client-var.environment-my-bucket",
						awsresources.S3Type.String())},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "sqs",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_sqs_queue",
							Name:   "my_sqs",
							Labels: []string{"aws_sqs_queue", "my_sqs"},
							Attributes: map[string]any{
								"name": "my-queue",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "my-queue", awsresources.SQSType.String())},
				Relationships: []resources.Relationship{},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := NewTransformer(tc.fields.yamlConfig, tc.fields.tfConfig).Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromLambdaToResourceFromEnvar(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *hcl.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", awsresources.LambdaType.String())
	bqResource := resources.NewGenericResource("2", "google", awsresources.GoogleBQType.String())
	dbResource := resources.NewGenericResource("2", "var.doc_db_host", awsresources.DatabaseType.String())
	kinesisResource := resources.NewGenericResource("2", "MyStream", awsresources.KinesisType.String())
	restfulAPIResource := resources.NewGenericResource("2", "MyRestful", awsresources.RestfulAPIType.String())
	s3BucketResource := resources.NewGenericResource("2", "my-bucket", awsresources.S3Type.String())
	sqsResource := resources.NewGenericResource("2", "var.variable1-my-queue", awsresources.SQSType.String())

	kinesisStreamTerraform := &hcl.Resource{
		Type:   "aws_kinesis_stream",
		Name:   "my_stream_kinesis",
		Labels: []string{"aws_kinesis_stream", "my_stream_kinesis"},
		Attributes: map[string]any{
			"name": "MyStream",
		},
	}

	s3BucketTerraform := &hcl.Resource{
		Type:   "aws_s3_bucket",
		Name:   "my_bucket",
		Labels: []string{"aws_s3_bucket", "my_bucket"},
		Attributes: map[string]any{
			"bucket": "my-bucket",
		},
	}

	sqsTerraform := &hcl.Resource{
		Type:   "aws_sqs_queue",
		Name:   "my_queue",
		Labels: []string{"aws_sqs_queue", "my_queue"},
		Attributes: map[string]any{
			"name":                       "var.variable1-my-queue",
			"visibility_timeout_seconds": "720",
		},
	}

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "lambda as resource with google BQ",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"environment": map[string]map[string]any{
									"variables": {
										"GOOGLE_BQ_PROJECT_ID": "google",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, bqResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: bqResource}},
			},
		},
		{
			name: "lambda as module with google BQ",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"lambda_function_env_vars": map[string]any{
									"BQ_PROJECT_ID": "google",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, bqResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: bqResource}},
			},
		},
		{
			name: "lambda as resource with database",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"environment": map[string]map[string]any{
									"variables": {
										"DOCDB_HOST": "var.doc_db_host",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, dbResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: dbResource}},
			},
		},
		{
			name: "lambda as module with database",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"lambda_function_env_vars": map[string]any{
									"DOCDB_HOST": "var.doc_db_host",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, dbResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: dbResource}},
			},
		},
		{
			name: "lambda as resource with kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_STREAM_KINESIS_STREAM_URL": "aws_kinesis_stream.my_stream_kinesis.name",
									},
								},
							},
						},
						kinesisStreamTerraform,
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, kinesisResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: kinesisResource}},
			},
		},
		{
			name: "lambda as module with kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"lambda_function_env_vars": map[string]any{
									"KINESIS_STREAM_URL": "aws_kinesis_stream.my_stream_kinesis.name",
								},
							},
						},
					},
					Resources: []*hcl.Resource{kinesisStreamTerraform},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, kinesisResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: kinesisResource}},
			},
		},
		{
			name: "lambda as resource with restful API",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_RESTFUL_API_BASE_URL": "MyRestful",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, restfulAPIResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: restfulAPIResource}},
			},
		},
		{
			name: "lambda as module with restful API",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"lambda_function_env_vars": map[string]any{
									"API_BASE_URL": "MyRestful",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, restfulAPIResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: restfulAPIResource}},
			},
		},
		{
			name: "lambda as resource with S3 Bucket",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_BUCKET_S3_BUCKET": "my-bucket",
									},
								},
							},
						},
						s3BucketTerraform,
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, s3BucketResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: s3BucketResource}},
			},
		},
		{
			name: "lambda as module with S3 Bucket",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"lambda_function_env_vars": map[string]any{
									"BUCKET_NAME": "my-bucket",
								},
							},
						},
					},
					Resources: []*hcl.Resource{s3BucketTerraform},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, s3BucketResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: s3BucketResource}},
			},
		},
		{
			name: "lambda as resource with SQS",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_QUEUE_SQS_QUEUE_URL": "var.variable1-my-queue",
									},
								},
							},
						},
						sqsTerraform,
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, sqsResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: sqsResource}},
			},
		},
		{
			name: "lambda as module with SQS",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
								"lambda_function_env_vars": map[string]any{
									"SQS_QUEUE_URL": "var.variable1-my-queue",
								},
							},
						},
					},
					Resources: []*hcl.Resource{sqsTerraform},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, sqsResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: sqsResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromCronToResource(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *hcl.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", awsresources.LambdaType.String())
	cronResource := resources.NewGenericResource("2", "cron(0 3 * * ? *)", awsresources.CronType.String())

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "from cron to lambda",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Modules: []*hcl.Module{
						{
							Labels: []string{"my_receiving_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
							},
						},
					},
					Resources: []*hcl.Resource{
						{
							Type:   "aws_cloudwatch_event_rule",
							Name:   "cron",
							Labels: []string{"aws_cloudwatch_event_rule", "cron"},
							Attributes: map[string]any{
								"schedule_expression": "cron(0 3 * * ? *)",
							},
						},
						{
							Type:   "aws_cloudwatch_event_target",
							Name:   "cron_event",
							Labels: []string{"aws_cloudwatch_event_target", "cron_event"},
							Attributes: map[string]any{
								"rule": "aws_cloudwatch_event_rule.cron.name",
								"arn":  "module.my_receiving_lambda.function_arn",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, cronResource},
				Relationships: []resources.Relationship{{Source: cronResource, Target: lambdaResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformEndpointAPIGatewayLambda(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *hcl.Config
	}

	endpointResource := resources.NewGenericResource("1", "local.api_domain", awsresources.EndpointType.String())
	apigResource := resources.NewGenericResource("2", "POST /v1/examples", awsresources.APIGatewayType.String())
	lambdaResource := resources.NewGenericResource("3", "myReceiver", awsresources.LambdaType.String())

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "happy path",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &hcl.Config{
					Resources: []*hcl.Resource{
						{
							Type:   "aws_apigatewayv2_domain_name",
							Name:   "my_restful_api",
							Labels: []string{"aws_apigatewayv2_domain_name", "my_stack_api"},
							Attributes: map[string]any{
								"domain_name": "local.api_domain",
							},
						},
						{
							Type:   "aws_apigatewayv2_route",
							Name:   "apigw_route_example_api_receiver",
							Labels: []string{"aws_apigatewayv2_route", "apigw_route_example_api_receiver"},
							Attributes: map[string]any{
								"api_id":    "aws_apigatewayv2_api.my_stack_api.id",
								"route_key": "POST /v1/examples",
								"target":    "integrations/${aws_apigatewayv2_integration.my_receiver.id}",
							},
						},
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "myReceiver",
							},
						},
						{
							Type:   "aws_apigatewayv2_integration",
							Name:   "my_receiver",
							Labels: []string{"aws_apigatewayv2_integration", "my_receiver"},
							Attributes: map[string]any{
								"api_id":          "aws_apigatewayv2_api.my_stack_api.id",
								"integration_uri": "aws_lambda_function.my_receiver_lambda.invoke_arn",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{endpointResource, apigResource, lambdaResource},
				Relationships: []resources.Relationship{
					{Source: endpointResource, Target: apigResource},
					{Source: apigResource, Target: lambdaResource},
				},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_hasResourceMatched(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *hcl.Config
	}

	type args struct {
		res     resources.Resource
		filters config.Filters
	}

	lambdaResource := resources.NewGenericResource("id", "myLambda", awsresources.LambdaType.String())

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "match",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &hcl.Config{},
			},
			args: args{
				res: lambdaResource,
				filters: config.Filters{
					"lambda": config.Filter{Match: []string{"^my"}},
				},
			},
			want: true,
		},
		{
			name: "not match",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &hcl.Config{},
			},
			args: args{
				res: lambdaResource,
				filters: config.Filters{
					"lambda": config.Filter{NotMatch: []string{"^my"}},
				},
			},
			want: false,
		},
		{
			name: "nil resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &hcl.Config{},
			},
			args: args{
				res: nil,
				filters: config.Filters{
					"lambda": config.Filter{NotMatch: []string{"^my"}},
				},
			},
			want: false,
		},
		{
			name: "no filter",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &hcl.Config{},
			},
			args: args{
				res:     lambdaResource,
				filters: nil,
			},
			want: true,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(tc.fields.yamlConfig, tc.fields.tfConfig)

			got := tr.hasResourceMatched(tc.args.res, tc.args.filters)

			require.Equal(t, tc.want, got)
		})
	}
}
