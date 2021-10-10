package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("Unable to get IP", func(t *testing.T) {
		EndpointGetMyIP = "http://127.0.0.1:12345"

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}
	})

	t.Run("Non 200 Response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		EndpointGetMyIP = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err != nil && err.Error() != ErrNon200Response.Error() {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
	})

	t.Run("Unable decode IP", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		EndpointGetMyIP = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid HTTP response")
		}
	})

	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		EndpointGetMyIP = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})

	t.Run("Check succesful 'Hello World' response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		EndpointGetMyIP = ts.URL

		want := "Hello World, 127.0.0.1"

		resp, _ := handler(events.APIGatewayProxyRequest{})
		fmt.Println(resp)
		got := resp.Body
		assert.Equal(t, want, got, "Want and got should be the same.")

	})
}
