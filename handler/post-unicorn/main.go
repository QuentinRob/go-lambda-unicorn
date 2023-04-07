package main

import (
	"connector"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"models"
)

type PostUnicornDto struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mysqlConnector := connector.MysqlConnector{}
	db, err := mysqlConnector.GetDB()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	var postUnicornDto PostUnicornDto

	err = json.Unmarshal([]byte(event.Body), &postUnicornDto)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	unicorn := &models.Unicorn{
		Name:  postUnicornDto.Name,
		Color: postUnicornDto.Color,
	}

	db.Create(unicorn)

	body, err := json.Marshal(*unicorn)

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
