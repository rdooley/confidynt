package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/rdooley/confidynt/types"
)

var propRe = regexp.MustCompile(`^(\w+)=(.*)$`)

// Write a given config file to a dynamo db table
func Write(table string, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	config := types.Config{}
	var key string
	var value string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// Ignore blank lines
			continue
		}
		if propRe.MatchString(line) {
			if key != "" {
				config[key] = value
			}
			matches := propRe.FindStringSubmatch(line)
			key = matches[1]
			value = matches[2]
		} else {
			// Continuation of a previous property "^  indented..."
			value += "\n" + strings.TrimRight(line, " \t")
		}
	}
	// Catch final prop
	if key != "" {
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	var item = map[string]*dynamodb.AttributeValue{}
	for k, v := range config {
		item[k] = &dynamodb.AttributeValue{S: aws.String(v)}
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(table),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				fmt.Println(dynamodb.ErrCodeTransactionConflictException, aerr.Error())
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
	fmt.Printf("%s written to %s\n", path, table)
}
