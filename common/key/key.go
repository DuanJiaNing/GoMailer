package key

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"GoMailer/common/utils"
	"GoMailer/conf"
)

const (
	EPKeyName             = "EPKey"
	ReCaptchaTokenKeyName = "grecaptcha_token"
)

func ReCaptchaKeyFromRequest(r *http.Request) string {
	token := r.URL.Query().Get(ReCaptchaTokenKeyName)
	if !utils.IsBlankStr(token) {
		return token
	}

	token = r.Header.Get(ReCaptchaTokenKeyName)
	if !utils.IsBlankStr(token) {
		return token
	}

	return getFromForm(r, ReCaptchaTokenKeyName)
}

func getFromForm(r *http.Request, keyName string) string {
	for k, vs := range r.Form {
		if k == keyName {
			return vs[0]
		}
	}

	return ""
}

func EPKeyFromRequest(r *http.Request) string {
	appKey := r.URL.Query().Get(EPKeyName)
	if !utils.IsBlankStr(appKey) {
		return appKey
	}

	appKey = r.Header.Get(EPKeyName)
	if !utils.IsBlankStr(appKey) {
		return appKey
	}

	return getFromForm(r, EPKeyName)
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

func VerifyReCaptcha(token string) (bool, error) {
	if utils.IsBlankStr(token) {
		return false, errors.New("reCaptcha token is empty")
	}

	const addr = "https://recaptcha.net/recaptcha/api/siteverify?secret=%s&response=%s"
	resp, err := http.Get(fmt.Sprintf(addr, conf.ReCaptchaSecret(), token))
	if err != nil {
		return false, err
	}
	m := struct {
		ChallengeTs time.Time `json:"challenge_ts"`
		Score       float32   `json:"score"`
		Hostname    string    `json:"hostname"`
		Success     bool      `json:"success"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return false, err
	}

	return m.Success, nil
}
