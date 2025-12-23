// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package room

import (
	"context"

	"yusi-backend/internal/svc"
	"yusi-backend/internal/types"
	"yusi-backend/internal/utils"
	"yusi-backend/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建房间
func NewCreateRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoomLogic {
	return &CreateRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoomLogic) CreateRoom(req *types.CreateRoomRequest) (resp *types.Response, err error) {
	// 验证参数
	if req.MaxMembers < 2 || req.MaxMembers > 10 {
		return &types.Response{
			Code:    400,
			Message: "房间人数必须在2-10之间",
		}, nil
	}

	// 生成唯一房间代码
	var code string
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		code = utils.GenerateRoomCode()

		// 检查代码是否已存在
		var existingRoom model.SituationRoom
		result := l.svcCtx.DB.Where("code = ?", code).First(&existingRoom)
		if result.Error != nil {
			// 代码不存在，可以使用
			break
		}

		if i == maxRetries-1 {
			return &types.Response{
				Code:    500,
				Message: "生成房间代码失败，请重试",
			}, nil
		}
	}

	// 创建房间
	room := model.SituationRoom{
		Code:       code,
		OwnerId:    req.OwnerId,
		MaxMembers: req.MaxMembers,
		Status:     "waiting",
	}

	if err := l.svcCtx.DB.Create(&room).Error; err != nil {
		return &types.Response{
			Code:    500,
			Message: "创建房间失败",
		}, nil
	}

	// 房主自动加入房间
	member := model.RoomMember{
		Code:   code,
		UserId: req.OwnerId,
	}

	if err := l.svcCtx.DB.Create(&member).Error; err != nil {
		// 如果添加成员失败，删除已创建的房间
		l.svcCtx.DB.Delete(&room)
		return &types.Response{
			Code:    500,
			Message: "加入房间失败",
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "创建成功",
		Data: map[string]interface{}{
			"code":       code,
			"ownerId":    room.OwnerId,
			"maxMembers": room.MaxMembers,
			"status":     room.Status,
		},
	}, nil
}
