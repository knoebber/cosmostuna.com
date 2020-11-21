package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/knoebber/comptche-shop/lambda/util"
	"github.com/stripe/stripe-go/coupon"
)

// CouponResponse is the structure for the response body.
type CouponResponse struct {
	// Coupons maps coupon ID's to  their prices.
	Coupons map[string]int64 `json:"coupons"`

	// Offers describe when coupons are automatically applied.
	Offers [3]util.BulkOffer `json:"offers"`
}

// HandleRequest processes a Lambda request.
// Creates a JSON response body of coupon ID's mapped to their discount amounts.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	if err = util.SetStripeKey(request.RequestContext.Stage); err != nil {
		response.StatusCode = 500
		return
	}

	responseBody := CouponResponse{
		Coupons: make(map[string]int64),
		Offers:  util.BulkOffers,
	}

	i := coupon.List(nil)
	for i.Next() {
		curr := i.Coupon()
		responseBody.Coupons[curr.ID] = curr.AmountOff
	}

	util.SetResponseBody(&response, &responseBody)
	return
}

func main() {
	lambda.Start(HandleRequest)
}
