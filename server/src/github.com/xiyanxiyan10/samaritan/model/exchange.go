package model

import (
	"github.com/xiyanxiyan10/samaritan/util"
	"time"
)

// Exchange struct
type Exchange struct {
	ID        int64      `gorm:"primary_key" json:"id"`
	UserID    int64      `gorm:"index" json:"userId"`
	Name      string     `gorm:"type:varchar(50)" json:"name"`
	Type      string     `gorm:"type:varchar(50)" json:"type"`
	AccessKey string     `gorm:"type:varchar(200)" json:"accessKey"`
	SecretKey string     `gorm:"type:varchar(200)" json:"secretKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// ListExchange ...
func (user User) ListExchange(size, page int64, order string) (total int64, exchanges []Exchange, err error) {
	_, users, err := user.ListUser(-1, 1, "id")
	if err != nil {
		return
	}
	userIDs := []int64{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	err = DB.Model(&Exchange{}).Where("user_id in (?)", userIDs).Count(&total).Error
	if err != nil {
		return
	}
	if size < 0 {
		err = DB.Where("user_id in (?)", userIDs).Order(util.ToUnderScoreCase(order)).Limit(size).Find(&exchanges).Error
	} else {
		err = DB.Where("user_id in (?)", userIDs).Order(util.ToUnderScoreCase(order)).Limit(size).Offset((page - 1) * size).Find(&exchanges).Error
	}
	return
}
