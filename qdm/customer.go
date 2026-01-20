package qdm

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/mirror520/qdm-sync/orders"
)

type CustomerCountData struct {
	Count int64 `json:"count"` // 會員筆數
}

func (d *CustomerCountData) UnmarshalJSON(data []byte) error {
	var raw struct {
		Count json.Number `json:"count"` // 會員筆數
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

type CustomerData struct {
	Count          int               `json:"count"`           // 擷取會員數
	TotalCount     int               `json:"total_count"`     // 符合會員數
	SearchCriteria ResultPagination  `json:"search_criteria"` // 分頁參數
	Result         []orders.Customer `json:"result"`          // 會員集合
}

type CustomerParams struct {
	CreatedAtMin time.Time // 起始時間 (required)
	CreatedAtMax time.Time // 結束時間 (required)
	PageSize     int       // 每頁筆數
	PageNumber   int       // 從第幾頁開始
}

func (p *CustomerParams) Values() url.Values {
	values := make(url.Values)
	values.Set("created_at_min", p.CreatedAtMin.Format(TIME_LAYOUT))
	values.Set("created_at_max", p.CreatedAtMax.Format(TIME_LAYOUT))
	values.Set("page_size", strconv.Itoa(p.PageSize))
	values.Set("page_number", strconv.Itoa(p.PageNumber))

	return values
}
