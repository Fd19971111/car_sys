package v1

import "github.com/gogf/gf/v2/frame/g"

type DriverBalance struct {
	DriverId     int64   `json:"driverId"`
	DriverName   string  `json:"driverName"`
	Phone        string  `json:"phone"`
	TripCount    int     `json:"tripCount"`
	EarnedTotal  float64 `json:"earnedTotal"`
	SettledTotal float64 `json:"settledTotal"`
	Payable      float64 `json:"payable"`
}

type DriversReq struct {
	g.Meta `path:"/driver-finance/drivers" method:"get" tags:"司机结算" summary:"司机待结余额"`
}

type DriversRes struct {
	Records []*DriverBalance `json:"records"`
}

type SettlementItem struct {
	Id              int64   `json:"id"`
	SettlementNo    string  `json:"settlementNo"`
	DriverId        int64   `json:"driverId"`
	DriverName      string  `json:"driverName"`
	SettlementDate  string  `json:"settlementDate"`
	Amount          float64 `json:"amount"`
	FundAccountName string  `json:"fundAccountName"`
	VoucherNo       string  `json:"voucherNo"`
	Remark          string  `json:"remark"`
}

type SettlementsReq struct {
	g.Meta   `path:"/driver-finance/settlements" method:"get" tags:"司机结算" summary:"司机结算记录"`
	DriverId int64 `json:"driverId"`
	Current  int   `json:"current" d:"1"`
	Size     int   `json:"size" d:"20"`
}

type SettlementsRes struct {
	Records []*SettlementItem `json:"records"`
	Total   int               `json:"total"`
}

type CreateSettlementReq struct {
	g.Meta         `path:"/driver-finance/settlement/create" method:"post" tags:"司机结算" summary:"结算司机车费"`
	DriverId       int64   `json:"driverId" v:"required|min:1#请选择司机|司机无效"`
	Amount         float64 `json:"amount" v:"required|min:0.01#请输入结算金额|结算金额必须大于0"`
	SettlementDate string  `json:"settlementDate" v:"required#请选择结算日期"`
	FundAccountId  int64   `json:"fundAccountId" v:"required|min:1#请选择付款账户|付款账户无效"`
	Remark         string  `json:"remark" v:"max-length:255#备注最多255个字符"`
}

type CreateSettlementRes struct {
	Id           int64  `json:"id"`
	SettlementNo string `json:"settlementNo"`
	VoucherNo    string `json:"voucherNo"`
}
