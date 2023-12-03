package kendra

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	"golang.org/x/exp/slog"
)

var Client *kendra.Client

const region = "eu-west-1"

// Init the kendra client
func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		slog.Error("Error loading AWS config", "error", err)
		os.Exit(1)
	}
	Client = kendra.NewFromConfig(cfg)
}

func Retrieve(client *kendra.Client, query string) (*kendra.RetrieveOutput, error) {
	// use the retrieve api to query kendra
	// https://docs.aws.amazon.com/sdk-for-go/api/service/kendra/#Kendra.RetrieveDocument
	// https://docs.aws.amazon.com/sdk-for-go/api/service/kendra/#Kendra.Query
	// https://docs.aws.amazon.com/sdk-for-go/api/service/kendra/#Kendra.QueryResultItem
	// set parameters
	index := os.Getenv("KENDRA_ID")

	parameters := &kendra.RetrieveInput{
		IndexId:   &index,
		QueryText: &query,
		AttributeFilter: &types.AttributeFilter{
			AndAllFilters: []types.AttributeFilter{
				{
					EqualsTo: &types.DocumentAttribute{
						Key: aws.String("_language_code"),
						Value: &types.DocumentAttributeValue{
							StringValue: aws.String("en"),
						},
					},
				},
			},
		},
	}
	// do retrieve
	resp, err := client.Retrieve(context.Background(), parameters)
	if err != nil {
		slog.Error("Error retrieving document", "error", err)
		os.Exit(1)
	}

	return resp, nil
}
