package logic

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"GoMeeting/pkg/structs"
	"GoMeeting/pkg/structs/message"
	"GoMeeting/rpcs/meeting/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"strconv"
	"time"

	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinMeetingLogic {
	return &JoinMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinMeetingLogic) JoinMeeting(in *meeting.JoinMeetingReq) (*meeting.Result, error) {
	//获取会议信息
	meetingInfo, err := l.svcCtx.MeetingInfoModel.FindOneByMeetingId(l.ctx, in.MeetingId)
	//数据库操作出错
	if err != nil && err != sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.MeetingNotExistCode,
		}, nil
	}
	//判断会议状态是否已结束
	if meetingInfo.Status == 0 {
		//未到达会议时间
		if meetingInfo.StartTime.Valid && meetingInfo.StartTime.Time.After(time.Now()) {
			return &meeting.Result{
				Code: code.MeetingNotStartedCode,
			}, nil
		}
		//已过会议时间且已更新会议结束时间
		if meetingInfo.EndTime.Valid && meetingInfo.EndTime.Time.Before(time.Now()) {
			return &meeting.Result{
				Code: code.MeetingAlreadyEndedCode,
			}, nil
		}
		return &meeting.Result{
			Code: code.MeetingNotStartedCode,
		}, nil
	}
	//判断加入方式
	//验证密码
	if meetingInfo.JoinType == 1 && meetingInfo.MeetingPassword != in.Password {
		return &meeting.Result{
			Code: code.MeetingJoinCodeErrorCode,
		}, nil
	}

	//查询会议成员记录
	meetingMember, err := l.svcCtx.MeetingMemberModel.FindOneByMeetingIdUserId(l.ctx, meetingInfo.Id, in.UserIndex)
	//数据库操作出错
	if err != nil && err != sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == sqlc.ErrNotFound {
		//第一次加入新增信息
		_, err := l.svcCtx.MeetingMemberModel.Insert(l.ctx, &models.MeetingMember{
			MeetingId: meetingInfo.Id,
			UserId:    in.UserIndex,
			LastJoinTime: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			return &meeting.Result{
				Code: code.ErrDbOpCode,
				Msg:  "新增会议成员信息失败",
			}, nil
		}
		meetingMember, err = l.svcCtx.MeetingMemberModel.FindOneByMeetingIdUserId(l.ctx, meetingInfo.Id, in.UserIndex)
		//数据库操作出错
		if err != nil && err != sqlc.ErrNotFound {
			return &meeting.Result{
				Code: code.ErrDbOpCode,
			}, nil
		}
	} else {
		//判断是否被拉黑
		if meetingMember.UserStatus != 0 {
			return &meeting.Result{
				Code: code.MeetingUserNotAllowedCode,
			}, nil
		}
		//非第一次加入修改入会时间
		meetingMember.LastJoinTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err = l.svcCtx.MeetingMemberModel.Update(l.ctx, meetingMember)
		if err != nil {
			return &meeting.Result{
				Code: code.ErrDbOpCode,
				Msg:  "修改会议成员信息失败",
			}, nil
		}
	}

	// 将会议室成员id和状态信息分别添加到Redis

	// 记录系统中正在开会的用户，存储成员Id到 Redis Hash (成员ID -> 会议ID)
	err = l.svcCtx.Redis.HsetCtx(l.ctx, ctxdata.OnMeetingUserPrefix, strconv.FormatUint(in.UserId, 10), strconv.FormatUint(in.MeetingId, 10))
	if err != nil {
		l.Logger.Errorf("Failed to store member ID in RedisHash: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}

	// 记录在会议中的用户，存储成员Id到 Redis Set
	memberSetKey := fmt.Sprintf(ctxdata.MeetingMemberPrefix, in.MeetingId)
	val, err := l.svcCtx.Redis.SaddCtx(l.ctx, memberSetKey, in.UserId)
	if err != nil {
		l.Logger.Errorf("Failed to store member ID in RedisList: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}
	//val为受影响的元素
	if val == 0 {
		return &meeting.Result{
			Code: code.MeetingAlreadyInCode,
		}, nil
	}
	// 发布Redis事件
	if val > 0 {
		// 发布新成员加入事件到 Kafka 或其他消息队列
		eventData := message.MeetingMemberJoinNoticeData{
			MeetingId: in.MeetingId,
			UserId:    in.UserId,
		}
		eventDataJson, _ := json.Marshal(eventData)
		event := message.Message{
			MessageType: message.Notification_Message,
			Method:      message.Meeting_Member_Join_Notice_Method,
			Data:        eventDataJson,
		}
		eventJson, _ := json.Marshal(event)
		// 使用 KafkaWsPusher 发布事件
		err := l.svcCtx.KafkaWsPusher.Push(l.ctx, string(eventJson))
		if err != nil {
			l.Logger.Errorf("Failed to publish event to Kafka: %v", err)
			return &meeting.Result{
				Code: code.ErrKafkaPushCode,
			}, nil
		}
	}
	//存储成员信息到 Redis Hash
	memberStatusKey := fmt.Sprintf(ctxdata.MeetingMemberDetailPrefix, in.MeetingId)
	memberStatus := structs.MemberStatus{
		UserId:       strconv.FormatUint(in.UserId, 10),
		Username:     in.Username,
		Sex:          in.Sex,
		Email:        in.Email,
		UserStatus:   meetingMember.UserStatus,
		UserType:     meetingMember.UserType,
		MicStatus:    in.MicStatus,
		CameraStatus: in.CameraStatus,
		ScreenStatus: in.ScreenStatus,
	}
	statusData, _ := json.Marshal(memberStatus)
	err = l.svcCtx.Redis.HsetCtx(l.ctx, memberStatusKey, strconv.FormatUint((in.UserId), 10), string(statusData))
	if err != nil {
		l.Logger.Errorf("Failed to store member status in RedisHash: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}

	//入会成功
	return &meeting.Result{
		Code: code.SUCCESSCode,
	}, nil
}
