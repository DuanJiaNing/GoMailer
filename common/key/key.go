package key

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"GoMailer/common/utils"
)

const (
	appKeyPartsCount = 4
)

type AppKey struct {
	UserId     int64
	AppId      int64
	EndpointId int64
	UnixNano   int64
}

func EncodeAppKey(userId, appId, endpointId int64) string {
	keyStr := fmt.Sprintf("%d:%d:%d:%d", userId, appId, endpointId, time.Now().UnixNano())
	key := base64.RawURLEncoding.EncodeToString([]byte(keyStr))
	return key
}

func DecodeAppKey(key string) (*AppKey, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	keyStr := string(bytes)
	parts := strings.Split(keyStr, ":")
	if len(parts) != appKeyPartsCount {
		return nil, errors.New("not a illegal app key")
	}

	userId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	appId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	endpointId, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, err
	}
	unixNano, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return nil, err
	}
	return &AppKey{
		UserId:     userId,
		AppId:      appId,
		EndpointId: endpointId,
		UnixNano:   unixNano,
	}, nil
}

func AppKeyFromRequest(r *http.Request) (*AppKey, error) {
	appKey := r.URL.Query().Get("app_key")
	if utils.IsBlankStr(appKey) {
		return nil, errors.New("app_key is missing")
	}

	ak, err := DecodeAppKey(appKey)
	if err != nil {
		return nil, errors.New("illegal app_key")
	}

	return ak, nil
}
