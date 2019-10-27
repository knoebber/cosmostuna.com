package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/product"
)

// TODO make route parameter.
const (
	testMode        = true
	testModeKeyName = "/test-secret-stripe-api-key"
	bucketName      = "cosmostuna-backend"
)

// TODO golinting.
// ProductResponse is the json response that HandleRequest responds with.
type ProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		bytes []byte
	)

	// TODO only cosmostuna.com
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	if err = initStripe(); err != nil {
		response.StatusCode = 500
		return
	}

	responseBody := []ProductResponse{}

	i := product.List(new(stripe.ProductListParams))
	for i.Next() {
		p := i.Product()
		responseBody = append(responseBody, ProductResponse{p.ID, p.Name})
	}

	fmt.Printf("products: %+v\n", responseBody)

	bytes, err = json.Marshal(&responseBody)
	if err != nil {
		return
	}

	response.Body = string(bytes)
	response.StatusCode = 200
	return
}

func initStripe() error {
	// Create a S3 client
	session := session.Must(session.NewSession())
	svc := s3.New(session)

	getInput := s3.GetObjectInput{
		Bucket: aws.String(bucketName),
	}

	// Get the Stripe secret api key
	if testMode {
		getInput.Key = aws.String(testModeKeyName)
	}

	output, err := svc.GetObject(&getInput)
	if err != nil {
		return errors.Wrap(err, "failed to get stripe api secret S3 object")
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(output.Body); err != nil {
		return errors.Wrap(err, "failed to read body from S3 object body")
	}

	key := buf.String()

	// Remove the newline.
	stripe.Key = key[:len(key)-1]
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
