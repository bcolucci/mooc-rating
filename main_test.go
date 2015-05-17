package main

import (
	//"fmt"
	"time"
	//"reflect"
	"strconv"
	"net/http"
	"testing"
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

func CurrentTime() int64 {
	return time.Now().Unix()
}
func CurrentTimeStr() string {
	return strconv.FormatInt(CurrentTime(), 10)
}

func AddAuth(r *http.Request) {
	ts := CurrentTimeStr()
	r.Header.Set("ts", ts)
	r.Header.Set("key", string(api.BuildKey(api.Conf.ApiKey, ts)))
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
	rating := &Rating{"MSPI", "Products", "someProductId", 3, 5}
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
	rating := &Rating{"MSPI", "Products", "someProductId", 3, 5}
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
	//rRating := []Rating{}
	//r.DecodeJsonPayload(rRating)
	//fmt.Println(rating)
	//fmt.Println(rRating)
	//if !reflect.DeepEqual(rating, rRating) {
	//	t.Fatal("Saved Rating is not equal")
	//}
}