package orders

type ProductImage struct {
	ImageURL              string `json:"image_url" bson:"image_url"`                             // 圖片CDN網址
	ImageRelativeFilename string `json:"image_relative_filename" bson:"image_relative_filename"` // 相對路徑圖片檔名
	Position              int    `json:"position" bson:"position"`                               // 位置排序 (1=主圖)
}

type ProductVariantOption struct {
}

type ProductVariantInventory struct {
	OptionRefID    string `json:"option_ref_id" bson:"option_ref_id"`       // 選項值編號組合
	SKU            string `json:"sku" bson:"sku"`                           // SKU
	Name           string `json:"name" bson:"name"`                         // 選項值名稱組合
	Quantity       int    `json:"quantity" bson:"quantity"`                 // 組合選項庫存
	MinQuantity    int    `json:"min_quantity" bson:"min_quantity"`         // 低庫存警示量
	OptionImageURL string `json:"option_image_url" bson:"option_image_url"` // 關聯商品圖網址
}

type ProductSpecialOffer struct {
	CustomerGroupID  int     `json:"customer_group_id" bson:"customer_group_id"` // 適用會員群組編號
	Price            int     `json:"price" bson:"price"`                         // 特賣直購價
	DateAvailable    QDMTime `json:"date_available" bson:"date_available"`       // 特賣開始日期
	DateDisavailable QDMTime `json:"date_disavailable" bson:"date_disavailable"` // 特賣截止日期
}

type ProductPurchaseRewards struct {
	CustomerGroupID   int    `json:"customer_group_id" bson:"customer_group_id"`     // 適用會員群組編號
	CustomerGroupName string `json:"customer_group_name" bson:"customer_group_name"` // 會員群組名稱
	Price             int    `json:"price" bson:"price"`                             // 回饋點數
}

type ProductCategory struct {
	CategoryID   int    `json:"category_id" bson:"category_id"`     // 分類編號ID
	CategoryName string `json:"category_name" bson:"category_name"` // 分類名稱
}

type ProductSEOInformation struct {
	MetaDescription string   `json:"meta_description" bson:"meta_description"` // 商品短描述
	MetaKeyword     string   `json:"meta_keyword" bson:"meta_keyword"`         // 商品名稱 & 商品頁標題
	Tag             []string `json:"tag" bson:"tag"`                           // 商品關鍵字標籤
}

type Product struct {
	ProductID                 int                       `json:"product_id" bson:"product_id"`                                   // 商品編號
	URL                       string                    `json:"url" bson:"url"`                                                 // 商品網址
	Name                      string                    `json:"name" bson:"name"`                                               // 商品名稱
	Feature                   string                    `json:"feature" bson:"feature"`                                         // 商品特色
	Description               string                    `json:"description" bson:"description"`                                 // 商品介紹
	Price                     int                       `json:"price" bson:"price"`                                             // 基本售價
	IndependentDeliveryCharge string                    `json:"independent_delivery_charge" bson:"independent_delivery_charge"` // 獨立運費
	SKU                       string                    `json:"sku" bson:"sku"`                                                 // SKU
	EAN                       string                    `json:"ean" bson:"ean"`                                                 // EAN
	Model                     string                    `json:"model" bson:"model"`                                             // 規格型號
	Weight                    float64                   `json:"weight" bson:"weight"`                                           // 重量(kg)
	Length                    float64                   `json:"length" bson:"length"`                                           // 長(cm)
	Width                     float64                   `json:"width" bson:"width"`                                             // 寬(cm)
	Height                    float64                   `json:"height" bson:"height"`                                           // 高(cm)
	TaxClassID                int                       `json:"tax_class_id" bson:"tax_class_id"`                               // 課稅別 (0=應稅, 2=零稅率, 3=免稅)
	DateAvailable             QDMTime                   `json:"date_available" bson:"date_available"`                           // 上架日期
	DateDisavailable          QDMTime                   `json:"date_disavailable" bson:"date_disavailable"`                     // 下架日期
	Status                    bool                      `json:"status" bson:"status"`                                           // 啟用狀態
	NoIndex                   bool                      `json:"no_index" bson:"no_index"`                                       // 隱藏模式
	CustomValue1              string                    `json:"custom_value_1" bson:"custom_value_1"`                           // 自訂資料1
	CustomValue2              string                    `json:"custom_value_2" bson:"custom_value_2"`                           // 自訂資料2
	Images                    []ProductImage            `json:"images" bson:"images"`                                           // 商品圖片集
	VariantsOptions           []ProductVariantOption    `json:"variants_options" bson:"variants_options"`                       // 購物選項組合屬性
	VariantsInventory         []ProductVariantInventory `json:"variants_inventory" bson:"variants_inventory"`                   // 購物選項組合庫存
	SpecialOffer              []ProductSpecialOffer     `json:"special_offer" bson:"special_offer"`                             // 限時特賣集合
	PurchaseRewards           []ProductPurchaseRewards  `json:"purchase_rewards" bson:"purchase_rewards"`                       // 紅利積點回饋
	Categories                []ProductCategory         `json:"categories" bson:"categories"`                                   // 商品分類
	SEOInformation            []ProductSEOInformation   `json:"seo_information" bson:"seo_information"`                         // SEO微資料標記
	DateModified              QDMTime                   `json:"date_modified" bson:"date_modified"`                             // 最近修改時間
}
