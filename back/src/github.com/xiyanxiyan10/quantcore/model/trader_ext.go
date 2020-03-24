package model

// TraderExt ...
type TraderExt struct {
	ID       int64  `gorm:"primary_key" json:"id"`
	BindID   int64  `gorm:"bind_id" json:"bind_id"`
	BindType int64  `gorm:"bind_type" json:"bind_type"`
	Type     int64  `gorm:"default:1;not null" json:"type"` //1-int64  2-float63 3-string
	Val      string `gorm:"type:text" json:"val"`
	Content  string `gorm:"type:text" json:"content"`
	Desc     string `gorm:"type:text" json:"desc"`
}

// ListParameters ...
func (user User) ListParameters(bindID, BindType int64) (extends []TraderExt, err error) {
	var bindQuery = ""
	if BindType == 0 {
		bindQuery = "trader_id = ?"
	} else {
		bindQuery = "algorithm_id = ?"
	}
	err = DB.Where(bindQuery, bindID).Find(&extends).Error
	return
}

// DeleteParameters ...
func (user User) DeleteParameters(bindID, BindType int64) (err error) {
	var bindQuery = ""
	if BindType == 0 {
		bindQuery = "trader_id = ?"
	} else {
		bindQuery = "algorithm_id = ?"
	}
	err = DB.Where(bindQuery, bindID).Delete(&TraderExt{}).Error
	return
}

// DeleteParameter ...
func (user User) DeleteParameter(parameterID int64) (err error) {
	err = DB.Where("id = ?", parameterID).Delete(&TraderExt{}).Error
	return
}
