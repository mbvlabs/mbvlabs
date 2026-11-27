package siteapi

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSearchResponseWireShape(t *testing.T) {
	input := `{"searchParameters":{"q":"term","type":"search","engine":"google"},"organic":[{"title":"title","link":"url","snippet":"text","date":"","position":1,"sitelinks":[{"title":"child","link":"child-url"}]}],"peopleAlsoAsk":[{"question":"why","snippet":"answer","title":"source","link":"url"}],"relatedSearches":[{"query":"other"}],"credits":1}`
	var response APISearchResponse
	if err := json.Unmarshal([]byte(input), &response); err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	var wantValue, gotValue any
	if err := json.Unmarshal([]byte(input), &wantValue); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(wantValue, gotValue) {
		t.Fatalf("response changed: %s", got)
	}
}

func TestScrapeRequestPointerOmission(t *testing.T) {
	var request APIScrapeRequest
	if err := json.Unmarshal([]byte(`{"url":"https://example.com","onlyMainContent":false,"maxAge":0,"formats":["markdown"]}`), &request); err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"url":"https://example.com","formats":["markdown"],"onlyMainContent":false,"maxAge":0}`
	if string(got) != want {
		t.Fatalf("request = %s; want %s", got, want)
	}
}
