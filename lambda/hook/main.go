package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/knoebber/comptcheshop/lambda/util"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

const (
	orderSucceeded = "order.payment_succeeded"
	orderUpdated   = "order.updated"

	sender  = "comptcheshop@gmail.com"
	charSet = "UTF-8"
)

type eventResponse struct {
	Message string `json: "message"`
}

func buildEvent(requestBody []byte, signature string) (*stripe.Event, error) {
	hookSecret, err := util.ReadStringKey(util.TestHookSecretName)
	if err != nil {
		return nil, errors.Errorf("reading hook secret, %v", err)
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(requestBody, signature, hookSecret)

	if err != nil {
		return nil, errors.Errorf("constructing event, %v", err)
	}
	return &event, nil
}

// HandleRequest processes a lambda request.
func HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	var (
		responseBody eventResponse
		event        *stripe.Event
		o            stripe.Order
		subject      string
		body         string
	)

	signature, ok := request.Headers["Stripe-Signature"]
	if !ok {
		err = errors.New("stripe signature header missing")
		response.StatusCode = 400
		return
	}

	event, err = buildEvent([]byte(request.Body), signature)
	if err != nil {
		response.StatusCode = 400
		return
	}

	err = json.Unmarshal(event.Data.Raw, &o)
	if err != nil {
		err = fmt.Errorf("parsing webhook JSON, %v", err)
		response.StatusCode = 500
		return
	}

	switch stripe.OrderStatus(o.Status) {
	case stripe.OrderStatusPaid:
		subject = "Your cosmostuna.com order has proccessed"
		body = "Thank you! We’ll send a confirmation when your tuna ships."
	case stripe.OrderStatusFulfilled:
		subject = "Your cosmostuna.com order has shipped"
		body = "Your tuna will arrive soon."
	case stripe.OrderStatusCanceled:
		subject = "Your cosmostuna.com order has been canceled"
		body = "You will be refunded."
	default:
		responseBody.Message = fmt.Sprintf("order status %#v", o.Status)
		util.SetResponseBody(&response, &responseBody)
		return
	}

	if err = sendEmail(o.Email, subject, body, o.Shipping.TrackingNumber); err != nil {
		responseBody.Message = err.Error()
	} else {
		responseBody.Message = fmt.Sprintf("processed %s, order status is %s", event.Type, o.Status)
	}

	util.SetResponseBody(&response, &responseBody)
	return
}

func sendEmail(address string, subject, body, tracking string) error {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{Region: aws.String(util.AWSRegion)})
	if err != nil {
		return fmt.Errorf("starting AWS session in %s, %v", util.AWSRegion, err)
	}

	// Create a SES session.
	svc := ses.New(sess)

	htmlBody := fmt.Sprintf(`<h3>cosmostuna.com</h3><p>%s</p>`, body)
	if tracking != "" {
		htmlBody = fmt.Sprintf(`
%s
<a 
   href="https://tools.usps.com/go/TrackConfirmAction_input?strOrigTrackNum=%s" 
   target="_blank"
>
USPS Tracking Number %s
</a>`, htmlBody, tracking, tracking)

		// For plain text emails.
		body += "\n USPS tracking number: " + tracking
	}

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(address)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}
	if _, err := svc.SendEmail(input); err != nil {
		return fmt.Errorf("failed to send email to %#v, %v", address, err)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
