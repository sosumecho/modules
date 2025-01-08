package utils

import (
	"github.com/sony/sonyflake"
	"strconv"
	"time"
)

var snowFlake *sonyflake.Sonyflake

// UniqueID 全局唯一ID
func UniqueID() uint64 {
	if snowFlake == nil {
		snowFlake = sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			// StartTime: time.Now(),
		})
	}
	id, _ := snowFlake.NextID()
	return id
}

// UniqueIDStr id
func UniqueIDStr() string {
	id := UniqueID()
	return strconv.FormatUint(id, 10)
}
