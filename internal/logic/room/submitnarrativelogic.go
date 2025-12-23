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

type SubmitNarrativeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交叙述
func NewSubmitNarrativeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitNarrativeLogic {
	return &SubmitNarrativeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitNarrativeLogic) SubmitNarrative(req *types.SubmitNarrativeRequest) (resp *types.Response, err error) {
	// 验证参数
	if req.Code == "" || req.Narrative == "" {
		return &types.Response{
			Code:    400,
			Message: "房间代码和叙述内容不能为空",
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

	// 检查房间状态
	if room.Status != "running" {
		return &types.Response{
			Code:    400,
			Message: "房间未开始或已结束",
		}, nil
	}

	// 检查用户是否在房间中
	var member model.RoomMember
	if err := l.svcCtx.DB.Where("code = ? AND user_id = ?", req.Code, req.UserId).First(&member).Error; err != nil {
		return &types.Response{
			Code:    403,
			Message: "您不在此房间中",
		}, nil
	}

	// 检查用户是否已提交过叙述
	var existingNarrative model.RoomNarrative
	result := l.svcCtx.DB.Where("code = ? AND user_id = ?", req.Code, req.UserId).First(&existingNarrative)
	if result.Error == nil {
		return &types.Response{
			Code:    400,
			Message: "您已提交过叙述",
		}, nil
	}

	// 保存叙述
	narrative := model.RoomNarrative{
		Code:      req.Code,
		UserId:    req.UserId,
		Narrative: req.Narrative,
	}

	if err := l.svcCtx.DB.Create(&narrative).Error; err != nil {
		return &types.Response{
			Code:    500,
			Message: "提交叙述失败",
		}, nil
	}

	// 检查是否所有成员都已提交
	var memberCount, narrativeCount int64
	l.svcCtx.DB.Model(&model.RoomMember{}).Where("code = ?", req.Code).Count(&memberCount)
	l.svcCtx.DB.Model(&model.RoomNarrative{}).Where("code = ?", req.Code).Count(&narrativeCount)

	allSubmitted := memberCount == narrativeCount

	return &types.Response{
		Code:    200,
		Message: "提交成功",
		Data: map[string]interface{}{
			"code":         req.Code,
			"allSubmitted": allSubmitted,
		},
	}, nil
}
