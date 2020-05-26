package main

import (
	"testing"
)

func TestSearchResponse(t *testing.T) {
	var cases = []struct {
		testname        string
		requestBody     []byte
		expectedSuccess bool
		expectedMessage string
	}{
		{
			testname:        "Search for recipes using an empty request body",
			requestBody:     []byte(`{}`),
			expectedSuccess: false,
			expectedMessage: "Unexpected post body received",
		},
		{
			testname:        "Search for recipes using a valid request body",
			requestBody:     []byte(`[{"ingredient":"mushroom","quantity":100},{"ingredient":"spinach","quantity":100},{"ingredient":"chicken","quantity":400}]`),
			expectedSuccess: true,
			expectedMessage: "OK",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.testname, func(t *testing.T) {
			res := DoSearchRequest(t, c.requestBody)

			if res.Success != c.expectedSuccess {
				t.Errorf("Actual response success %t expected response success %t result mismatch", res.Success, c.expectedSuccess)
			}

			if res.Message != c.expectedMessage {
				t.Errorf("Actual response message %s expected response message %s result mismatch", res.Message, c.expectedMessage)
			}
		})
	}
}
