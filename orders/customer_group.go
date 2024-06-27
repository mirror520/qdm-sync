package orders

type CustomerGroup struct {
	CustomerGroupID int    `json:"customer_group_id" bson:"customer_group_id"` // 會員群組編號
	Name            string `json:"name" bson:"name"`                           // 名稱
	Description     string `json:"description" bson:"description"`             // 描述
	Approval        int    `json:"approval" bson:"approval"`                   // 加入會員是否需要審核 (0=免審核)
	EffectivePeriod int    `json:"effective_period" bson:"effective_period"`   // 會員資格效期(年) (0=永久)
	RenewalByAmount int    `json:"renewal_by_amount" bson:"renewal_by_amount"` // 最低購買次數(續會) (0=無限制)
	RenewalByTotal  int    `json:"renewal_by_total" bson:"renewal_by_total"`   // 最低購買金額(續會) (0=無限制)
}
