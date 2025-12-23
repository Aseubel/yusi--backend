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

type StartRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 开始房间
func NewStartRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartRoomLogic {
	return &StartRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartRoomLogic) StartRoom(req *types.StartRoomRequest) (resp *types.Response, err error) {
	// 验证参数
	if req.Code == "" {
		return &types.Response{
			Code:    400,
			Message: "房间代码不能为空",
		}, nil
	}

	// 查询房间
	var room model.SituationRoom
	if err := l.svcCtx.DB.Where("code = ?", req.Code).First(&room).Error; err != nil {
		return &types.Response{
			Code:    404,
			Message: "房间不存在",
		}, nil
	}

	// 验证权限：只有房主可以开始房间
	if room.OwnerId != req.OwnerId {
		return &types.Response{
			Code:    403,
			Message: "只有房主可以开始房间",
		}, nil
	}

	// 检查房间状态
	if room.Status != "waiting" {
		return &types.Response{
			Code:    400,
			Message: "房间已开始或已结束",
		}, nil
	}

	// 检查房间人数
	var memberCount int64
	l.svcCtx.DB.Model(&model.RoomMember{}).Where("code = ?", req.Code).Count(&memberCount)
	if memberCount < 2 {
		return &types.Response{
			Code:    400,
			Message: "房间人数不足，至少需要2人",
		}, nil
	}

	// 更新房间状态和场景ID
	updates := map[string]interface{}{
		"status":      "running",
		"scenario_id": req.ScenarioId,
	}

	if err := l.svcCtx.DB.Model(&room).Updates(updates).Error; err != nil {
		return &types.Response{
			Code:    500,
			Message: "开始房间失败",
		}, nil
	}

	return &types.Response{
		Code:    200,
		Message: "房间已开始",
		Data: map[string]interface{}{
			"code":       room.Code,
			"status":     "running",
			"scenarioId": req.ScenarioId,
		},
	}, nil
}
