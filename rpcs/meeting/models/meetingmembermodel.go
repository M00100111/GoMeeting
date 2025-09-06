package models

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MeetingMemberModel = (*customMeetingMemberModel)(nil)

type (
	// MeetingMemberModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMeetingMemberModel.
	MeetingMemberModel interface {
		meetingMemberModel
	}

	customMeetingMemberModel struct {
		*defaultMeetingMemberModel
	}
)

// NewMeetingMemberModel returns a model for the database table.
func NewMeetingMemberModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MeetingMemberModel {
	return &customMeetingMemberModel{
		defaultMeetingMemberModel: newMeetingMemberModel(conn, c, opts...),
	}
}
