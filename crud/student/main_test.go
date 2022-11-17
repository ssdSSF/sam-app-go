package main

import (
	httpmodel "crud/model/http"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	_ "github.com/go-sql-driver/mysql"
)

func TestHandler(t *testing.T) {

	conn, err := sql.Open("mysql", "root:Welcome1@tcp(localhost:3306)/crud")
	if err != nil {
		fmt.Printf("err: %s\n", err)
		panic(err)
	}
	db = conn

	t.Run("Response cannot be null", func(t *testing.T) {
		_, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
	})

	t.Run("json typo", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			Body: `{"first_name": "lucas", "last_name": "shi", "email": "lucas@srvusd.net"}`,
		})
		if err != nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
		if response.StatusCode != 400 {
			t.Fatalf("response is not 400")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			Body:       `{"FirstName": "lucas", "LastName": "shi", "email": "lucas@srvusd.net"}`,
			HTTPMethod: "POST",
		})
		if err != nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
		if response.StatusCode != 400 {
			t.Fatalf("response is not 400")
		}
		var createJSON httpmodel.CreateJSON
		json.NewDecoder(strings.NewReader(response.Body)).Decode(&createJSON)
		if createJSON.Success == true {
			t.Fatalf("success should not be true")
		}
		if createJSON.Message == "" {
			t.Fatalf("Message is empty")
		}
	})
}
