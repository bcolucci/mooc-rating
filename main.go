package main

import (
    "log"
    "net/http"
    "github.com/bcolucci/moocapic-rating/rating"
)

var api *rating.Api

func main() {
    conf := rating.DevConf()
    api = rating.NewApi(conf)
    defer api.Session.Close()
    if api.Conf.Env == "dev" && api.Conf.DB.Drop {
        api.Database.DropDatabase()
    }
    log.Fatal(http.ListenAndServe(api.Conf.Host.Port, api.MakeHandler()))
}
