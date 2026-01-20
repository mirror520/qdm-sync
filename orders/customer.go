package orders

type AddressInfo struct {
	Postcode string `json:"postcode" bson:"postcode"` // 郵遞區號
	City     string `json:"city" bson:"city"`         // 縣市
	Address  string `json:"address" bson:"address"`   // 地址
}

type CustomerUpgradeHistory struct {
	PreviousCustomerGroupID int     `json:"previous_customer_group_id" bson:"previous_customer_group_id"` // 升等前(舊的)會員群組
	RecentCustomerGroupID   int     `json:"recent_customer_group_id" bson:"recent_customer_group_id"`     // 升等後(新的)會員群組
	EffectiveDate           QDMTime `json:"effective_date" bson:"effective_date"`                         // 升等生效時間
}

// type CartItem Product
type CartItem struct {
}

type RewardRow struct {
	OrderID     int     `json:"order_id" bson:"order_id"`       // 訂單編號
	Description string  `json:"description" bson:"description"` // 紅利項目說明
	Points      int     `json:"points" bson:"points"`           // 紅利積點或折抵 (-負號表示折抵)
	DateAdded   QDMTime `json:"date_added" bson:"date_added"`   // 項目新增時間
}

type Reward struct {
	Total int         `json:"total" bson:"total"` // 累積可用紅利
	Rows  []RewardRow `json:"rows" bson:"rows"`   // 紅利累積與折抵紀錄
}

type Customer struct {
	CustomerID          int                      `json:"customer_id" bson:"customer_id"`                     // 會員編號
	CustomerGroupID     int                      `json:"customer_group_id" bson:"customer_group_id"`         // 會員群組編號
	CountryCode         string                   `json:"country_code" bson:"country_code"`                   // 國別代號 ISO 3166‑1 二字母碼（A2）
	FacebookConnected   int                      `json:"facebook_connected" bson:"facebook_connected"`       // 連動 Facebook 帳號登入 (0=無, 1=使用臉書登入)
	FacebookUserID      string                   `json:"facebook_user_id" bson:"facebook_user_id"`           // Facebook User ID
	LineConnected       int                      `json:"line_connected" bson:"line_connected"`               // 連動 LINE 帳號登入 (0=無, 1=使用LINE登入)
	LineUserID          string                   `json:"line_user_id" bson:"line_user_id"`                   // LINE User ID
	Name                string                   `json:"name" bson:"name"`                                   // 姓名
	Birthday            string                   `json:"birthday" bson:"birthday"`                           // 生日 (YYYY-MM-DD)
	Sex                 int                      `json:"sex" bson:"sex"`                                     // 性別 (1=女, 2=男)
	Email               string                   `json:"email" bson:"email"`                                 // 電子郵件
	Telephone           string                   `json:"telephone" bson:"telephone"`                         // 聯絡手機
	Newsletter          int                      `json:"newsletter" bson:"newsletter"`                       // 訂閱電子報 (0=無, 1=訂閱)
	Approved            int                      `json:"approved" bson:"approved"`                           // 會員資格審核 (0=等待審核, 1=核准)
	DateAdded           QDMTime                  `json:"date_added" bson:"date_added"`                       // 加入日期
	DateModified        QDMTime                  `json:"date_modified" bson:"date_modified"`                 // 資料異動時間
	AddressInfo         []AddressInfo            `json:"address_info" bson:"address_info"`                   // 預設地址(會員中心)
	UpgradeHistoryInfo  []CustomerUpgradeHistory `json:"upgrade_history_info" bson:"upgrade_history_info"`   // 會員升等(群組異動)歷程
	TotalPurchaseAmount int                      `json:"total_purchase_amount" bson:"total_purchase_amount"` // 累積消費金額 (已付款, 不含已取消)
	Cart                []CartItem               `json:"cart" bson:"cart"`                                   // 購物車商品集合
	CartLastModified    QDMTime                  `json:"cart_last_modified" bson:"cart_last_modified"`       // 購物車最近異動時間
	Tags                []string                 `json:"tags" bson:"tags"`                                   // 會員標籤
	CustomValue1        string                   `json:"custom_value_1" bson:"custom_value_1"`               // 自訂資料1
	CustomValue2        string                   `json:"custom_value_2" bson:"custom_value_2"`               // 自訂資料2
	Reward              Reward                   `json:"reward" bson:"reward"`                               // 紅利
}
