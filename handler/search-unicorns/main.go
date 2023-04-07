package main

import (
	"connector"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm"
	"models"
	"strconv"
)

type Paginator struct {
	Page  int
	Limit int
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

	page, err := strconv.Atoi(event.QueryStringParameters["page"])

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	limit, err := strconv.Atoi(event.QueryStringParameters["limit"])

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	var unicorns []models.Unicorn
	paginator := Paginator{
		Page:  page,
		Limit: limit,
	}

	db.Scopes(Paginate(&paginator)).Find(&unicorns)

	body, err := json.Marshal(unicorns)

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

func Paginate(paginator *Paginator) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := paginator.Page
		if page <= 0 {
			page = 1
		}

		pageSize := paginator.Limit
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func main() {
	lambda.Start(handler)
}
