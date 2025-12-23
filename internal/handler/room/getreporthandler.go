// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package room

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yusi-backend/internal/logic/room"
	"yusi-backend/internal/svc"
)

// 获取报告
func GetReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL路径获取code
		code := r.URL.Query().Get(":code")
		if code == "" {
			// 尝试从URL路径获取
			code = r.URL.Path[len("/api/room/report/"):]
		}

		l := room.NewGetReportLogic(r.Context(), svcCtx, code)
		resp, err := l.GetReport()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
