package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 库存管理 =====================

// StockListReq 库存分页列表请求（跨项目，项目ID为 0 时表示全部项目）
type StockListReq struct {
	g.Meta `path:"/inventory/stock/list" method:"get" tags:"库存管理" summary:"库存列表"`

	ProjectId   int64  `json:"projectId" dc:"项目ID，0 表示全部项目"`
	WarehouseId int64  `json:"warehouseId" dc:"仓库ID，0 表示全部仓库"`
	Keyword     string `json:"keyword"   dc:"材料名称 / 规格 / 供应商关键字"`
	Current     int    `json:"current"   d:"1"  dc:"当前页码"`
	Size        int    `json:"size"      d:"10" dc:"每页条数"`
}

// StockItem 库存行
type StockItem struct {
	Id             int64   `json:"id"             dc:"库存ID（材料ID）"`
	ProjectId      int64   `json:"projectId"      dc:"项目ID"`
	ProjectName    string  `json:"projectName"    dc:"项目名称"`
	WarehouseId    int64   `json:"warehouseId"    dc:"仓库ID"`
	WarehouseName  string  `json:"warehouseName"  dc:"仓库名称"`
	Name           string  `json:"name"           dc:"材料名称"`
	Spec           string  `json:"spec"           dc:"规格 / 型号"`
	Unit           string  `json:"unit"           dc:"单位"`
	Quantity       float64 `json:"quantity"       dc:"库存数量（累计入库）"`
	UnitPrice      float64 `json:"unitPrice"      dc:"单价"`
	Amount         float64 `json:"amount"         dc:"金额"`
	UsedQuantity   float64 `json:"usedQuantity"   dc:"已领用数量（出库）"`
	RemainQuantity float64 `json:"remainQuantity" dc:"当前结存数量"`
	Supplier       string  `json:"supplier"       dc:"供应商"`
	Remark         string  `json:"remark"         dc:"备注"`
	UpdatedAt      string  `json:"updatedAt"      dc:"最后变动时间"`
}

// StockSummary 库存汇总
type StockSummary struct {
	Kinds    int     `json:"kinds"    dc:"库存种类数"`
	Amount   float64 `json:"amount"   dc:"库存总金额"`
	Shortage int     `json:"shortage" dc:"结存为 0 的种类数"`
}

// StockListRes 库存分页列表响应
type StockListRes struct {
	Current int           `json:"current" dc:"当前页码"`
	Size    int           `json:"size"    dc:"每页条数"`
	Total   int           `json:"total"   dc:"总条数"`
	Summary *StockSummary `json:"summary" dc:"库存汇总"`
	Records []*StockItem  `json:"records" dc:"库存列表"`
}

// StockForm 库存（材料）表单字段
type StockForm struct {
	ProjectId   int64   `json:"projectId"   v:"required#请选择项目"                              dc:"项目ID"`
	WarehouseId int64   `json:"warehouseId" v:"min:1#请选择仓库"                                 dc:"仓库ID，0 表示默认仓库"`
	Name        string  `json:"name"        v:"required|max-length:100#请填写材料名称|材料名称最多 100 个字符" dc:"材料名称"`
	Spec        string  `json:"spec"        v:"max-length:100#规格最多 100 个字符"             dc:"规格 / 型号"`
	Unit        string  `json:"unit"        v:"max-length:20#单位最多 20 个字符"               dc:"单位"`
	Quantity    float64 `json:"quantity"    v:"min:0#数量不能为负数"                          dc:"库存数量"`
	UnitPrice   float64 `json:"unitPrice"   v:"min:0#单价不能为负数"                          dc:"单价"`
	Supplier    string  `json:"supplier"    v:"max-length:100#供应商最多 100 个字符"            dc:"供应商"`
	Remark      string  `json:"remark"      v:"max-length:255#备注最多 255 个字符"             dc:"备注"`
}

// StockCreateReq 新增库存请求
type StockCreateReq struct {
	g.Meta `path:"/inventory/stock/create" method:"post" tags:"库存管理" summary:"新增库存"`
	StockForm
}

// StockCreateRes 新增库存响应
type StockCreateRes struct {
	Id int64 `json:"id" dc:"库存ID"`
}

// StockUpdateReq 更新库存请求
type StockUpdateReq struct {
	g.Meta `path:"/inventory/stock/update" method:"post" tags:"库存管理" summary:"更新库存"`
	Id     int64 `json:"id" v:"required#库存ID不能为空" dc:"库存ID"`
	StockForm
}

// StockUpdateRes 更新库存响应
type StockUpdateRes struct {
	Id int64 `json:"id" dc:"库存ID"`
}

// StockDeleteReq 删除库存请求
type StockDeleteReq struct {
	g.Meta `path:"/inventory/stock/delete" method:"post" tags:"库存管理" summary:"删除库存"`
	Id     int64 `json:"id" v:"required#库存ID不能为空" dc:"库存ID"`
}

// StockDeleteRes 删除库存响应
type StockDeleteRes struct{}

// StockInReq 材料入库请求：库存数量增加，可选择同时生成一张付款单。
type StockInReq struct {
	g.Meta `path:"/inventory/stock/in" method:"post" tags:"库存管理" summary:"材料入库"`

	ProjectId       int64   `json:"projectId"       v:"required#请选择项目"   dc:"项目ID"`
	MaterialId      int64   `json:"materialId"      v:"required#请选择材料"   dc:"材料ID（材料清单中的 item）"`
	Quantity        float64 `json:"quantity"        v:"required|gt:0#请填写入库数量|入库数量必须大于 0" dc:"入库数量"`
	UnitPrice       float64 `json:"unitPrice"       v:"min:0#单价不能为负数"  dc:"入库单价，0 表示使用材料清单单价"`
	GenerateFinance bool    `json:"generateFinance" dc:"是否生成付款单"`
	VendorId        int64   `json:"vendorId"        dc:"供应商ID（生成付款单时必填）"`
	FundAccountId   int64   `json:"fundAccountId"   dc:"资金账户ID（生成付款单时必填）"`
	BizDate         string  `json:"bizDate"         v:"required#请选择入库日期" dc:"入库日期 YYYY-MM-DD"`
	Remark          string  `json:"remark"          v:"max-length:255#备注最多 255 个字符" dc:"备注"`
}

// StockInRes 材料入库响应
type StockInRes struct {
	Id        int64  `json:"id"        dc:"出入库流水ID"`
	PayoutId  int64  `json:"payoutId"  dc:"生成的付款单ID，未生成时为 0"`
	VoucherNo string `json:"voucherNo" dc:"付款凭证号，未记账时为空"`
}

// StockOutReq 材料出库（项目领用）请求：写入领用记录并扣减结存，不生成任何财务单据。
type StockOutReq struct {
	g.Meta `path:"/inventory/stock/out" method:"post" tags:"库存管理" summary:"材料出库"`

	ProjectId  int64   `json:"projectId"  v:"required#请选择项目"  dc:"项目ID"`
	MaterialId int64   `json:"materialId" v:"required#请选择材料"  dc:"材料ID"`
	Quantity   float64 `json:"quantity"   v:"required|gt:0#请填写出库数量|出库数量必须大于 0" dc:"出库数量"`
	UnitPrice  float64 `json:"unitPrice"  v:"min:0#单价不能为负数"  dc:"出库单价（领用成本单价），0 表示使用材料清单单价"`
	Receiver   string  `json:"receiver"   v:"max-length:50#领用人最多 50 个字符" dc:"领用人"`
	BizDate    string  `json:"bizDate"    v:"required#请选择出库日期" dc:"出库日期 YYYY-MM-DD"`
	Remark     string  `json:"remark"     v:"max-length:255#备注最多 255 个字符" dc:"备注"`
}

// StockOutRes 材料出库响应
type StockOutRes struct {
	Id int64 `json:"id" dc:"出入库流水ID"`
}

// StockAdjustReq 库存补录请求：漏补历史数据时直接调整库存，不生成财务单据。
type StockAdjustReq struct {
	g.Meta `path:"/inventory/stock/adjust" method:"post" tags:"库存管理" summary:"库存补录"`

	RecordType string  `json:"recordType"  v:"required|in:IN,OUT#请选择补录方向|补录方向无效" dc:"补录方向 IN=补录入库 OUT=补录出库"`
	ProjectId  int64   `json:"projectId"   v:"required#请选择项目"   dc:"项目ID"`
	MaterialId int64   `json:"materialId"  v:"required#请选择材料"   dc:"材料ID"`
	Quantity   float64 `json:"quantity"    v:"required|gt:0#请填写数量|数量必须大于 0" dc:"调整数量"`
	UnitPrice  float64 `json:"unitPrice"   v:"min:0#单价不能为负数"  dc:"单价（补录入库时使用，0 表示材料清单单价）"`
	Receiver   string  `json:"receiver"    v:"max-length:50#领用人最多 50 个字符" dc:"领用人（补录出库时使用）"`
	BizDate    string  `json:"bizDate"     v:"required#请选择补录日期" dc:"补录日期 YYYY-MM-DD"`
	Remark     string  `json:"remark"      v:"max-length:255#备注最多 255 个字符" dc:"备注"`
}

// StockAdjustRes 库存补录响应
type StockAdjustRes struct {
	Id int64 `json:"id" dc:"出入库流水ID"`
}

// StockRecordItem 出入库流水行
type StockRecordItem struct {
	Id               int64   `json:"id"               dc:"流水ID"`
	RecordType       string  `json:"recordType"       dc:"IN=入库 OUT=出库"`
	BizType          string  `json:"bizType"          dc:"NORMAL=业务单据 ADJUST=补录"`
	ProjectName      string  `json:"projectName"      dc:"项目名称"`
	MaterialName     string  `json:"materialName"     dc:"材料名称"`
	Spec             string  `json:"spec"             dc:"规格 / 型号"`
	Unit             string  `json:"unit"             dc:"单位"`
	WarehouseName    string  `json:"warehouseName"    dc:"仓库名称"`
	Quantity         float64 `json:"quantity"         dc:"数量"`
	UnitPrice        float64 `json:"unitPrice"        dc:"单价"`
	Amount           float64 `json:"amount"           dc:"金额"`
	VendorName       string  `json:"vendorName"       dc:"供应商名称"`
	FundAccountName  string  `json:"fundAccountName"  dc:"资金账户名称"`
	PayoutId         int64   `json:"payoutId"         dc:"关联付款单ID"`
	PaymentId        int64   `json:"paymentId"        dc:"关联收款单ID"`
	PayoutVoucherNo  string  `json:"payoutVoucherNo"  dc:"付款凭证号"`
	PaymentVoucherNo string  `json:"paymentVoucherNo" dc:"收款凭证号"`
	BizDate          string  `json:"bizDate"          dc:"业务日期"`
	Remark           string  `json:"remark"           dc:"备注"`
	CreatedAt        string  `json:"createdAt"        dc:"创建时间"`
}

// StockRecordListReq 出入库流水分页请求
type StockRecordListReq struct {
	g.Meta `path:"/inventory/stock/records" method:"get" tags:"库存管理" summary:"出入库流水列表"`

	ProjectId  int64  `json:"projectId"  dc:"项目ID，0 表示全部"`
	MaterialId int64  `json:"materialId" dc:"材料ID，0 表示全部"`
	RecordType string `json:"recordType" dc:"IN / OUT，空表示全部"`
	Current    int    `json:"current"    d:"1"  dc:"当前页码"`
	Size       int    `json:"size"       d:"10" dc:"每页条数"`
}

// StockRecordListRes 出入库流水分页响应
type StockRecordListRes struct {
	Current int                `json:"current" dc:"当前页码"`
	Size    int                `json:"size"    dc:"每页条数"`
	Total   int                `json:"total"   dc:"总条数"`
	Records []*StockRecordItem `json:"records" dc:"流水列表"`
}
