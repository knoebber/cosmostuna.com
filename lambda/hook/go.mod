module hook

go 1.13

replace github.com/knoebber/comptcheshop/lambda/util => ../util

require (
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.25.19
	github.com/knoebber/comptcheshop/lambda/util v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.8.1
	github.com/stripe/stripe-go v67.3.0+incompatible
)
