package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type QDMTime time.Time

func (t *QDMTime) UnmarshalJSON(data []byte) error {
	var tsStr string
	if err := json.Unmarshal(data, &tsStr); err != nil {
		return err
	}

	ts, _ := time.ParseInLocation("2006-01-02T15:04:05", tsStr, time.Local)

	*t = QDMTime(ts)

	return nil
}

func (t QDMTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	dt := time.Time(t)
	return bson.MarshalValue(dt)
}

type Order struct {
	OrderID                int                `json:"order_id" bson:"order_id"`                                 // 訂單編號
	BuyButtonID            string             `json:"bb_id" bson:"bb_id"`                                       // Buy Button 專屬編號
	OrderStatus            int                `json:"order_status" bson:"order_status"`                         // 訂單狀態 (1=待處理, 2=處理中, 3=已配送, 4=已取消, 5=已完成, 6=調貨中。或店家另新增的狀態)
	DateAdded              QDMTime            `json:"date_added" bson:"date_added"`                             // 購買日期
	OrderItems             []OrderItem        `json:"order_items" bson:"order_items"`                           // 購買商品集合
	OrderSubtotals         []OrderSubtotals   `json:"order_subtotals" bson:"order_subtotals"`                   // 訂單小計項目
	Total                  float64            `json:"total" bson:"total"`                                       // 訂單總金額
	CurrencyCode           string             `json:"currency_code" bson:"currency_code"`                       // 幣別
	CouponCodeUsed         string             `json:"coupon_code_used" bson:"coupon_code_used"`                 // 使用優惠券代碼
	RewardUsed             int                `json:"reward_used" bson:"reward_used"`                           // 紅利折抵
	RewardBonus            int                `json:"reward_bonus" bson:"reward_bonus"`                         // 紅利獲得
	PromotionDiscountTotal int                `json:"promotion_discount_total" bson:"promotion_discount_total"` // 折扣總計
	ShippingFee            int                `json:"shipping_fee" bson:"shipping_fee"`                         // 運費
	PaymentCode            string             `json:"payment_code" bson:"payment_code"`                         // 付款方式代號
	PaymentMethod          string             `json:"payment_method" bson:"payment_method"`                     // 付款方式說明
	ShippingCode           string             `json:"shipping_code" bson:"shipping_code"`                       // 配送方式代號
	ShippingMethod         string             `json:"shipping_method" bson:"shipping_method"`                   // 配送方式說明
	CustomerID             int                `json:"customer_id" bson:"customer_id"`                           // 會員 ID
	CustomerGroupID        int                `json:"customer_group_id" bson:"customer_group_id"`               // 會員群組 ID
	RegularCustomer        int                `json:"regular_customer" bson:"regular_customer"`                 // 購買次數
	PaymentName            string             `json:"payment_name" bson:"payment_name"`                         // 購買人姓名
	PaymentEmail           string             `json:"payment_email" bson:"payment_email"`                       // 購買人電子郵件信箱
	PaymentTelephone       string             `json:"payment_telephone" bson:"payment_telephone"`               // 購買人手機號碼
	Comment                string             `json:"comment" bson:"comment"`                                   // 購買人留言備註
	CustomerIPAddress      string             `json:"customer_ip_address" bson:"customer_ip_address"`           // 連線 IP
	CustomerUserAgent      string             `json:"customer_user_agent" bson:"customer_user_agent"`           // 連線裝置 User-Agent
	ShippingName           string             `json:"shipping_name" bson:"shipping_name"`                       // 收件人姓名
	ShippingTelephone      string             `json:"shipping_telephone" bson:"shipping_telephone"`             // 收件人手機號碼
	ShippingCountry        string             `json:"shipping_country" bson:"shipping_country"`                 // 配送國家
	ShippingPostcode       int                `json:"shipping_postcode" bson:"shipping_postcode"`               // 郵遞區號
	ShippingZone           string             `json:"shipping_zone" bson:"shipping_zone"`                       // 配送縣市
	ShippingAddress        string             `json:"shipping_address" bson:"shipping_address"`                 // 配送地址
	PickupPoint            string             `json:"pickup_point" bson:"pickup_point"`                         // 取貨門市店號（或服務編號）
	ConvenientStoreType    string             `json:"convenient_store_type" bson:"convenient_store_type"`       // 超商類型
	ConvenientStoreNo      string             `json:"convenient_store_no" bson:"convenient_store_no"`           // 超商門市店號 (UNIMART=7-11, FAMI=全家, HILIFE=萊爾富)
	ConvenientStoreName    string             `json:"convenient_store_name" bson:"convenient_store_name"`       // 超商門市名稱
	ConvenientStoreAddress string             `json:"convenient_store_address" bson:"convenient_store_address"` // 超商門市地址
	ECPayAllPayLogisticsID string             `json:"ECPay_AllPayLogisticsID" bson:"ECPay_AllPayLogisticsID"`   // 綠界科技物流交易編號
	ECPayCVSPaymentNo      string             `json:"ECPay_CVSPaymentNo" bson:"ECPay_CVSPaymentNo"`             // 綠界科技寄貨編號
	ECPayCVSValidationNo   string             `json:"ECPay_CVSValidationNo" bson:"ECPay_CVSValidationNo"`       // 綠界科技驗證碼
	ShippingDeliverPeriod  string             `json:"shipping_deliver_period" bson:"shipping_deliver_period"`   // 指定到貨時段
	ShippingDeliverOndate  string             `json:"shipping_deliver_ondate" bson:"shipping_deliver_ondate"`   // 指定到貨日期
	ShippingStatus         string             `json:"shipping_status" bson:"shipping_status"`                   // 配送狀態 (PENDING=等待出貨, SHIPPED=已出貨, PICKREADY=貨到門市, DELIVERED=已取貨, EXCEPTION=配送失敗退回, ABANDONED=七天未取)
	TrackingNumber         string             `json:"tracking_number" bson:"tracking_number"`                   // 包裹追蹤碼
	PaymentStatus          string             `json:"payment_status" bson:"payment_status"`                     // 付款狀態 (PAID=已付款, UNPAID=尚未付款)
	PaymentTime            QDMTime            `json:"payment_time" bson:"payment_time"`                         // 付款時間
	PaymentBanks           string             `json:"payment_banks" bson:"payment_banks"`                       // 付款來源
	PaymentAccount         string             `json:"payment_account" bson:"payment_account"`                   // 付款帳號
	InvoiceNumber          string             `json:"invoice_number" bson:"invoice_number"`                     // 發票號碼
	InvoiceType            int                `json:"invoice_type" bson:"invoice_type"`                         // 發票類型 (0=不開發票, 2=二聯式, 3=三聯式)
	InvoiceRandomNumber    string             `json:"invoice_randomnumber" bson:"invoice_randomnumber"`         // 發票隨機碼
	InvoiceVatNumber       string             `json:"invoice_vat_number" bson:"invoice_vat_number"`             // 開立統編
	InvoiceCarrierNum      string             `json:"invoice_carrier_num" bson:"invoice_carrier_num"`           // 載具編號
	InvoiceLovecode        string             `json:"invoice_lovecode" bson:"invoice_lovecode"`                 // 捐贈愛心碼
	InvoiceTitle           string             `json:"invoice_title" bson:"invoice_title"`                       // 發票抬頭
	InvoiceAddress         string             `json:"invoice_address" bson:"invoice_address"`                   // 發票地址
	InvoiceCreatedAt       QDMTime            `json:"invoice_created_at" bson:"invoice_created_at"`             // 發票開立時間
	MaRID                  string             `json:"ma_rid" bson:"ma_rid"`                                     // 美安 rid
	MaClickID              string             `json:"ma_click_id" bson:"ma_click_id"`                           // 美安 click_id
	LineECID               string             `json:"line_ecid" bson:"line_ecid"`                               // LINE 購物 ecid
	ECID                   string             `json:"ec_id" bson:"ec_id"`                                       // 導購追蹤識別碼
	UTMSource              string             `json:"utm_source" bson:"utm_source"`                             // utm_source
	UTMMedium              string             `json:"utm_medium" bson:"utm_medium"`                             // utm_medium
	UTMCampaign            string             `json:"utm_campaign" bson:"utm_campaign"`                         // utm_campaign
	UTMContent             string             `json:"utm_content" bson:"utm_content"`                           // utm_content
	UTMTerm                string             `json:"utm_term" bson:"utm_term"`                                 // utm_term
	AffiliateID            string             `json:"affiliate_id" bson:"affiliate_id"`                         // 分銷合作 KOL 專屬編號
	DateModified           QDMTime            `json:"date_modified" bson:"date_modified"`                       // 最近修改時間
	ReturnRequest          int                `json:"return_request" bson:"return_request"`                     // 此單是否申請退換貨 (0=無, 1=是)
	ReturnID               int                `json:"return_id" bson:"return_id"`                               // 退換貨單號
	ReturnStatusID         int                `json:"return_status_id" bson:"return_status_id"`                 // 退換貨處理狀態 (1=待處理, 2=等待商品到貨, 3=已完成)
	ReturnName             string             `json:"return_name" bson:"return_name"`                           // 退換貨申請人
	ReturnEmail            string             `json:"return_email" bson:"return_email"`                         // 退換貨申請人電子郵件
	ReturnTelephone        string             `json:"return_telephone" bson:"return_telephone"`                 // 退換貨申請人手機號碼
	ReturnComment          string             `json:"return_comment" bson:"return_comment"`                     // 退換貨申請人留言
	ReturnAddress          string             `json:"return_address" bson:"return_address"`                     // 退換貨申請收貨地址
	ReturnReturnBank       string             `json:"return_return_bank" bson:"return_return_bank"`             // 退換貨申請退款銀行
	ReturnDateAdded        QDMTime            `json:"return_date_added" bson:"return_date_added"`               // 退換貨申請時間
	ReturnItems            []OrderReturnItems `json:"return_items" bson:"return_items"`                         // 退換貨商品
}

type OrderItem struct {
	ProductID int                  `json:"product_id" bson:"product_id"` // 商品編號
	Name      string               `json:"name" bson:"name"`             // 商品名稱
	Options   []OrderProductOption `json:"options" bson:"options"`       // 購物選項
	SKU       string               `json:"sku" bson:"sku"`               // 單一商品 SKU
	Price     int                  `json:"price" bson:"price"`           // 單價
	Quantity  int                  `json:"quantity" bson:"quantity"`     // 數量
	Total     int                  `json:"total" bson:"total"`           // 小計
	Giveaway  int                  `json:"giveaway" bson:"giveaway"`     // 是否為贈品
}

type OrderProductOptionID string

func (id *OrderProductOptionID) UnmarshalJSON(data []byte) error {
	var i any
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch val := i.(type) {
	case float64:
		str := fmt.Sprintf("%v", val)
		*id = OrderProductOptionID(str)

	case string:
		*id = OrderProductOptionID(val)

	default:
		return errors.New("invalid type")
	}

	return nil
}

func (id OrderProductOptionID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	i := string(id)
	return bson.MarshalValue(i)
}

type OrderProductOption struct {
	ID    OrderProductOptionID `json:"id" bson:"id"`       // 選項值編號
	Name  string               `json:"name" bson:"name"`   // 選項名稱
	Value string               `json:"value" bson:"value"` // 選項值
}

type OrderSubtotals struct {
	Code  string `json:"code" bson:"code"`   // 項目類型
	Name  string `json:"name" bson:"name"`   // 項目名稱
	Value int    `json:"value" bson:"value"` // 小計金額
}

type OrderReturnItemsProductID string

func (id *OrderReturnItemsProductID) UnmarshalJSON(data []byte) error {
	var i any
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch val := i.(type) {
	case float64:
		str := fmt.Sprintf("%v", val)
		*id = OrderReturnItemsProductID(str)

	case string:
		*id = OrderReturnItemsProductID(val)

	default:
		return errors.New("invalid type")
	}

	return nil
}

func (id OrderReturnItemsProductID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	i := string(id)
	return bson.MarshalValue(i)
}

type OrderReturnItems struct {
	ProductID        OrderReturnItemsProductID `json:"product_id" bson:"product_id"`                 // 商品編號
	Name             string                    `json:"name" bson:"name"`                             // 商品名稱
	SKU              string                    `json:"sku" bson:"sku"`                               // SKU
	Quantity         int                       `json:"quantity" bson:"quantity"`                     // 數量
	Price            int                       `json:"price" bson:"price"`                           // 價格
	ReturnReasonID   int                       `json:"return_reason_id" bson:"return_reason_id"`     // 退換貨申請原因代號
	ReturnReasonName string                    `json:"return_reason_name" bson:"return_reason_name"` // 退換貨申請原因
}
