package main

import (
    "io"
    "log"
    "bytes"
    "net/http"
    "crypto/md5"
    "github.com/ant0ine/go-json-rest/rest"   
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Rating struct {
    Tenant string
    Category string
    ItemId string
    Rating uint
    RatingOn uint
}

type RatingApiConf struct {
    Env string
    Port string
    DbHost string
    DbName string
    ColName string
    DropDb bool
    ApiKey string
}
func DefaultConf() *RatingApiConf {
    return &RatingApiConf{"dev", ":8080", "localhost", "ratings", "ratings", true, "RZPr/WqLfc#?:/^%-j%wl*d2Vg~v$8"}
}

type RatingApi struct {
    rest.Api
    Conf *RatingApiConf
    Session *mgo.Session
    Database *mgo.Database
    Ratings *mgo.Collection
}

type CheckPKeyMiddleware struct {}

var api *RatingApi

func main() {
    conf := DefaultConf() //TODO config files
    api = NewRatingApi(conf)
    defer api.Session.Close()
    if api.Conf.Env == "dev" && api.Conf.DropDb {
        api.Database.DropDatabase()
    }
    log.Fatal(http.ListenAndServe(api.Conf.Port, api.MakeHandler()))
}

func (pkm *CheckPKeyMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
        ts := r.Header.Get("ts")
        if ts == "" {
            rest.Error(w, "ts parameter required", http.StatusInternalServerError)
        }
        key := r.Header.Get("key")
        if key == "" {
            rest.Error(w, "key parameter required", http.StatusInternalServerError)
        }
        expKey := api.BuildKey(api.Conf.ApiKey, ts)
        if !bytes.Equal(expKey, []byte(key)) {
            rest.Error(w, "Authentification failed", http.StatusInternalServerError)
        }
        h(w, r)
	}
}

func NewRatingApi(conf *RatingApiConf) *RatingApi {
    api := &RatingApi{}
    api.Conf = conf
    if api.Conf.Env == "dev" {
        api.Use(rest.DefaultDevStack...)  
    } else {
        api.Use(rest.DefaultProdStack...)  
    }
    api.Use(&CheckPKeyMiddleware{})
	api.loadRoutes()
    api.loadDb()
    return api
}

func (api *RatingApi) BuildKey(pkey string, ts string) []byte {
	h := md5.New()
    io.WriteString(h, pkey)
    io.WriteString(h, ts)
    return h.Sum(nil)
}

func (api *RatingApi) loadRoutes() {
    router, err := rest.MakeRouter(
        rest.Get("/", api.GetAll),
        rest.Post("/", api.Save))
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
}

func (api *RatingApi) loadDb() {
    session, err := mgo.Dial(api.Conf.DbHost)
    if err != nil {
        panic(err)
    }
    session.SetMode(mgo.Monotonic, true)
    api.Session = session
    api.Database = api.Session.DB(api.Conf.DbName)
    api.Ratings = api.Database.C(api.Conf.ColName)
    //TODO set indexes
}

func (api *RatingApi) GetAll(w rest.ResponseWriter, r *rest.Request) {
    var ratings []Rating
	api.Ratings.Find(bson.M{}).All(&ratings)
	if len(ratings) == 0 {
        ratings = make([]Rating, 0)
    }
    w.WriteJson(&ratings)
}

func (api *RatingApi) Save(w rest.ResponseWriter, r *rest.Request) {
    rating := Rating{}
    if err := r.DecodeJsonPayload(&rating); err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
    }
    api.Ratings.Insert(rating)
    w.WriteJson(&rating)
}
