package qdm

import "github.com/mirror520/qdm-sync/orders"

type CustomerGroupData struct {
	Count  int                    `json:"count"`  // 會員群組數
	Result []orders.CustomerGroup `json:"result"` // 會員群組集合
}
