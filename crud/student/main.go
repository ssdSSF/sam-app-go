package main

import (
	"crud/model"
	httpmodel "crud/model/http"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch method := request.HTTPMethod; method {
	case "OPTIONS":
		return events.APIGatewayProxyResponse{
			Headers:    httpmodel.Cors,
			StatusCode: 200,
		}, nil
	case "POST":
		return handlePOST(request)
	case "GET":
		return handleGET(request)
	case "PUT":
		return handlePUT(request)
	case "DELETE":
		return handleDELETE(request)
	default:
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("UNKNOW HTTP METHOD %s", method),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}
}

func handleDELETE(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idParam, ok := request.QueryStringParameters["id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			Body:       "query param id is not presented:",
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err parsing query parameter id: %s\n", idParam),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}

	delete, err := db.Prepare("delete from students where id = ?")
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err prepare query student id: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer delete.Close()

	_, err = delete.Exec(id)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err query and scan student: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    httpmodel.Cors,
		StatusCode: 204,
	}, nil
}

func handlePUT(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idParam, ok := request.QueryStringParameters["id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			Body:       "query param id is not presented:",
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err parsing query parameter id: %s\n", idParam),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}

	queryExist, err := db.Prepare("select id from students where id = ?")
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err prepare query student id exists: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer queryExist.Close()

	err = queryExist.QueryRow(id).Scan(&id)

	if err == sql.ErrNoRows {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("student id '%d' does not exists", id),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err query student id exists: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	var student model.Student

	valid, response := validateJSON(request, &student)
	if !valid {
		return response, nil
	}

	update, err := db.Prepare("update students set first_name = ?, last_name = ?, email = ? where id = ?")
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err prepare query student id: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer update.Close()

	_, err = update.Exec(student.FirstName, student.LastName, student.Email, id)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err update student: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	student.Id = id
	body, err := json.Marshal(student)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("unable to marshal student: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    httpmodel.Cors,
		StatusCode: 200,
	}, nil
}

func handleGET(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idParam, ok := request.QueryStringParameters["id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			Body:       "query param id is not presented:",
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err parsing query parameter id: %s\n", idParam),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}, nil
	}

	query, err := db.Prepare("select first_name, last_name, email from students where id = ?")
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err prepare query student id: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer query.Close()

	var student model.Student
	err = query.QueryRow(id).Scan(&student.FirstName, &student.LastName, &student.Email)
	if err == sql.ErrNoRows {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("no students found: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 404,
		}, nil
	}
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("err query and scan student: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	body, err := json.Marshal(student)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("unable to marshal student: %s\n", err),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    httpmodel.Cors,
		StatusCode: 200,
	}, nil
}

func handlePOST(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var student model.Student

	valid, response := validateJSON(request, &student)
	if !valid {
		return response, nil
	}

	insert, err := db.Prepare("INSERT INTO students (first_name, last_name, email) VALUES( ?, ?, ? )")
	if err != nil {
		log.Printf("failed to prepare statement , err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "failed to prepare statement",
		})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer insert.Close()

	result, err := insert.Exec(student.FirstName, student.LastName, student.Email)
	if err != nil {
		log.Printf("insert into students failed, err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "insert into students failed",
		})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	studentId, err := result.LastInsertId()
	if err != nil {
		log.Printf("failed to get LastInsertId, err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "failed to get LastInsertId",
		})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	query, err := db.Prepare("select first_name, last_name, email from students where id = ?")
	if err != nil {
		log.Printf("failed to prepare query statement , err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "failed to prepare query statement",
		})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}
	defer query.Close()

	var firstName string
	var lastName string
	var email string
	err = query.QueryRow(studentId).Scan(&firstName, &lastName, &email)
	if err != nil {
		log.Printf("failed to query student , err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "failed to prepare query student",
		})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	body, err := json.Marshal(model.Student{
		Id:        studentId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	})

	if err != nil {
		log.Printf("failed to marshal student json, err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: fmt.Sprintf("failed to marshal student json, err: %s", err),
		})
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}, nil
	}

	httpmodel.Cors["content-type"] = "application-json"
	defer delete(httpmodel.Cors, "content-type")
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    httpmodel.Cors,
		StatusCode: 200,
	}, nil
}

func validateJSON(request events.APIGatewayProxyRequest, student *model.Student) (bool, events.APIGatewayProxyResponse) {

	err := json.NewDecoder(strings.NewReader(request.Body)).Decode(student)
	if err != nil {
		log.Printf("failed to decode json, err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "cannot decode json",
		})
		return false, events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}
	}

	if student.FirstName == "" {
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "FirstName cannot be empty",
		})
		return false, events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}
	}

	if student.LastName == "" {
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "LastName cannot be empty",
		})
		return false, events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}
	}

	if student.Email == "" {
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "Email cannot be empty",
		})
		return false, events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}
	}

	var id int64
	duplicate, err := db.Prepare("select id from students where email = ?")
	if err != nil {
		log.Printf("failed to prepare statement duplicate, err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "failed to prepare statement duplicate",
		})
		return false, events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}
	}
	defer duplicate.Close()

	err = duplicate.QueryRow(student.Email).Scan(&id)
	if err != sql.ErrNoRows && id != student.Id {
		// update: null, id == student.id -> pass
		// update: null, id != student.id -> fail
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: fmt.Sprintf("Email %s has been taken", student.Email),
		})
		return false, events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 400,
		}
	} else if err != sql.ErrNoRows && err != nil {
		// create, err has to be sql.ErrNoRows
		// update, update email, err has to be sql.ErrNoRows
		log.Printf("failed to check duplicate email, err: %s\n", err)
		body, _ := json.Marshal(httpmodel.CreateJSON{
			Success: false,
			Message: "failed to check duplicat email",
		})
		return false, events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    httpmodel.Cors,
			StatusCode: 500,
		}
	}

	return true, events.APIGatewayProxyResponse{}
}

func main() {
	conn, err := sql.Open("mysql", os.Getenv("ConnectionString"))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		panic(err)
	}
	db = conn
	defer db.Close()

	lambda.Start(handler)
}
