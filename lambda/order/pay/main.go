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
	var requestBody payRequest

	// TODO only cosmostuna.com
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	if err = json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		response.StatusCode = 400
		return
	}

	if err = util.SetStripeKey(request.RequestContext.Stage); err != nil {
		response.StatusCode = 500
		return
	}

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
