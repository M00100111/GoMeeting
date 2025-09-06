package models

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MeetingInfoModel = (*customMeetingInfoModel)(nil)

type (
	// MeetingInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMeetingInfoModel.
	MeetingInfoModel interface {
		meetingInfoModel
	}

	customMeetingInfoModel struct {
		*defaultMeetingInfoModel
	}
)

// NewMeetingInfoModel returns a model for the database table.
func NewMeetingInfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MeetingInfoModel {
	return &customMeetingInfoModel{
		defaultMeetingInfoModel: newMeetingInfoModel(conn, c, opts...),
	}
}
