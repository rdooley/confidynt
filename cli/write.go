package cli

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/rdooley/confidynt/service"
	"github.com/rdooley/confidynt/types"
)

var propRe = regexp.MustCompile(`^(\w+)=(.*)$`)

// Write a given config file to a dynamo db table
func Write(table string, path string, ds service.Dynamo, w io.Writer) {
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
		if strings.HasPrefix(line, "#") {
			// Ignore comments in flat files
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
	err = ds.Write(table, config)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Fprintln(w, dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Fprintln(w, dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Fprintln(w, dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Fprintln(w, dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				fmt.Fprintln(w, dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Fprintln(w, dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Fprintln(w, dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Fprintln(w, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Fprintln(w, err.Error())
		}
		return
	}
	fmt.Fprintf(w, "%s written to %s\n", path, table)
}
