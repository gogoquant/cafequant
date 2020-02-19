package model

// TraderExt ...
type TraderExt struct {
	ID       int64  `gorm:"primary_key" json:"id"`
	TraderID int64  `gorm:"index"`
	Type     int64  `gorm:"default:1;not null" json:"type"` //1-int64  2-float63 3-string
	Val      string `gorm:"type:text" json:"val"`
	Content  string `gorm:"type:text" json:"content"`
	Desc     string `gorm:"type:text" json:"desc"`
}

// ListParameters ...
func (user User) ListParameters(traderID int64) (exts []TraderExt, err error) {
	err = DB.Where("trader_id = ?", traderID).Find(&exts).Error
	return
}

// DeleteParameters ...
func (user User) DeleteParameters(traderID int64) (err error) {
	err = DB.Where("trader_id = ?", traderID).Delete(&TraderExt{}).Error
	return
}

// DeleteParameter ...
func (user User) DeleteParameter(parameterID int64) (err error) {
	err = DB.Where("id = ?", parameterID).Delete(&TraderExt{}).Error
	return
}
