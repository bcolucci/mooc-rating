package rating

type ApiHostConf struct {
    Addr string
    Port string
}

type ApiDBConf struct {
    Host string
    Name string
    Col string
    Drop bool
}

type ApiConf struct {
    ApiKey string
    Env string
    Host ApiHostConf
    DB ApiDBConf
}

func DevConf() *ApiConf {
    var conf = ApiConf{}
    if err := LoadJSON("./config/dev.json", &conf); err != nil {
        panic(err)
    }
    return &conf
}
