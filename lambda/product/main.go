package main

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"
)

// TODO use stage parameter for testMode.
const (
	testMode        = true
	testModeKeyName = "/test-secret-stripe-api-key"
	bucketName      = "cosmostuna-backend"
)

// ProductResponse is the JSON response that HandleRequest responds with.
type ProductResponse struct {
	Name    string        `json:"name"`
	SKUList []*stripe.SKU `json:"SKUList"`
}

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		bytes []byte
		p     *stripe.Product
	)

	// TODO only cosmostuna.com
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	productID, ok := request.PathParameters["productID"]
	if !ok {
		response.StatusCode = 400
		err = errors.New("path parameter productID is required")
		return
	}

	if err = initStripe(); err != nil {
		response.StatusCode = 500
		return
	}

	p, err = product.Get(productID, nil)
	if err != nil {
		response.StatusCode = 500
		return
	}

	responseBody := ProductResponse{
		Name:    p.Name,
		SKUList: []*stripe.SKU{},
	}

	i := sku.List(&stripe.SKUListParams{Product: stripe.String(productID)})
	for i.Next() {
		responseBody.SKUList = append(responseBody.SKUList, i.SKU())
	}

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
