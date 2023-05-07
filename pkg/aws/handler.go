package aws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fi3te/monthly-planning-backend/pkg/auth"
	"github.com/fi3te/monthly-planning-backend/pkg/config"
)

const primarySlot = "A"
const secondarySlot = "B"
const currentDataSlot = "Z"

var cfg *config.Config
var db dynamodb.Client

func init() {
	var err error
	cfg, err = config.ReadConfig()
	if err != nil {
		panic(err)
	}
	db = createDynamoDbClient()
}

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	authorized := auth.IsAuthorized(cfg, request.Headers)
	if !authorized {
		return responseWithStatus(http.StatusUnauthorized)
	}

	slot := request.QueryStringParameters["slot"]
	if !contains([]string{"", primarySlot, secondarySlot}, slot) {
		return responseWithStatus(http.StatusBadRequest)
	}
	if slot == "" {
		slot = currentDataSlot
	}

	method := request.RequestContext.HTTP.Method
	if method == http.MethodGet {
		return handleGetRequest(slot)
	}
	if method == http.MethodPut {
		return handlePutRequest(slot, request.Body)
	}
	return responseWithStatus(http.StatusBadRequest)
}

func handleGetRequest(slot string) (events.LambdaFunctionURLResponse, error) {
	value, err := getItem(context.TODO(), cfg, &db, slot)
	if err != nil {
		log.Fatal(err)
		return responseWithStatus(http.StatusInternalServerError)
	}
	if value != nil {
		return events.LambdaFunctionURLResponse{Body: value.Data, StatusCode: 200}, nil
	} else {
		return responseWithStatus(http.StatusNotFound)
	}
}

func handlePutRequest(slot string, requestBody string) (events.LambdaFunctionURLResponse, error) {
	if slot == currentDataSlot {
		return responseWithStatus(http.StatusBadGateway)
	}
	dataSlice := []MonthlyPlanningData{
		{Slot: slot, Data: requestBody},
		{Slot: currentDataSlot, Data: requestBody},
	}
	err := putItems(context.TODO(), cfg, &db, dataSlice)
	if err != nil {
		log.Fatal(err)
		return responseWithStatus(http.StatusInternalServerError)
	}
	return events.LambdaFunctionURLResponse{StatusCode: 204}, nil
}

func responseWithStatus(code int) (events.LambdaFunctionURLResponse, error) {
	body, err := toJsonString(http.StatusText(code))
	return events.LambdaFunctionURLResponse{
		StatusCode: code,
		Body:       body,
	}, err
}

func toJsonString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

func contains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}
