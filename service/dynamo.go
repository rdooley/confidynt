package service

//go:generate mockgen -destination=dynamo_mock.go -package=service -source dynamo.go

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/rdooley/confidynt/types"
)

// DynamoService provides an interface for dynamodb
type Dynamo interface {
	Read(table, key, value string) (types.Config, error)
	Write(table string, c types.Config) error
}
type dynamo struct {
	region string
}

// NewDynamoService returns a new dynamo service
func NewDynamo() Dynamo {
	return &dynamo{
		region: "us-east-1",
	}
}

// Read a config matching key=value from the given table
func (dynamo) Read(table, key, value string) (c types.Config, err error) {
	// TODO region support
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				S: aws.String(value),
			},
		},
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :v", key)),
		TableName:              aws.String(table),
	}

	result, err := svc.Query(input)
	if err != nil {
		return c, err
	}

	c = types.Config{}
	for _, it := range result.Items {
		for k, v := range it {
			c[k] = *v.S
		}

	}
	if c[key] != value {
		return c, fmt.Errorf("Value mismatch between queried %s and retrived %s for key %s", value, c[key], key)
	}
	return c, err
}

// Write a config to the given table
func (dynamo) Write(table string, c types.Config) (err error) {
	// TODO region support
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	var item = map[string]*dynamodb.AttributeValue{}
	for k, v := range c {
		item[k] = &dynamodb.AttributeValue{S: aws.String(v)}
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(table),
	}

	_, err = svc.PutItem(input)
	return err
}
