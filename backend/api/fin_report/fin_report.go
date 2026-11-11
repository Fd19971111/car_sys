package fin_report

import (
	"context"

	v1 "engineServ/api/fin_report/v1"
)

// IFinReportV1 财务报表接口（登录后可查看）。
type IFinReportV1 interface {
	Overview(ctx context.Context, req *v1.ReportOverviewReq) (res *v1.ReportOverviewRes, err error)
	Flows(ctx context.Context, req *v1.ReportFlowListReq) (res *v1.ReportFlowListRes, err error)
	ProjectProfitList(ctx context.Context, req *v1.ProjectProfitListReq) (res *v1.ProjectProfitListRes, err error)
	ProjectProfitDetail(ctx context.Context, req *v1.ProjectProfitDetailReq) (res *v1.ProjectProfitDetailRes, err error)
}
