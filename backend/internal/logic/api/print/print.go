// Package print 打印设置服务：每个公司一份全局配置，打印单据时作为公共的
// 抬头 / 页脚 / 日期 / 电子章样式来源。
package print

import (
	"context"

	"engineServ/internal/consts"
	v1 "engineServ/api/print/v1"
	uploadLogic "engineServ/internal/logic/upload"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// Service 打印设置服务
type Service struct{}

// New 创建打印设置服务
func New() *Service { return &Service{} }

// Get 获取当前公司的打印设置；从未保存过时返回默认值。
func (s *Service) Get(ctx context.Context) (*v1.PrintSettingModel, error) {
	row, err := g.DB().GetOne(ctx, `SELECT title,sub_title,footer,show_date,show_seal,seal_image
		FROM print_setting WHERE company_id=?`, currentCompanyId(ctx))
	if err != nil {
		glog.Errorf(ctx, "查询打印设置失败: %v", err)
		return nil, gerror.New("查询打印设置失败，请稍后重试")
	}

	// 未配置时给出一套开箱即用的默认值
	if row == nil || row.IsEmpty() {
		return &v1.PrintSettingModel{
			Title:    "",
			SubTitle: "",
			Footer:   "本单据由项目管理系统生成",
			ShowDate: true,
			ShowSeal: false,
		}, nil
	}

	return &v1.PrintSettingModel{
		Title:        row["title"].String(),
		SubTitle:     row["sub_title"].String(),
		Footer:       row["footer"].String(),
		ShowDate:     row["show_date"].Bool(),
		ShowSeal:     row["show_seal"].Bool(),
		SealImage:    row["seal_image"].String(),
		SealImageUrl: uploadLogic.AbsoluteURL(ctx, row["seal_image"].String()),
	}, nil
}

// Save 保存打印设置（按公司维度 upsert）。
func (s *Service) Save(ctx context.Context, in *v1.SettingSaveReq) error {
	companyId := currentCompanyId(ctx)

	count, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM print_setting WHERE company_id=?", companyId)
	if err != nil {
		glog.Errorf(ctx, "校验打印设置失败: %v", err)
		return gerror.New("保存打印设置失败，请稍后重试")
	}

	if count.Int() > 0 {
		if _, err = g.DB().Exec(ctx, `UPDATE print_setting SET
			title=?,sub_title=?,footer=?,show_date=?,show_seal=?,seal_image=?
			WHERE company_id=?`,
			in.Title, in.SubTitle, in.Footer, in.ShowDate, in.ShowSeal, in.SealImage, companyId); err != nil {
			glog.Errorf(ctx, "更新打印设置失败: %v", err)
			return gerror.New("保存打印设置失败，请稍后重试")
		}

		return nil
	}

	if _, err = g.DB().Exec(ctx, `INSERT INTO print_setting
		(company_id,title,sub_title,footer,show_date,show_seal,seal_image)
		VALUES(?,?,?,?,?,?,?)`,
		companyId, in.Title, in.SubTitle, in.Footer, in.ShowDate, in.ShowSeal, in.SealImage); err != nil {
		glog.Errorf(ctx, "写入打印设置失败: %v", err)
		return gerror.New("保存打印设置失败，请稍后重试")
	}

	return nil
}

func currentCompanyId(ctx context.Context) uint {
	return gconv.Uint(ctx.Value(consts.CtxCompanyIdKey))
}
