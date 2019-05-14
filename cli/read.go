package cli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/rdooley/confidynt/types"
)

// Read a config from dynamo db
func Read(table, key, value string) {
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	config := types.Config{}
	for _, it := range result.Items {
		for k, v := range it {
			config[k] = *v.S
		}

	}
	if config[key] != value {
		// This is bad
		return
	}
	fmt.Printf("%s=%s\n", key, config[key])
	for k, v := range config {
		if k != key {
			fmt.Printf("%s=%s\n", k, v)
		}
	}

}
