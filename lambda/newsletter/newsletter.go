package newsletter

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	tableName   = "cosmos-tuna-newsletter"
	tokenLength = 32
)

// Subscription represents a dynamodb object.
// Email is the primary key.
type Subscription struct {
	Email          string  `json:"email" dynamodbav:"email"`
	TestMode       bool    `json:"test_mode" dynamodbav:"test_mode"`
	Token          string  `json:"-" dynamodbav:"token"`
	CreatedAt      string  `json:"-" dynamodbav:"created_at"`
	ConfirmedAt    *string `json:"-" dynamodbav:"confirmed_at"`
	UnsubscribedAt *string `json:"-" dynamodbav:"unsubscribed_at"`
}

// Subscribe subscribes an email for the newsletter.
// Returns an error wrapping dynamodb.ConditionalCheckFailedException if email already exists.
func Subscribe(svc *dynamodb.DynamoDB, email string, testMode bool) error {
	token, err := Token()
	if err != nil {
		return err
	}

	s := &Subscription{
		Email:     email,
		CreatedAt: time.Now().String(),
		TestMode:  testMode,
		Token:     token,
	}

	item, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		return fmt.Errorf("marshaling newsletter subscription for %q: %w", s.Email, err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(email)"),
		Item:                item,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("putting newsletter subscription for %q: %w", s.Email, err)
	}

	return nil
}

func primaryKey(email string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{"email": {S: aws.String(email)}}
}

// SetUnsubscribedAt sets or removes a subscription's unscribed_at field.
// When unsubscribed_at is nil creates a new token.
func SetUnsubscribedAt(svc *dynamodb.DynamoDB, email string, unsubscribedAt *time.Time) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key:       primaryKey(email),
	}

	if unsubscribedAt == nil {
		token, err := Token()
		if err != nil {
			return err
		}

		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(token),
			},
		}
		input.UpdateExpression = aws.String("SET token = :t REMOVE unsubscribed_at")
	} else {
		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":u": {
				S: aws.String(unsubscribedAt.String()),
			},
		}
		input.UpdateExpression = aws.String("SET unsubscribed_at = :u")
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("updating newsletter unsubscribed time for %q: %w", email, err)
	}

	return nil
}

// SetConfirmedAt sets a subscription's confirmed time.
func SetConfirmedAt(svc *dynamodb.DynamoDB, email string, confirmedAt time.Time) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key:       primaryKey(email),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(confirmedAt.String()),
			},
		},
		UpdateExpression: aws.String("SET confirmed_at = :c REMOVE token"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("updating newsletter unsubscribed time for %q: %w", email, err)
	}

	return nil
}

// Get returns a newsletter subscription if it exists.
// Returns nil, nil when email is not found.
func Get(svc *dynamodb.DynamoDB, email string) (*Subscription, error) {
	var res Subscription

	output, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       primaryKey(email),
	})
	if err != nil {
		return nil, fmt.Errorf("getting dynamodb entry for %q: %w", email, err)
	}

	if output.Item == nil {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling newsletter subscription for %q: %w", email, err)
	}

	return &res, nil
}

func randomBytes(n int) ([]byte, error) {
	buff := make([]byte, n)

	if _, err := io.ReadFull(rand.Reader, buff); err != nil {
		return nil, err
	}

	return buff, nil
}

// Token generates a random string for use in email subscriptions.
func Token() (string, error) {
	buff, err := randomBytes(32)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buff), nil
}
