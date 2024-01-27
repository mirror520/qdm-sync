package qdm

import (
	"encoding/json"
	"time"
)

type AuthData struct {
	StoreUID    string    // 商店專屬代號
	AccessToken string    // JSON Web Token (JWT)
	TokenType   string    // Token 的類型 (Bearer)
	ExpiresIn   time.Time // Token 的效期 (一個小時)
	Message     string    // 成功取得一組 API Access Token
}

func (a *AuthData) UnmarshalJSON(data []byte) error {
	var raw struct {
		StoreUID    string `json:"store_uid"`
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
		Message     string `json:"message"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	a.StoreUID = raw.StoreUID
	a.AccessToken = raw.AccessToken
	a.TokenType = raw.TokenType
	a.ExpiresIn = time.Unix(raw.ExpiresIn, 0).Local()
	a.Message = raw.Message
	return nil
}
