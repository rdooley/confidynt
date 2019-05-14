package cli

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/rdooley/confidynt/service"
)

// Read a config from dynamo db
func Read(table, key, value string, ds service.Dynamo, w io.Writer) {
	config, err := ds.Read(table, key, value)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Fprintln(w, dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Fprintln(w, dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
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

	fmt.Fprintf(w, "%s=%s\n", key, config[key])
	for k, v := range config {
		if k != key {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}

}
