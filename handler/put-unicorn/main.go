package main

import (
	"connector"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"models"
	"strconv"
)

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mysqlConnector := connector.MysqlConnector{}
	db, err := mysqlConnector.GetDB()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	id, err := strconv.Atoi(event.PathParameters["id"])

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	var putUnicornDto map[string]interface{}

	err = json.Unmarshal([]byte(event.Body), &putUnicornDto)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	var unicorn models.Unicorn

	db.First(&unicorn, id)
	db.Model(&unicorn).Updates(putUnicornDto)

	body, err := json.Marshal(unicorn)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
