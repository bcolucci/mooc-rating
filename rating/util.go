package rating

import (
    "os"
	"time"
	"strconv"
    "encoding/json"
)

func CurrentTime() int64 {
	return time.Now().Unix()
}

func CurrentTimeStr() string {
	return strconv.FormatInt(CurrentTime(), 10)
}

func LoadJSON(file string, v interface{}) error {
    p, err := os.Open(file)
    if err != nil {
        return err
    }
    decoder := json.NewDecoder(p)
    return decoder.Decode(v)
}
