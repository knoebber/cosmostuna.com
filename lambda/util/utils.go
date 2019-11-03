package util

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

const backendBucket = "cosmostuna-backend"

const (
	// TestModeKeyName is the key in backendBucket that holds the Stripe test mode secret key.
	TestModeKeyName = "/test-secret-stripe-api-key"
)

// BulkOffer maps a quantity of items to a discount.
// If a customer buys that quantity, the discount will be applied.
type BulkOffer struct {
	Quantity int64   `json:"quantity"`
	CouponID *string `json:"couponID"`
}

// BulkOffers are the current deals that the store offers.
var BulkOffers = [4]BulkOffer{
	{
		Quantity: 1,
		CouponID: nil,
	},
	{
		Quantity: 3,
		CouponID: aws.String("3_CAN"),
	},
	{
		Quantity: 12,
		CouponID: aws.String("HALF_CASE"),
	},
	{
		Quantity: 24,
		CouponID: aws.String("FULL_CASE"),
	},
}

// ReadStringKey reads and returns the contents of key in the backend bucket.
func ReadStringKey(key string) (string, error) {
	// Create a S3 client
	session := session.Must(session.NewSession())
	svc := s3.New(session)

	getInput := s3.GetObjectInput{
		Bucket: aws.String(backendBucket),
	}

	// Get the Stripe secret api key
	getInput.Key = aws.String(key)

	output, err := svc.GetObject(&getInput)
	if err != nil {
		return "", errors.Wrap(err, "failed to get stripe api secret S3 object")
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(output.Body); err != nil {
		return "", errors.Wrap(err, "failed to read body from S3 object body")
	}

	value := buf.String()

	// Remove the newline. TODO why is there a newline.
	return value[:len(value)-1], nil
}

// SetResponseBody attempts to marshal body into the Lambda response body.
func SetResponseBody(response *events.APIGatewayProxyResponse, responseBody interface{}) {
	bytes, err := json.Marshal(responseBody)
	if err != nil {
		response.StatusCode = 500
		return
	}

	response.Body = string(bytes)
	response.StatusCode = 200
}
