package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	gorm.Model
	ID       uint64 `json:"id" ;gorm:"primary_key"`
	Text     string `json:"text" ;gorm:"column:text"`
	NewsId   uint64 `json:"news_id" ;gorm:"column:news_id"`
	Censored bool   `json:"censored" ;gorm:"column:censored"`
}

type CommentChild struct {
	gorm.Model
	ID       uint64 `json:"id" ;gorm:"primary_key"`
	ParentId uint64 `json:"parent_id" ;gorm:"column:parent_id"`
	ChildId  uint64 `json:"child_id" ;gorm:"column:child_id" ;gorm:"CommentRefer"`
	NewsId   uint64 `json:"news_id" ;gorm:"column:news_id"`
}

type CommentUnion struct {
	ID        uint64    `json:"id" ;gorm:"primary_key"`
	Text      string    `json:"text" ;gorm:"column:text"`
	Censored  bool      `json:"censored" ;gorm:"column:censored"`
	ParentId  uint64    `json:"parent_id" ;gorm:"column:parent_id"`
	ChildId   uint64    `json:"child_id" ;gorm:"column:child_id"`
	NewsId    uint64    `json:"news_id" ;gorm:"column:news_id"`
	CreatedAt time.Time `json:"created_at" ;gor:"column:created_at"`
}
