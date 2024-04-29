package gormx

import (
	"github.com/sony/sonyflake"
	"github.com/zeddy-go/zeddy/container"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

var _ callbacks.BeforeCreateInterface = (*SnowflakeID)(nil)

type UnixTimestampMilli struct {
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type CommonField struct {
	SnowflakeID
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`
}

type SnowflakeID struct {
	ID uint64 `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
}

func (s *SnowflakeID) BeforeCreate(tx *gorm.DB) (err error) {
	snowflake, err := container.Resolve[*sonyflake.Sonyflake]()
	if err != nil {
		return
	}
	if s.ID == 0 {
		s.ID, err = snowflake.NextID()
	}
	return
}
