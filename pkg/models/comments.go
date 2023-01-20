package models

type Comment struct {
	ID       uint   `json:"id" ;gorm:"primary_key"`
	Text     string `json:"text" ;gorm:"column:text"`
	ParentId uint64 `json:"parent_id" ;gorm:"column:parent_id"`
	ChildId  uint64 `json:"child_id" ;gorm:"column:child_id"`
	NewsId   uint64 `json:"news_id" ;gorm:"column:news_id"`
	Censored bool   `json:"censored" ;gorm:"column:censored"`
}
