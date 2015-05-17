package rating

import (
    "log"
    "net/http"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/ant0ine/go-json-rest/rest"
)

type Api struct {
    rest.Api
    Conf *ApiConf
    PKeyMiddleware *PKeyMiddleware
    Session *mgo.Session
    Database *mgo.Database
    Ratings *mgo.Collection
}

func NewApi(conf *ApiConf) *Api {
    api := &Api{}
    api.Conf = conf
    api.PKeyMiddleware = NewPKeyMiddleware(conf.ApiKey)
    if api.Conf.Env == "dev" {
        api.Use(rest.DefaultDevStack...)  
    } else {
        api.Use(rest.DefaultProdStack...)  
    }
    api.Use(api.PKeyMiddleware)
	api.loadRoutes()
    api.loadDb()
    return api
}

func (api *Api) loadRoutes() {
    router, err := rest.MakeRouter(
        rest.Get("/", api.GetAll),
        rest.Post("/", api.Save))
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
}

func (api *Api) loadDb() {
    session, err := mgo.Dial(api.Conf.DB.Host)
    if err != nil {
        panic(err)
    }
    session.SetMode(mgo.Monotonic, true)
    api.Session = session
    api.Database = api.Session.DB(api.Conf.DB.Name)
    api.Ratings = api.Database.C(api.Conf.DB.Col)
    //TODO set indexes
}

func (api *Api) GetAll(w rest.ResponseWriter, r *rest.Request) {
    var ratings []Rating
	api.Ratings.Find(bson.M{}).All(&ratings)
	if len(ratings) == 0 {
        ratings = make([]Rating, 0)
    }
    w.WriteJson(&ratings)
}

func (api *Api) Save(w rest.ResponseWriter, r *rest.Request) {
    rating := Rating{}
    if err := r.DecodeJsonPayload(&rating); err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
    }
    api.Ratings.Insert(rating)
    w.WriteJson(rating.ID)
}
