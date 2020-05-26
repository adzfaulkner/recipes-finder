package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func DoSearchRequest(t *testing.T, requestBody []byte) SearchResponse {
	resp, err := http.Post("http://test_server:8080/search", "application/json", bytes.NewReader(requestBody))

	if err != nil {
		t.Fatalf("http post error: %+v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("body read error: %+v", err)
	}

	var res SearchResponse
	err = json.Unmarshal(body, &res)

	if err != nil {
		t.Fatalf("body json decode error: %+v", err)
	}

	return res
}

func DoResultsRequest(t *testing.T, id string) ResultResponse {
	resp, err := http.Get(fmt.Sprintf("http://test_server:8080/results?id=%s", id))

	if err != nil {
		t.Fatalf("http get error: %+v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("body read error: %+v", err)
	}

	var res ResultResponse
	err = json.Unmarshal(body, &res)

	if err != nil {
		t.Fatalf("body unmarshal error: %+v", err)
	}

	return res
}
