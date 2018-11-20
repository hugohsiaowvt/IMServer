package util

import (
	"math/rand"
	"time"
	"fmt"
	"strconv"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetRandomCode(n int8) string {
	randomCode := ""
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := int8(0); i < n; i ++ {
		randomCode +=fmt.Sprintf("%v", rand.Intn(10))
	}
	return randomCode
}

func GenerateKey(key string) string {
	o_id := key + GetRandomCode(10) + strconv.FormatInt(time.Now().Unix(), 10)
	return MD5(o_id)
}