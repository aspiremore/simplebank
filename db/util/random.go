package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano()) //to make sure at every execution, the seed is different
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)

	}
	return sb.String()
}

func RandomOwner() string  {
	return RandomString(7)
}

func RandomCurrency() string  {
	currency_options := []string{EUR, USD, CAD}
	return currency_options[rand.Intn(len(currency_options))]
}

func RandomMoney() int64 {
	return RandomInt(1,1000)

}

func RandomAccountID() int64 {
	return RandomInt(1000,10000)

}