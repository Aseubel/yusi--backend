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

type GetReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	code   string
}

// 获取报告
func NewGetReportLogic(ctx context.Context, svcCtx *svc.ServiceContext, code string) *GetReportLogic {
	return &GetReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		code:   code,
	}
}

func (l *GetReportLogic) GetReport() (resp *types.Response, err error) {
	// 验证参数
	if l.code == "" {
		return &types.Response{
			Code:    400,
			Message: "房间代码不能为空",
		}, nil
	}

	// 查询房间
	var room model.SituationRoom
	if err := l.svcCtx.DB.Where("code = ?", l.code).First(&room).Error; err != nil {
		return &types.Response{
			Code:    404,
			Message: "房间不存在",
		}, nil
	}

	// 检查房间状态（可以查看正在运行或已结束的房间报告）
	if room.Status == "waiting" {
		return &types.Response{
			Code:    400,
			Message: "房间尚未开始",
		}, nil
	}

	// 获取所有成员
	var members []model.RoomMember
	l.svcCtx.DB.Where("code = ?", l.code).Find(&members)

	// 获取所有叙述
	var narratives []model.RoomNarrative
	l.svcCtx.DB.Where("code = ?", l.code).Find(&narratives)

	// 组织叙述数据
	narrativeMap := make(map[string]string)
	for _, n := range narratives {
		narrativeMap[n.UserId] = n.Narrative
	}

	// 生成报告摘要（简单版本，实际可能需要AI生成）
	summary := "本次情景房间活动已完成"
	if room.Status == "running" {
		summary = "本次情景房间活动正在进行中"
	}

	return &types.Response{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			"code":       room.Code,
			"ownerId":    room.OwnerId,
			"scenarioId": room.ScenarioId,
			"status":     room.Status,
			"summary":    summary,
			"narratives": narrativeMap,
			"totalMembers": len(members),
			"submittedCount": len(narratives),
		},
	}, nil
}
