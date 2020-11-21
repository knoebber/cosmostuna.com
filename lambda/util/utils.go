package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stripe/stripe-go"
)

const (
	// AWSRegion is where AWS resources are created.
	AWSRegion = "us-west-2"

	// Keys are held in environment variables in Lambda.
	prodStripeKey = "prod_stripe_key"
	testStripeKey = "test_stripe_key"

	prodHookSecret = "prod_webhook_secret"
	testHookSecret = "test_webhook_secret"
)

// BulkOffer maps a quantity of items to a discount.
// If a customer buys that quantity, the discount will be applied.
type BulkOffer struct {
	Quantity int64   `json:"quantity"`
	CouponID *string `json:"couponID"`
}

// BulkOffers are the current deals that the store offers.
var BulkOffers = [3]BulkOffer{
	{
		Quantity: 1,
		CouponID: nil,
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

func envValue(key string) (value string, err error) {
	if value = os.Getenv(key); value == "" {
		err = fmt.Errorf("$%s is not set or empty", key)
	}
	return
}

// SetStripeKey sets the stripe API key.
func SetStripeKey(stageName string) error {
	var (
		stripeKey string
		err       error
	)

	if stageName == "prod" {
		stripeKey, err = envValue(prodStripeKey)
	} else {
		stripeKey, err = envValue(testStripeKey)
	}

	stripe.Key = stripeKey
	return err
}

// HookSecret returns the stripe webhook secret.
func HookSecret(stageName string) (string, error) {
	if stageName == "prod" {
		return envValue(prodHookSecret)
	}

	return envValue(testHookSecret)
}

// SetResponseMessage sets a body like {"message": "..."}.
func SetResponseMessage(response *events.APIGatewayProxyResponse, message string) {
	SetResponseBody(response, map[string]string{"message": message})
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
