package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"crud/model"
	httpmodel "crud/model/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const PER_PAGE = 3

type ListResult struct {
	Count    int64
	Students []*model.Student
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	pageParam, ok := request.QueryStringParameters["page"]
	if !ok {
		pageParam = "0"
	}

	page, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err parsing query parameter page: %s\n", pageParam),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}

	query, err := db.Prepare("select id, first_name, last_name, email from students limit ?, ?")
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err opening prepare statement: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer query.Close()

	rows, err := query.Query(page*PER_PAGE, PER_PAGE)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Query failed, err: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	students := make([]*model.Student, 0)
	for {
		hasNext := rows.Next()
		if !hasNext {
			break
		}
		var id int64
		var firstName string
		var lastName string
		var email string
		err := rows.Scan(&id, &firstName, &lastName, &email)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Scan failed, err: %s\n", err),
				Headers:    httpmodel.Cors,
				StatusCode: 500,
			}, nil
		}
		students = append(students, &model.Student{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		})
	}

	queryCount, err := db.Prepare("select count(1) from students")
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err opening prepare statement for count: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer queryCount.Close()

	var count int64
	err = queryCount.QueryRow().Scan(&count)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err getting count: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	buf := bytes.NewBufferString("")
	encode := json.NewEncoder(buf)
	encode.SetIndent("", "  ")
	err = encode.Encode(ListResult{
		Count:    count,
		Students: students,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("JSON Encode failed, err: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       buf.String(),
		Headers:    httpmodel.Cors,
		StatusCode: 200,
	}, nil
}

func main() {
	conn, err := sql.Open("mysql", os.Getenv("ConnectionString"))
	if err != nil {
		log.Printf("err: %s\n", err)
		panic(err)
	}
	db = conn
	defer db.Close()
	lambda.Start(handler)
}
