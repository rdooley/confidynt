package main

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
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app    = kingpin.New("confidynt", "A command-line application for 12 factor config dynamo db management")
	dryrun = app.Flag("dryrun", "Enable dryrun mode.").Bool()
	table  = app.Flag("table", "Table name").Required().String()

	read      = app.Command("read", "Read a config from dynamo")
	readKey   = read.Arg("key", "Query key").Required().String()
	readValue = read.Arg("value", "Query value").Required().String()

	write     = app.Command("write", "Write a config file to dynamo")
	writeFile = write.Arg("config", "Config file to write").ExistingFile()

	propRe = regexp.MustCompile(`^(\w+)=(.*)$`)
)

type Config map[string]string

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

	config := Config{}
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

func Write(table string, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	config := Config{}
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

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Read
	case read.FullCommand():
		Read(*table, *readKey, *readValue)

	// Write
	case write.FullCommand():
		Write(*table, *writeFile)
	}
}
