package main

import (
	"time"
	"strconv"
)

func CurrentTime() int64 {
	return time.Now().Unix()
}

func CurrentTimeStr() string {
	return strconv.FormatInt(CurrentTime(), 10)
}