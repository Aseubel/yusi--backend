package model

import (
	"time"
)

// User 用户模型
type User struct {
	UserId     string    `gorm:"column:user_id;primaryKey" json:"userId"`
	UserName   string    `gorm:"column:user_name" json:"userName"`
	Password   string    `gorm:"column:password" json:"-"`
	Email      string    `gorm:"column:email" json:"email"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (User) TableName() string {
	return "user"
}

// Diary 日记模型
type Diary struct {
	DiaryId    string    `gorm:"column:diary_id;primaryKey" json:"diaryId"`
	UserId     string    `gorm:"column:user_id;index" json:"userId"`
	Title      string    `gorm:"column:title" json:"title"`
	Content    string    `gorm:"column:content;type:text" json:"content"`
	Visibility bool      `gorm:"column:visibility" json:"visibility"`
	EntryDate  time.Time `gorm:"column:entry_date" json:"entryDate"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (Diary) TableName() string {
	return "diary"
}

// SituationRoom 情景房间模型
type SituationRoom struct {
	Code       string    `gorm:"column:code;primaryKey" json:"code"`
	OwnerId    string    `gorm:"column:owner_id;index" json:"ownerId"`
	MaxMembers int       `gorm:"column:max_members" json:"maxMembers"`
	Status     string    `gorm:"column:status;default:'waiting'" json:"status"` // waiting, running, finished
	ScenarioId string    `gorm:"column:scenario_id" json:"scenarioId"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (SituationRoom) TableName() string {
	return "situation_room"
}

// RoomMember 房间成员模型
type RoomMember struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code       string    `gorm:"column:code;index" json:"code"`
	UserId     string    `gorm:"column:user_id;index" json:"userId"`
	JoinTime   time.Time `gorm:"column:join_time;autoCreateTime" json:"joinTime"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
}

func (RoomMember) TableName() string {
	return "room_member"
}

// RoomNarrative 房间叙述模型
type RoomNarrative struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code       string    `gorm:"column:code;index" json:"code"`
	UserId     string    `gorm:"column:user_id;index" json:"userId"`
	Narrative  string    `gorm:"column:narrative;type:text" json:"narrative"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
}

func (RoomNarrative) TableName() string {
	return "room_narrative"
}
