module newsletter-subscribe

go 1.13

require (
	github.com/aws/aws-lambda-go v1.20.0
	github.com/aws/aws-sdk-go v1.35.23
	github.com/go-playground/validator/v10 v10.4.1
	github.com/knoebber/comptche-shop/lambda/newsletter v0.0.0
	github.com/knoebber/comptche-shop/lambda/util v0.0.0
)

replace (
	github.com/knoebber/comptche-shop/lambda/newsletter => ../
	github.com/knoebber/comptche-shop/lambda/util => ../../util/
)
