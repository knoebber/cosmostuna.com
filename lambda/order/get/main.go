package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/knoebber/comptcheshop/lambda/util"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/order"
)

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		stripeKey string
		o         *stripe.Order
	)

	// TODO only cosmostuna.com
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	orderID, ok := request.PathParameters["orderID"]
	if !ok {
		response.StatusCode = 400
		err = errors.New("path parameter orderID is required")
		return
	}

	// TODO Create stage param to choose which key to use.
	stripeKey, err = util.ReadStringKey(util.TestModeKeyName)
	if err != nil {
		response.StatusCode = 500
		return
	}
	stripe.Key = stripeKey
	o, err = order.Get(orderID, nil)
	if err != nil {
		response.StatusCode = 500
		return
	}

	util.SetResponseBody(&response, o)
	return
}

func main() {
	lambda.Start(HandleRequest)
}
