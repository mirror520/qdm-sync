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
		ExpiresIn   struct {
			Date         string `json:"date"`
			TimezoneType int    `json:"timezone_type"`
			Timezone     string `json:"timezone"`
		} `json:"expires_in"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	a.StoreUID = raw.StoreUID
	a.AccessToken = raw.AccessToken
	a.TokenType = raw.TokenType
	a.Message = raw.Message

	// 解析新的時間格式
	if raw.ExpiresIn.Date != "" {
		loc, err := time.LoadLocation(raw.ExpiresIn.Timezone)
		if err != nil {
			loc = time.Local // 如果時區解析失敗，使用本地時區
		}

		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05.000000", raw.ExpiresIn.Date, loc)
		if err != nil {
			return err
		}

		a.ExpiresIn = parsedTime.Local()
	}

	return nil
}
