package fin_report

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// 项目成本 / 毛利口径说明
//
//   - 收入侧：结算金额 = 合同金额 + 变更净额（增项 - 减项）；
//   - 成本侧：按「付款」归集，只统计关联到项目的付款（project_id>0），
//     公司公共支出（project_id=0）不摊入项目成本，单独展示；
//   - 材料耗用：按「领用数量 × 材料单价」计算，仅作对照参考，
//     不重复计入成本（材料款付款已经在成本里）；
//   - 预计毛利 = 结算金额 - 成本合计；在建项目成本尚未完全发生，毛利为预计值；
//   - 现金结余 = 已收款 - 成本合计，反映项目当前是垫资还是回正。

// 成本类型展示顺序（未知类型追加在后面）
var costTypeOrder = []string{"材料款", "人工费", "施工费", "设计费", "管理费", "销售费", "其他"}

// ProfitCostItem 成本构成项
type ProfitCostItem struct {
	Type    string
	Amount  float64
	Percent float64
}

// ProjectProfit 项目毛利行
type ProjectProfit struct {
	ProjectId        int64
	ProjectName      string
	CustomerName     string
	Status           int
	StartDate        string
	EndDate          string
	ContractAmount   float64
	ChangeNet        float64
	SettlementAmount float64
	ReceivedAmount   float64
	UnpaidAmount     float64
	ReceivedRate     float64
	CostTotal        float64
	CostItems        []*ProfitCostItem
	MaterialPaid     float64
	MaterialUsed     float64
	GrossProfit      float64
	GrossRate        float64
	CashProfit       float64
}

// ProjectProfitSummary 项目毛利汇总
type ProjectProfitSummary struct {
	ProjectCount      int
	SettlementTotal   float64
	ReceivedTotal     float64
	UnpaidTotal       float64
	CostTotal         float64
	GrossProfit       float64
	GrossRate         float64
	LossCount         int
	CompanyCost       float64
	MaterialUsedTotal float64
}

// ProjectProfitList 项目毛利列表
type ProjectProfitList struct {
	Records []*ProjectProfit
	Summary *ProjectProfitSummary
}

// ProfitReceipt 项目收款明细
type ProfitReceipt struct {
	Id          int64
	PaymentType string
	Amount      float64
	PaymentDate string
	Remark      string
}

// ProfitPayout 项目付款明细
type ProfitPayout struct {
	Id         int64
	PayoutType string
	Amount     float64
	PayoutDate string
	VendorName string
	Remark     string
}

// ProfitMaterial 项目材料领用明细
type ProfitMaterial struct {
	Id           int64
	MaterialName string
	Unit         string
	Quantity     float64
	UnitPrice    float64
	Amount       float64
	UsedAt       string
	Receiver     string
}

// ProfitChange 项目变更明细
type ProfitChange struct {
	Id         int64
	ChangeType string
	Title      string
	Amount     float64
	ChangeDate string
	Reason     string
}

// ProjectProfitDetail 项目毛利明细
type ProjectProfitDetail struct {
	ProjectProfit
	Receipts  []*ProfitReceipt
	Payouts   []*ProfitPayout
	Materials []*ProfitMaterial
	Changes   []*ProfitChange
}

// ProjectProfits 项目毛利列表。
//
// 数据量小（一个公司几十个项目），一次性返回全部行，前端做排序 / 只看亏损。
func (s *Service) ProjectProfits(ctx context.Context, keyword string, status *int, onlyLoss bool) (*ProjectProfitList, error) {
	rows, err := s.loadProjectProfitBase(ctx, nil, keyword, status)
	if err != nil {
		return nil, err
	}

	projectIds := make([]int64, 0, len(rows))
	for _, row := range rows {
		projectIds = append(projectIds, row.ProjectId)
	}

	costs, err := s.loadProjectCosts(ctx, projectIds)
	if err != nil {
		return nil, err
	}

	records := make([]*ProjectProfit, 0, len(rows))
	for _, row := range rows {
		item := buildProjectProfit(row, costs[row.ProjectId])
		if onlyLoss && item.GrossProfit >= 0 {
			continue
		}
		records = append(records, item)
	}

	// 默认把亏钱的排在前面，方便老板先看问题项目
	sort.SliceStable(records, func(i, j int) bool {
		if records[i].GrossProfit == records[j].GrossProfit {
			return records[i].ProjectId > records[j].ProjectId
		}
		return records[i].GrossProfit < records[j].GrossProfit
	})

	companyCost, err := g.DB().GetValue(ctx,
		"SELECT COALESCE(SUM(amount),0) FROM project_payout WHERE deleted_at IS NULL AND project_id=0")
	if err != nil {
		glog.Errorf(ctx, "统计公司公共支出失败: %v", err)
	}

	return &ProjectProfitList{
		Records: records,
		Summary: buildProjectProfitSummary(records, companyCost.Float64()),
	}, nil
}

// ProjectProfitDetail 单个项目的成本归集与明细。
func (s *Service) ProjectProfitDetail(ctx context.Context, projectId int64) (*ProjectProfitDetail, error) {
	rows, err := s.loadProjectProfitBase(ctx, &projectId, "", nil)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, gerror.New("项目不存在或已被删除")
	}

	costs, err := s.loadProjectCosts(ctx, []int64{projectId})
	if err != nil {
		return nil, err
	}

	detail := &ProjectProfitDetail{
		ProjectProfit: *buildProjectProfit(rows[0], costs[projectId]),
		Receipts:      make([]*ProfitReceipt, 0),
		Payouts:       make([]*ProfitPayout, 0),
		Materials:     make([]*ProfitMaterial, 0),
		Changes:       make([]*ProfitChange, 0),
	}

	if err = s.fillProfitReceipts(ctx, projectId, detail); err != nil {
		return nil, err
	}
	if err = s.fillProfitPayouts(ctx, projectId, detail); err != nil {
		return nil, err
	}
	if err = s.fillProfitMaterials(ctx, projectId, detail); err != nil {
		return nil, err
	}
	if err = s.fillProfitChanges(ctx, projectId, detail); err != nil {
		return nil, err
	}

	return detail, nil
}

// profitBaseRow 项目 + 收款 / 变更 / 材料耗用 的基础行。
type profitBaseRow struct {
	ProjectId      int64
	ProjectName    string
	CustomerName   string
	Status         int
	StartDate      string
	EndDate        string
	ContractAmount float64
	ChangeNet      float64
	ReceivedAmount float64
	MaterialUsed   float64
}

// loadProjectProfitBase 查询项目基础数据（项目、客户、变更净额、已收款、材料耗用）。
func (s *Service) loadProjectProfitBase(ctx context.Context, projectId *int64, keyword string, status *int) ([]*profitBaseRow, error) {
	where := []string{"p.deleted_at IS NULL"}
	args := make([]interface{}, 0, 4)

	if projectId != nil {
		where = append(where, "p.id=?")
		args = append(args, *projectId)
	}
	if status != nil && *status > 0 {
		where = append(where, "p.status=?")
		args = append(args, *status)
	}
	if text := strings.TrimSpace(keyword); text != "" {
		like := "%" + text + "%"
		where = append(where, "(p.name LIKE ? OR c.name LIKE ?)")
		args = append(args, like, like)
	}

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT
			p.id,p.name,p.status,p.amount,
			COALESCE(c.name,'') AS customer_name,
			DATE_FORMAT(p.start_date,'%%Y-%%m-%%d') AS start_date,
			DATE_FORMAT(p.end_date,'%%Y-%%m-%%d') AS end_date,
			COALESCE(ch.change_net,0) AS change_net,
			COALESCE(pay.received,0) AS received,
			COALESCE(mu.material_used,0) AS material_used
		FROM project p
		LEFT JOIN customer c ON c.id=p.customer_id
		LEFT JOIN (
			SELECT project_id,SUM(CASE WHEN change_type='ADD' THEN amount ELSE -amount END) AS change_net
			FROM project_change WHERE deleted_at IS NULL GROUP BY project_id
		) ch ON ch.project_id=p.id
		LEFT JOIN (
			SELECT project_id,SUM(amount) AS received
			FROM project_payment WHERE deleted_at IS NULL GROUP BY project_id
		) pay ON pay.project_id=p.id
		LEFT JOIN (
			SELECT u.project_id,SUM(u.quantity*m.unit_price) AS material_used
			FROM project_material_usage u
			INNER JOIN project_material m ON m.id=u.material_id
			WHERE u.deleted_at IS NULL
			GROUP BY u.project_id
		) mu ON mu.project_id=p.id
		WHERE %s
		ORDER BY p.id DESC
		LIMIT 500`, strings.Join(where, " AND ")), args...)
	if err != nil {
		glog.Errorf(ctx, "查询项目毛利基础数据失败: %v", err)
		return nil, gerror.New("查询项目毛利失败，请稍后重试")
	}

	rows := make([]*profitBaseRow, 0, len(records))
	for _, record := range records {
		rows = append(rows, &profitBaseRow{
			ProjectId:      record["id"].Int64(),
			ProjectName:    record["name"].String(),
			CustomerName:   record["customer_name"].String(),
			Status:         record["status"].Int(),
			StartDate:      record["start_date"].String(),
			EndDate:        record["end_date"].String(),
			ContractAmount: round2(record["amount"].Float64()),
			ChangeNet:      round2(record["change_net"].Float64()),
			ReceivedAmount: round2(record["received"].Float64()),
			MaterialUsed:   round2(record["material_used"].Float64()),
		})
	}

	return rows, nil
}

// loadProjectCosts 按项目归集付款成本：projectId → 成本类型 → 金额。
func (s *Service) loadProjectCosts(ctx context.Context, projectIds []int64) (map[int64]map[string]float64, error) {
	result := make(map[int64]map[string]float64, len(projectIds))
	if len(projectIds) == 0 {
		return result, nil
	}

	args := make([]interface{}, 0, len(projectIds))
	for _, id := range projectIds {
		args = append(args, id)
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(projectIds)), ",")

	records, err := g.DB().GetAll(ctx, fmt.Sprintf(`SELECT project_id,payout_type,SUM(amount) AS amount
		FROM project_payout
		WHERE deleted_at IS NULL AND project_id>0 AND project_id IN (%s)
		GROUP BY project_id,payout_type`, placeholders), args...)
	if err != nil {
		glog.Errorf(ctx, "归集项目付款成本失败: %v", err)
		return nil, gerror.New("归集项目成本失败，请稍后重试")
	}

	for _, record := range records {
		projectId := record["project_id"].Int64()
		if result[projectId] == nil {
			result[projectId] = make(map[string]float64)
		}

		payoutType := strings.TrimSpace(record["payout_type"].String())
		if payoutType == "" {
			payoutType = "其他"
		}
		result[projectId][payoutType] += round2(record["amount"].Float64())
	}

	return result, nil
}

// buildProjectProfit 计算单行毛利指标。
func buildProjectProfit(row *profitBaseRow, costMap map[string]float64) *ProjectProfit {
	costTotal := 0.0
	for _, amount := range costMap {
		costTotal += amount
	}
	costTotal = round2(costTotal)

	settlement := round2(row.ContractAmount + row.ChangeNet)
	unpaid := round2(settlement - row.ReceivedAmount)
	if unpaid < 0 {
		unpaid = 0
	}

	profit := round2(settlement - costTotal)
	cashProfit := round2(row.ReceivedAmount - costTotal)

	item := &ProjectProfit{
		ProjectId:        row.ProjectId,
		ProjectName:      row.ProjectName,
		CustomerName:     row.CustomerName,
		Status:           row.Status,
		StartDate:        row.StartDate,
		EndDate:          row.EndDate,
		ContractAmount:   row.ContractAmount,
		ChangeNet:        row.ChangeNet,
		SettlementAmount: settlement,
		ReceivedAmount:   row.ReceivedAmount,
		UnpaidAmount:     unpaid,
		CostTotal:        costTotal,
		MaterialPaid:     round2(costMap["材料款"]),
		MaterialUsed:     row.MaterialUsed,
		GrossProfit:      profit,
		CashProfit:       cashProfit,
		CostItems:        buildCostItems(costMap, costTotal),
	}

	if settlement > 0 {
		item.ReceivedRate = round1(row.ReceivedAmount / settlement * 100)
		item.GrossRate = round1(profit / settlement * 100)
	}

	return item
}

// buildCostItems 生成成本构成（按固定顺序 + 未知类型追加）。
func buildCostItems(costMap map[string]float64, total float64) []*ProfitCostItem {
	if len(costMap) == 0 {
		return make([]*ProfitCostItem, 0)
	}

	items := make([]*ProfitCostItem, 0, len(costMap))
	used := make(map[string]struct{}, len(costMap))

	for _, costType := range costTypeOrder {
		amount := round2(costMap[costType])
		if amount == 0 {
			continue
		}
		used[costType] = struct{}{}
		items = append(items, &ProfitCostItem{Type: costType, Amount: amount})
	}

	extra := make([]*ProfitCostItem, 0)
	for costType, amount := range costMap {
		if _, ok := used[costType]; ok {
			continue
		}
		if value := round2(amount); value != 0 {
			extra = append(extra, &ProfitCostItem{Type: costType, Amount: value})
		}
	}
	sort.Slice(extra, func(i, j int) bool { return extra[i].Amount > extra[j].Amount })
	items = append(items, extra...)

	if total > 0 {
		for _, item := range items {
			item.Percent = round1(item.Amount / total * 100)
		}
	}

	return items
}

// buildProjectProfitSummary 汇总（只统计过滤后的行）。
func buildProjectProfitSummary(records []*ProjectProfit, companyCost float64) *ProjectProfitSummary {
	summary := &ProjectProfitSummary{
		ProjectCount: len(records),
		CompanyCost:  round2(companyCost),
	}

	for _, item := range records {
		summary.SettlementTotal += item.SettlementAmount
		summary.ReceivedTotal += item.ReceivedAmount
		summary.UnpaidTotal += item.UnpaidAmount
		summary.CostTotal += item.CostTotal
		summary.GrossProfit += item.GrossProfit
		summary.MaterialUsedTotal += item.MaterialUsed

		if item.GrossProfit < 0 && (item.SettlementAmount > 0 || item.CostTotal > 0) {
			summary.LossCount++
		}
	}

	summary.SettlementTotal = round2(summary.SettlementTotal)
	summary.ReceivedTotal = round2(summary.ReceivedTotal)
	summary.UnpaidTotal = round2(summary.UnpaidTotal)
	summary.CostTotal = round2(summary.CostTotal)
	summary.GrossProfit = round2(summary.GrossProfit)
	summary.MaterialUsedTotal = round2(summary.MaterialUsedTotal)

	if summary.SettlementTotal > 0 {
		summary.GrossRate = round1(summary.GrossProfit / summary.SettlementTotal * 100)
	}

	return summary
}

// fillProfitReceipts 项目收款明细。
func (s *Service) fillProfitReceipts(ctx context.Context, projectId int64, detail *ProjectProfitDetail) error {
	records, err := g.DB().GetAll(ctx, `SELECT id,payment_type,amount,remark,
			DATE_FORMAT(payment_date,'%Y-%m-%d') AS payment_date
		FROM project_payment
		WHERE deleted_at IS NULL AND project_id=?
		ORDER BY payment_date DESC,id DESC
		LIMIT 200`, projectId)
	if err != nil {
		glog.Errorf(ctx, "查询项目收款明细失败, projectId=%d: %v", projectId, err)
		return gerror.New("查询项目收款明细失败，请稍后重试")
	}

	for _, record := range records {
		detail.Receipts = append(detail.Receipts, &ProfitReceipt{
			Id:          record["id"].Int64(),
			PaymentType: record["payment_type"].String(),
			Amount:      round2(record["amount"].Float64()),
			PaymentDate: record["payment_date"].String(),
			Remark:      record["remark"].String(),
		})
	}

	return nil
}

// fillProfitPayouts 项目付款明细（成本）。
func (s *Service) fillProfitPayouts(ctx context.Context, projectId int64, detail *ProjectProfitDetail) error {
	records, err := g.DB().GetAll(ctx, `SELECT p.id,p.payout_type,p.amount,COALESCE(v.name,'') AS vendor_name,p.remark,
			DATE_FORMAT(payout_date,'%Y-%m-%d') AS payout_date
		FROM project_payout p
		LEFT JOIN vendor v ON v.id=p.vendor_id
		WHERE p.deleted_at IS NULL AND p.project_id=?
		ORDER BY p.payout_date DESC,p.id DESC
		LIMIT 200`, projectId)
	if err != nil {
		glog.Errorf(ctx, "查询项目付款明细失败, projectId=%d: %v", projectId, err)
		return gerror.New("查询项目付款明细失败，请稍后重试")
	}

	for _, record := range records {
		detail.Payouts = append(detail.Payouts, &ProfitPayout{
			Id:         record["id"].Int64(),
			PayoutType: record["payout_type"].String(),
			Amount:     round2(record["amount"].Float64()),
			PayoutDate: record["payout_date"].String(),
			VendorName: record["vendor_name"].String(),
			Remark:     record["remark"].String(),
		})
	}

	return nil
}

// fillProfitMaterials 项目材料领用明细（对照参考，不计入成本）。
func (s *Service) fillProfitMaterials(ctx context.Context, projectId int64, detail *ProjectProfitDetail) error {
	records, err := g.DB().GetAll(ctx, `SELECT u.id,u.quantity,u.receiver,u.remark,
			DATE_FORMAT(u.used_at,'%Y-%m-%d') AS used_at,
			COALESCE(m.name,'') AS material_name,
			COALESCE(m.unit,'') AS unit,
			COALESCE(m.unit_price,0) AS unit_price,
			ROUND(u.quantity*COALESCE(m.unit_price,0),2) AS amount
		FROM project_material_usage u
		LEFT JOIN project_material m ON m.id=u.material_id
		WHERE u.deleted_at IS NULL AND u.project_id=?
		ORDER BY u.used_at DESC,u.id DESC
		LIMIT 200`, projectId)
	if err != nil {
		glog.Errorf(ctx, "查询项目材料领用失败, projectId=%d: %v", projectId, err)
		return gerror.New("查询项目材料领用失败，请稍后重试")
	}

	for _, record := range records {
		detail.Materials = append(detail.Materials, &ProfitMaterial{
			Id:           record["id"].Int64(),
			MaterialName: record["material_name"].String(),
			Unit:         record["unit"].String(),
			Quantity:     round2(record["quantity"].Float64()),
			UnitPrice:    round2(record["unit_price"].Float64()),
			Amount:       round2(record["amount"].Float64()),
			UsedAt:       record["used_at"].String(),
			Receiver:     record["receiver"].String(),
		})
	}

	return nil
}

// fillProfitChanges 项目变更明细。
func (s *Service) fillProfitChanges(ctx context.Context, projectId int64, detail *ProjectProfitDetail) error {
	records, err := g.DB().GetAll(ctx, `SELECT id,change_type,title,amount,reason,
			DATE_FORMAT(change_date,'%Y-%m-%d') AS change_date
		FROM project_change
		WHERE deleted_at IS NULL AND project_id=?
		ORDER BY change_date DESC,id DESC
		LIMIT 200`, projectId)
	if err != nil {
		glog.Errorf(ctx, "查询项目变更明细失败, projectId=%d: %v", projectId, err)
		return gerror.New("查询项目变更明细失败，请稍后重试")
	}

	for _, record := range records {
		detail.Changes = append(detail.Changes, &ProfitChange{
			Id:         record["id"].Int64(),
			ChangeType: record["change_type"].String(),
			Title:      record["title"].String(),
			Amount:     round2(record["amount"].Float64()),
			ChangeDate: record["change_date"].String(),
			Reason:     record["reason"].String(),
		})
	}

	return nil
}
