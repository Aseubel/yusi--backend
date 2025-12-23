// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package room

import (
	"context"

	"yusi-backend/internal/svc"
	"yusi-backend/internal/types"
	"yusi-backend/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 加入房间
func NewJoinRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinRoomLogic {
	return &JoinRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinRoomLogic) JoinRoom(req *types.JoinRoomRequest) (resp *types.Response, err error) {
	// 验证参数
	if req.Code == "" {
		return &types.Response{
			Code:    400,
			Message: "房间代码不能为空",
		}, nil
	}

	// 查询房间是否存在
	var room model.SituationRoom
	if err := l.svcCtx.DB.Where("code = ?", req.Code).First(&room).Error; err != nil {
		return &types.Response{
			Code:    404,
			Message: "房间不存在",
		}, nil
	}

	// 检查房间状态
	if room.Status != "waiting" {
		return &types.Response{
			Code:    400,
			Message: "房间已开始或已结束，无法加入",
		}, nil
	}

	// 检查用户是否已在房间中
	var existingMember model.RoomMember
	result := l.svcCtx.DB.Where("code = ? AND user_id = ?", req.Code, req.UserId).First(&existingMember)
	if result.Error == nil {
		return &types.Response{
			Code:    400,
			Message: "您已在房间中",
		}, nil
	}

	// 检查房间人数是否已满
	var memberCount int64
	l.svcCtx.DB.Model(&model.RoomMember{}).Where("code = ?", req.Code).Count(&memberCount)
	if int(memberCount) >= room.MaxMembers {
		return &types.Response{
			Code:    400,
			Message: "房间已满",
		}, nil
	}

	// 加入房间
	member := model.RoomMember{
		Code:   req.Code,
		UserId: req.UserId,
	}

	if err := l.svcCtx.DB.Create(&member).Error; err != nil {
		return &types.Response{
			Code:    500,
			Message: "加入房间失败",
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "加入成功",
		Data: map[string]interface{}{
			"code":       room.Code,
			"ownerId":    room.OwnerId,
			"maxMembers": room.MaxMembers,
			"status":     room.Status,
		},
	}, nil
}
