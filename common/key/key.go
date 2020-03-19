package key

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"GoMailer/common/utils"
)

func AppKeyFromRequest(r *http.Request) (string, error) {
	const key = "app_key"
	appKey := r.URL.Query().Get(key)
	if !utils.IsBlankStr(appKey) {
		return appKey, nil
	}

	appKey = r.Header.Get(key)
	if !utils.IsBlankStr(appKey) {
		return appKey, nil
	}

	return "", errors.New("app_key is missing")
}

func GenerateKey() string {
	const str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const keyLen = 10
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < keyLen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}
