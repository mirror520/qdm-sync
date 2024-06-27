package qdm

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/mirror520/qdm-sync/orders"
)

type OrderCountData struct {
	Count int64 `json:"count"` // 訂單筆數
}

func (d *OrderCountData) UnmarshalJSON(data []byte) error {
	var raw struct {
		Count json.Number `json:"count"` // 訂單筆數
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	count, err := raw.Count.Int64()
	if err != nil {
		return err
	}

	d.Count = count
	return nil
}

type OrderData struct {
	Count          int            `json:"count"`           // 擷取訂單數
	TotalCount     int            `json:"total_count"`     // 總訂單數
	SearchCriteria SearchCriteria `json:"search_criteria"` // 分頁參數
	Result         []orders.Order `json:"result"`          // 訂單集合
}

type SearchCriteria struct {
	PageSize   int `json:"page_size"`   // 每頁筆數
	PageNumber int `json:"page_number"` // 從第幾頁開始
	PageCount  int `json:"page_count"`  // 總頁數
}

type OrderParams struct {
	CreatedAtMin time.Time // 起始時間 (required)
	CreatedAtMax time.Time // 結束時間 (required)
	PageSize     int       // 每頁筆數
	PageNumber   int       // 從第幾頁開始
}

func (p *OrderParams) Values() url.Values {
	values := make(url.Values)
	values.Set("created_at_min", p.CreatedAtMin.Format(TIME_LAYOUT))
	values.Set("created_at_max", p.CreatedAtMax.Format(TIME_LAYOUT))
	values.Set("page_size", strconv.Itoa(p.PageSize))
	values.Set("page_number", strconv.Itoa(p.PageNumber))

	return values
}

type OrderOption interface {
	apply(*OrderParams)
}

func WithPageSize(size int) OrderOption {
	return pageSizeOption(size)
}

type pageSizeOption int

func (opt pageSizeOption) apply(p *OrderParams) {
	p.PageSize = int(opt)
}

func WithPageNumber(num int) OrderOption {
	return pageNumberOption(num)
}

type pageNumberOption int

func (opt pageNumberOption) apply(p *OrderParams) {
	p.PageNumber = int(opt)
}
