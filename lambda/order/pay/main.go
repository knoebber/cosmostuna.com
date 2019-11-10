package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/knoebber/comptcheshop/lambda/util"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/order"
)

type payRequest struct {
	Token   string `json:"token"`
	OrderID string `json:"orderID"`
}

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		stripeKey   string
		requestBody payRequest
	)

	// TODO only cosmostuna.com
	// TODO change over all go error handling to not 200 on errors.
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	if err = json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		response.StatusCode = 400
		return
	}
	// TODO Create stage param to choose which key to use.
	stripeKey, err = util.ReadStringKey(util.TestModeKeyName)
	if err != nil {
		response.StatusCode = 500
		return
	}
	stripe.Key = stripeKey
	params := &stripe.OrderPayParams{}
	if err = params.SetSource(requestBody.Token); err != nil {
		return
	}

	_, err = order.Pay(requestBody.OrderID, params)

	response.StatusCode = 200

	return
}

func main() {
	lambda.Start(HandleRequest)
}
