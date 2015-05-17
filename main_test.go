package main

import (
	//"fmt"
	"strconv"
	"testing"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest/test"
)

const HOST = "http://localhost/"

var handler http.Handler

func Setup() {
	conf := DefaultConf()
	api = NewRatingApi(conf)
	handler = api.MakeHandler()
	api.Database.DropDatabase()
}

func AddAuth(r *http.Request) {
	ts := CurrentTimeStr()
	key := string(api.PKeyMiddleware.BuildKey(api.Conf.ApiKey, ts))
	r.Header.Set("ts", ts)
	r.Header.Set("key", key)
}

func CreateRating() *Rating {
	return &Rating{
		Tenant: "MSPI",
		Category: "Products",
		ItemId: "someProductId",
		Rating: 3,
		RatingOn: 5}
}

func TestNoAuth(t *testing.T) {
	Setup()
	defer api.Session.Close()
	r := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", HOST, nil))
	r.CodeIs(500)
	r.ContentTypeIsJson()
}

func TestInvalidAuth(t *testing.T) {
	Setup()
	defer api.Session.Close()
	req := test.MakeSimpleRequest("GET", HOST, nil)
	AddAuth(req)
	req.Header.Set("ts", strconv.FormatInt(CurrentTime() + 1, 10))
	r := test.RunRequest(t, handler, req)
	r.CodeIs(500)
	r.ContentTypeIsJson()
}

func TestGetAllEmpty(t *testing.T) {
	Setup()
	defer api.Session.Close()
	req := test.MakeSimpleRequest("GET", HOST, nil)
	AddAuth(req)
	r := test.RunRequest(t, handler, req)
	r.CodeIs(200)
	r.ContentTypeIsJson()
	r.BodyIs("[]")
}

func TestSave(t *testing.T) {
	Setup()
	defer api.Session.Close()
	rating := CreateRating()
	req := test.MakeSimpleRequest("POST", HOST, rating)
	AddAuth(req)
	r := test.RunRequest(t, handler, req)
	r.CodeIs(200)
	r.ContentTypeIsJson()
}

func TestGetAll(t *testing.T) {
	Setup()
	defer api.Session.Close()
	
	// save one
	rating := CreateRating()
	req := test.MakeSimpleRequest("POST", HOST, rating)
	AddAuth(req)
	r := test.RunRequest(t, handler, req)
	r.CodeIs(200)
	
	// get all
	req = test.MakeSimpleRequest("GET", HOST, nil)
	AddAuth(req)
	r = test.RunRequest(t, handler, req)
	r.CodeIs(200)
	r.ContentTypeIsJson()
	
	//TODO compare returned Rating
	//ratings := []Rating{}
	//r.DecodeJsonPayload(ratings)
	//fmt.Println(ratings)
}