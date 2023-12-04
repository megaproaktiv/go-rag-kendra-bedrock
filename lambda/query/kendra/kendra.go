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

type Document struct {
	Excerpt *string `json:"excerpt"`
	Title   *string `json:"title"`
	Page    *int    `json:"page"`
	Link    *string `json:"link"`
}

type Query struct {
	Question string  `json:"question"`
	Category *string `json:"category"`
	Version  *string `json:"version"`
}

var Client *kendra.Client

var region string
var languageCode string

// Init the kendra client
func init() {
	languageCode = os.Getenv("KENDRA_LANGUAGE_CODE")
	region = os.Getenv("KENDRA_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		slog.Error("Error loading AWS config", "error", err)
		os.Exit(1)
	}
	Client = kendra.NewFromConfig(cfg)
}

func Retrieve(client *kendra.Client, query Query) (*kendra.RetrieveOutput, error) {
	// use the retrieve api to query kendra
	// https://docs.aws.amazon.com/sdk-for-go/api/service/kendra/#Kendra.RetrieveDocument
	// https://docs.aws.amazon.com/sdk-for-go/api/service/kendra/#Kendra.Query
	// https://docs.aws.amazon.com/sdk-for-go/api/service/kendra/#Kendra.QueryResultItem
	// set parameters
	index := os.Getenv("KENDRA_ID")

	// Set filter if necessary
	filter := &types.AttributeFilter{
		AndAllFilters: []types.AttributeFilter{
			{
				EqualsTo: &types.DocumentAttribute{
					Key: aws.String("_language_code"),
					Value: &types.DocumentAttributeValue{
						StringValue: aws.String(languageCode),
					},
				},
			},
		},
	}

	if query.Category != nil {
		filter.AndAllFilters = append(filter.AndAllFilters, types.AttributeFilter{
			EqualsTo: &types.DocumentAttribute{
				Key: aws.String("_category"),
				Value: &types.DocumentAttributeValue{
					StringValue: query.Category,
				},
			},
		})
	}

	if query.Version != nil {
		filter.AndAllFilters = append(filter.AndAllFilters, types.AttributeFilter{
			EqualsTo: &types.DocumentAttribute{
				Key: aws.String("_version"),
				Value: &types.DocumentAttributeValue{
					StringValue: query.Version,
				},
			},
		})
	}

	slog.Info("Filter", "filter", filter)
	parameters := &kendra.RetrieveInput{
		IndexId:         &index,
		QueryText:       &query.Question,
		PageSize:        aws.Int32(20),
		AttributeFilter: filter,
	}
	// do retrieve
	resp, err := client.Retrieve(context.Background(), parameters)
	if err != nil {
		slog.Error("Error retrieving document", "error", err)
		os.Exit(1)
	}

	return resp, nil
}
