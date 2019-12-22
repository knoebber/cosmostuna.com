module coupon

go 1.13

require (
	github.com/aws/aws-lambda-go v1.13.2
	github.com/knoebber/comptcheshop/lambda/util v0.0.0
	github.com/stripe/stripe-go v66.1.1+incompatible
)

replace github.com/knoebber/comptcheshop/lambda/util => ../util
