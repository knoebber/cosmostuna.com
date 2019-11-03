package main

import (
	"encoding/json"
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
		bytes     []byte
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

	bytes, err = json.Marshal(&o)
	if err != nil {
		return
	}

	response.Body = string(bytes)
	response.StatusCode = 200
	return
}

func main() {
	lambda.Start(HandleRequest)
}
