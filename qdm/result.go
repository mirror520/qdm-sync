package qdm

import (
	"encoding/json"
	"errors"
	"io"
)

const TIME_LAYOUT string = "2006-01-02T15:04:05"

var (
	EOF = io.EOF
)

type ResultPagination struct {
	PageSize   int `json:"page_size"`   // 每頁筆數
	PageNumber int `json:"page_number"` // 從第幾頁開始
	PageCount  int `json:"page_count"`  // 總頁數
}

type Result struct {
	Meta struct {
		Error  bool
		Status int
	}
	Data json.RawMessage
}

func (r *Result) Error() error {
	if !r.Meta.Error {
		return nil
	}

	var data struct {
		Message string
	}

	if err := json.Unmarshal(r.Data, &data); err != nil {
		return err
	}

	return errors.New(data.Message)
}

func (r *Result) AuthData() (data *AuthData, err error) {
	err = json.Unmarshal(r.Data, &data)
	return
}

func (r *Result) OrderCountData() (data *OrderCountData, err error) {
	err = json.Unmarshal(r.Data, &data)
	return
}

func (r *Result) OrderData() (data *OrderData, err error) {
	err = json.Unmarshal(r.Data, &data)
	return
}

func (r *Result) CustomerCountData() (data *CustomerCountData, err error) {
	err = json.Unmarshal(r.Data, &data)
	return
}

func (r *Result) CustomerData() (data *CustomerData, err error) {
	err = json.Unmarshal(r.Data, &data)
	return
}

func (r *Result) CustomerGroupData() (data *CustomerGroupData, err error) {
	err = json.Unmarshal(r.Data, &data)
	return
}
