package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-playground/validator/v10"
	"github.com/knoebber/comptche-shop/lambda/newsletter"
	"github.com/knoebber/comptche-shop/lambda/util"
)

// HandleRequest adds a new newsletter subscription dynamodb.
// Email must not already exist.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		alreadyExists *dynamodb.ConditionalCheckFailedException
		body          struct {
			Email    string `json:"email"`
			TestMode bool   `json:"testMode"`
		}
	)
	validate := validator.New()

	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	if err = json.Unmarshal([]byte(request.Body), &body); err != nil {
		response.StatusCode = 400
		return
	}
	if body.Email == "" {
		err = errors.New("email is required")
		response.StatusCode = 400
		return
	}
	log.Printf("subscribing %q to newsletter", body.Email)

	if err = validate.Var(body.Email, "email"); err != nil {
		response.StatusCode = 400
		return
	}

	svc := dynamodb.New(session.Must(session.NewSession()))
	if err = newsletter.Subscribe(svc, body.Email, body.TestMode); err == nil {
		response.StatusCode = 200
		return
	}

	log.Print(err)

	if errors.As(err, &alreadyExists) {
		util.SetResponseMessage(&response, "Already subscribed to the newsletter.")
		err = nil
		return
	}

	response.StatusCode = 500
	return
}

func main() {
	lambda.Start(HandleRequest)
}
