package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	doc, err := goquery.NewDocument("https://www.mozilla.org/en-US/firefox/releases/")
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	version := ""
	doc.Find("html").Each(func(_ int, s *goquery.Selection) {
		version, _ = s.Attr("data-latest-firefox")
	})

	body, err := json.Marshal(map[string]interface{}{
		"response_type": "in_channel",
		"text":          fmt.Sprintf("latest firefox is %v", version),
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
