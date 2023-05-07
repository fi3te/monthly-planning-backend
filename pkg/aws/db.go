package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fi3te/monthly-planning-backend/pkg/config"
)

type MonthlyPlanningData struct {
	Slot string `dynamodbav:"slot"`
	Data string `dynamodbav:"data"`
}

func createDynamoDbClient() dynamodb.Client {
	sdkConfig, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return *dynamodb.NewFromConfig(sdkConfig)
}

func putItems(ctx context.Context, cfg *config.Config, db *dynamodb.Client, dataSlice []MonthlyPlanningData) error {
	putRequests := make([]types.WriteRequest, 0)
	for _, data := range dataSlice {
		item, err := attributevalue.MarshalMap(data)
		if err != nil {
			return err
		}
		putRequests = append(putRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			cfg.TableName: putRequests,
		},
	}

	_, err := db.BatchWriteItem(ctx, input)
	return err
}

func getItem(ctx context.Context, cfg *config.Config, db *dynamodb.Client, slot string) (*MonthlyPlanningData, error) {
	slotAttributeValue, err := attributevalue.Marshal(slot)
	if err != nil {
		return nil, err
	}
	key := map[string]types.AttributeValue{"slot": slotAttributeValue}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(cfg.TableName),
		Key:       key,
	}

	res, err := db.GetItem(ctx, input)
	if err != nil || res.Item == nil {
		return nil, err
	}

	data := new(MonthlyPlanningData)
	err = attributevalue.UnmarshalMap(res.Item, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
