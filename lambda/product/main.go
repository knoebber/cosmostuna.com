package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/knoebber/comptche-shop/lambda/util"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"
)

type productResponse struct {
	Name    string `json:"name"`
	SKUList []SKU  `json:"SKUList"`
}

// SKU is a shop keeping unit for a product.
// Provides a leaner structure for stripe.SKU information.
type SKU struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
	Price    int64  `json:"price"`
	Image    string `json:"imageURL"`
	Flavor   string `json:"flavor"`
}

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var p *stripe.Product

	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	productID, ok := request.PathParameters["productID"]
	if !ok {
		response.StatusCode = 400
		err = errors.New("path parameter productID is required")
		return
	}

	if err = util.SetStripeKey(request.RequestContext.Stage); err != nil {
		response.StatusCode = 500
		return
	}

	p, err = product.Get(productID, nil)
	if err != nil {
		response.StatusCode = 500
		return
	}

	responseBody := productResponse{
		Name:    p.Name,
		SKUList: []SKU{},
	}

	i := sku.List(&stripe.SKUListParams{Product: stripe.String(productID)})
	for i.Next() {
		curr := i.SKU()
		responseBody.SKUList = append(responseBody.SKUList, SKU{
			ID:       curr.ID,
			Quantity: curr.Inventory.Quantity,
			Price:    curr.Price,
			Image:    curr.Image,
			Flavor:   curr.Attributes["flavor"],
		})
	}

	util.SetResponseBody(&response, &responseBody)
	return
}

func main() {
	lambda.Start(HandleRequest)
}
