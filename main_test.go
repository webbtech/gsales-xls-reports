package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestPingHandler(t *testing.T) {

	os.Setenv("Stage", "test")

	var msg string
	t.Run("Successful ping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		r, err := handler(events.APIGatewayProxyRequest{Path: "/"})

		expectedMsg := "Pong"
		msg = extractMessage(r.Body)
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}

func extractMessage(b string) (msg string) {
	var dat map[string]string
	_ = json.Unmarshal([]byte(b), &dat)
	return dat["message"]
}
