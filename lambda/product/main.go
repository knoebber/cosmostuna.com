package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/knoebber/comptcheshop/lambda/util"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"
)

// ProductResponse is the structure for the response body.
type ProductResponse struct {
	Name    string `json:"name"`
	SKUList []SKU  `json:"SKUList"`
}

// SKU is a shop keeping unit for a product.
// Provides a leaner structure for stripe.SKU information.
type SKU struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quanitity"`
	Price    int64  `json:"price"`
	Image    string `json:"imageURL"`
	Flavor   string `json:"flavor"`
}

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		stripeKey string
		bytes     []byte
		p         *stripe.Product
	)

	// TODO only cosmostuna.com
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	productID, ok := request.PathParameters["productID"]
	if !ok {
		response.StatusCode = 400
		err = errors.New("path parameter productID is required")
		return
	}

	// TODO Create stage param to choose which key to use.
	stripeKey, err = util.ReadStringKey(util.TestModeKeyName)
	if err != nil {
		response.StatusCode = 500
		return
	}
	stripe.Key = stripeKey

	p, err = product.Get(productID, nil)
	if err != nil {
		response.StatusCode = 500
		return
	}

	responseBody := ProductResponse{
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

	bytes, err = json.Marshal(&responseBody)
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
