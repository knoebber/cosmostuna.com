package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/knoebber/comptcheshop/lambda/util"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/order"
)

// OrderResponse is the JSON response for the handler.
type OrderResponse struct {
	OrderID string `json:"orderID"`
	Message string `json:"message"`
	Target  string `json:"target"`
}

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		stripeKey    string
		requestBody  stripe.OrderParams
		responseBody OrderResponse
		o            *stripe.Order
	)

	// TODO only cosmostuna.com
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	if err = json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		response.StatusCode = 400
		return
	}

	if !shippable(requestBody.Shipping.Address.State) {
		responseBody.Message = "We only ship to the lower 48 US states."
		responseBody.Target = "state"
		util.SetResponseBody(&response, &responseBody)
		return
	}

	// TODO Create stage param to choose which key to use.
	stripeKey, err = util.ReadStringKey(util.TestModeKeyName)
	if err != nil {
		response.StatusCode = 500
		return
	}
	stripe.Key = stripeKey

	requestBody.Currency = stripe.String(string(stripe.CurrencyUSD))
	requestBody.Items[0].Type = stripe.String(string(stripe.OrderItemTypeSKU))
	requestBody.Shipping.Address.Country = stripe.String("US")

	var validOrder bool
	for _, offer := range util.BulkOffers {
		// Only accept orders that match one of our offers.
		if offer.Quantity == *requestBody.Items[0].Quantity {
			requestBody.Coupon = offer.CouponID
			validOrder = true
		}
	}
	if !validOrder {
		responseBody.Message = "Order is invalid."
		responseBody.Target = "product-grid"
		util.SetResponseBody(&response, &responseBody)
		return
	}

	o, err = order.New(&requestBody)
	if err != nil {
		responseBody.Message = "failed to POST order to Stripe"
	} else {
		responseBody.OrderID = o.ID
	}

	util.SetResponseBody(&response, &responseBody)
	return
}

func main() {
	lambda.Start(HandleRequest)
}
