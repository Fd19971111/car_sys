package v1

import "github.com/gogf/gf/v2/frame/g"

// ===================== 材料清单 =====================

// MaterialListReq 材料清单分页列表请求
type MaterialListReq struct {
	g.Meta `path:"/project/material/list" method:"get" tags:"项目资料" summary:"材料清单列表"`

	ProjectId int64  `json:"projectId" v:"required#请选择项目" dc:"项目ID"`
	Keyword   string `json:"keyword"   dc:"材料名称 / 规格 / 供应商关键字"`
	Current   int    `json:"current"   d:"1"  dc:"当前页码"`
	Size      int    `json:"size"      d:"10" dc:"每页条数"`
}

// MaterialItem 材料清单行
type MaterialItem struct {
	Id             int64   `json:"id"             dc:"材料ID"`
	CatalogId      int64   `json:"catalogId"      dc:"公司级材料档案ID"`
	ProjectId      int64   `json:"projectId"      dc:"项目ID"`
	WarehouseId    int64   `json:"warehouseId"    dc:"仓库ID"`
	WarehouseName  string  `json:"warehouseName"  dc:"仓库名称"`
	Name           string  `json:"name"           dc:"材料名称"`
	Spec           string  `json:"spec"           dc:"规格 / 型号"`
	Unit           string  `json:"unit"           dc:"单位"`
	Quantity       float64 `json:"quantity"       dc:"清单数量"`
	UnitPrice      float64 `json:"unitPrice"      dc:"单价"`
	TaxMode        string  `json:"taxMode"        dc:"税务方式"`
	TaxRate        float64 `json:"taxRate"        dc:"税率（百分比）"`
	Amount         float64 `json:"amount"         dc:"金额"`
	UsedQuantity   float64 `json:"usedQuantity"   dc:"已领用数量"`
	RemainQuantity float64 `json:"remainQuantity" dc:"剩余数量"`
	Supplier       string  `json:"supplier"       dc:"供应商"`
	Remark         string  `json:"remark"         dc:"备注"`
	CreatedAt      string  `json:"createdAt"      dc:"创建时间"`
}

// MaterialListRes 材料清单分页列表响应
type MaterialListRes struct {
	Current int             `json:"current" dc:"当前页码"`
	Size    int             `json:"size"    dc:"每页条数"`
	Total   int             `json:"total"   dc:"总条数"`
	Records []*MaterialItem `json:"records" dc:"材料清单"`
}

// MaterialForm 材料清单表单字段
type MaterialForm struct {
	ProjectId   int64   `json:"projectId"   v:"required#请选择项目"                              dc:"项目ID"`
	CatalogId   int64   `json:"catalogId"   dc:"公司级材料档案ID"`
	WarehouseId int64   `json:"warehouseId" v:"min:1#请选择仓库"                                 dc:"仓库ID，0 表示默认仓库"`
	Name        string  `json:"name"        v:"required|max-length:100#请填写材料名称|材料名称最多 100 个字符" dc:"材料名称"`
	Spec        string  `json:"spec"        v:"max-length:100#规格最多 100 个字符"             dc:"规格 / 型号"`
	Unit        string  `json:"unit"        v:"max-length:20#单位最多 20 个字符"               dc:"单位"`
	Quantity    float64 `json:"quantity"    v:"min:0#数量不能为负数"                          dc:"清单数量"`
	UnitPrice   float64 `json:"unitPrice"   v:"min:0#单价不能为负数"                          dc:"单价"`
	TaxMode     string  `json:"taxMode"     v:"required|in:NO_INVOICE,TAX_INCLUDED,TAX_EXCLUDED#请选择税务方式|税务方式无效" dc:"税务方式"`
	TaxRate     float64 `json:"taxRate"     v:"min:0|max:100#税率不能超过 100%"                 dc:"税率（百分比）"`
	Supplier    string  `json:"supplier"    v:"max-length:100#供应商最多 100 个字符"            dc:"供应商"`
	Remark      string  `json:"remark"      v:"max-length:255#备注最多 255 个字符"             dc:"备注"`
}

// MaterialCreateReq 新增材料请求
type MaterialCreateReq struct {
	g.Meta `path:"/project/material/create" method:"post" tags:"项目资料" summary:"新增材料"`
	MaterialForm
}

// MaterialCreateRes 新增材料响应
type MaterialCreateRes struct {
	Id int64 `json:"id" dc:"材料ID"`
}

// MaterialUpdateReq 更新材料请求
type MaterialUpdateReq struct {
	g.Meta `path:"/project/material/update" method:"post" tags:"项目资料" summary:"更新材料"`
	Id     int64 `json:"id" v:"required#材料ID不能为空" dc:"材料ID"`
	MaterialForm
}

// MaterialUpdateRes 更新材料响应
type MaterialUpdateRes struct {
	Id int64 `json:"id" dc:"材料ID"`
}

// MaterialDeleteReq 删除材料请求
type MaterialDeleteReq struct {
	g.Meta `path:"/project/material/delete" method:"post" tags:"项目资料" summary:"删除材料"`
	Id     int64 `json:"id" v:"required#材料ID不能为空" dc:"材料ID"`
}

// MaterialDeleteRes 删除材料响应
type MaterialDeleteRes struct{}

// ===================== 材料领用记录 =====================

// MaterialUsageListReq 材料领用记录请求
type MaterialUsageListReq struct {
	g.Meta `path:"/project/material/usageList" method:"get" tags:"项目资料" summary:"材料领用记录"`

	ProjectId  int64 `json:"projectId"  v:"required#请选择项目" dc:"项目ID"`
	MaterialId int64 `json:"materialId" dc:"材料ID，0 表示全部"`
	Current    int   `json:"current"    d:"1"  dc:"当前页码"`
	Size       int   `json:"size"       d:"10" dc:"每页条数"`
}

// MaterialUsageItem 领用记录行
type MaterialUsageItem struct {
	Id            int64   `json:"id"            dc:"领用记录ID"`
	ProjectId     int64   `json:"projectId"     dc:"项目ID"`
	MaterialId    int64   `json:"materialId"    dc:"材料ID"`
	MaterialName  string  `json:"materialName"  dc:"材料名称"`
	WarehouseId   int64   `json:"warehouseId"   dc:"出库仓库ID，0 表示沿用材料所属仓库"`
	WarehouseName string  `json:"warehouseName" dc:"仓库名称"`
	Quantity      float64 `json:"quantity"      dc:"领用数量"`
	UsedAt        string  `json:"usedAt"        dc:"领用日期"`
	Receiver      string  `json:"receiver"      dc:"领用人"`
	Remark        string  `json:"remark"        dc:"备注"`
	StockedAt     string  `json:"stockedAt"     dc:"出库时间，空表示未出入库"`
	CreatedBy     int64   `json:"createdBy"     dc:"创建人ID，0 表示历史数据"`
	CreatedAt     string  `json:"createdAt"     dc:"登记时间"`
}

// MaterialUsageListRes 领用记录响应
type MaterialUsageListRes struct {
	Current int                  `json:"current" dc:"当前页码"`
	Size    int                  `json:"size"    dc:"每页条数"`
	Total   int                  `json:"total"   dc:"总条数"`
	Records []*MaterialUsageItem `json:"records" dc:"领用记录"`
}

// MaterialUsageForm 领用记录表单字段
type MaterialUsageForm struct {
	ProjectId  int64   `json:"projectId"  v:"required#请选择项目"                     dc:"项目ID"`
	MaterialId int64   `json:"materialId" v:"required#请选择材料"                     dc:"材料ID"`
	Quantity   float64 `json:"quantity"   v:"required|min:0.01#请填写领用数量|领用数量必须大于 0" dc:"领用数量"`
	UsedAt     string  `json:"usedAt"     v:"required#请选择领用日期"                  dc:"领用日期（YYYY-MM-DD）"`
	Receiver   string  `json:"receiver"   v:"max-length:50#领用人最多 50 个字符"        dc:"领用人"`
	Remark     string  `json:"remark"     v:"max-length:255#备注最多 255 个字符"        dc:"备注"`
}

// MaterialUsageCreateReq 新增领用记录请求
type MaterialUsageCreateReq struct {
	g.Meta `path:"/project/material/usageCreate" method:"post" tags:"项目资料" summary:"新增材料领用"`
	MaterialUsageForm
}

// MaterialUsageCreateRes 新增领用记录响应
type MaterialUsageCreateRes struct {
	Id int64 `json:"id" dc:"领用记录ID"`
}

// MaterialUsageDeleteReq 删除领用记录请求
type MaterialUsageDeleteReq struct {
	g.Meta `path:"/project/material/usageDelete" method:"post" tags:"项目资料" summary:"删除材料领用"`
	Id     int64 `json:"id" v:"required#领用记录ID不能为空" dc:"领用记录ID"`
}

// MaterialUsageDeleteRes 删除领用记录响应
type MaterialUsageDeleteRes struct{}

// MaterialUsageStockReq 领用单出库请求：选择仓库后直接补出库流水。
type MaterialUsageStockReq struct {
	g.Meta `path:"/project/material/usageStock" method:"post" tags:"库存" summary:"领用单出库"`

	Id          int64  `json:"id"          v:"required#领用记录ID不能为空" dc:"领用记录ID"`
	WarehouseId int64  `json:"warehouseId" v:"required|min:1#请选择仓库"    dc:"出库仓库ID"`
	Remark      string `json:"remark"      v:"max-length:255#备注最多 255 个字符" dc:"备注"`
}

// MaterialUsageStockRes 领用单出库响应
type MaterialUsageStockRes struct {
	Id       int64 `json:"id"       dc:"领用记录ID"`
	RecordId int64 `json:"recordId" dc:"出入库流水ID"`
}

// ===================== 材料采购单 =====================

type PurchaseOrderListReq struct {
	g.Meta `path:"/project/purchaseOrder/list" method:"get" tags:"库存" summary:"采购单列表"`

	ProjectId int64 `json:"projectId" dc:"项目ID，0 表示全部项目"`
	Current   int   `json:"current" d:"1" dc:"当前页码"`
	Size      int   `json:"size" d:"10" dc:"每页条数"`
}

type PurchaseOrderPendingListReq struct {
	g.Meta `path:"/project/purchaseOrder/pendingList" method:"get" tags:"库存" summary:"待付款采购单列表"`
}

type PurchaseOrderItem struct {
	Id            int64   `json:"id"`
	OrderNo       string  `json:"orderNo"`
	ProjectId     int64   `json:"projectId"`
	ProjectName   string  `json:"projectName"`
	MaterialId    int64   `json:"materialId"`
	VendorId      int64   `json:"vendorId"`
	PurchaseDate  string  `json:"purchaseDate"`
	MaterialName  string  `json:"materialName"`
	Spec          string  `json:"spec"`
	Unit          string  `json:"unit"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unitPrice"`
	TaxMode       string  `json:"taxMode"`
	TaxRate       float64 `json:"taxRate"`
	NetAmount     float64 `json:"netAmount"`
	TaxAmount     float64 `json:"taxAmount"`
	Amount        float64 `json:"amount"`
	Supplier      string  `json:"supplier"`
	Remark        string  `json:"remark"`
	Status        string  `json:"status"`
	PayoutId      int64   `json:"payoutId"`
	WarehouseId   int64   `json:"warehouseId"`
	WarehouseName string  `json:"warehouseName"`
	StockedAt     string  `json:"stockedAt"`
	CreatedBy     int64   `json:"createdBy"`
	CreatedAt     string  `json:"createdAt"`
}

type PurchaseOrderListRes struct {
	Current int                  `json:"current"`
	Size    int                  `json:"size"`
	Total   int                  `json:"total"`
	Records []*PurchaseOrderItem `json:"records"`
}

type PurchaseOrderPendingListRes struct {
	Records []*PurchaseOrderItem `json:"records"`
}

type PurchaseOrderCreateReq struct {
	g.Meta `path:"/project/purchaseOrder/create" method:"post" tags:"库存" summary:"创建采购单"`

	ProjectId    int64   `json:"projectId" v:"required#请选择项目"`
	MaterialId   int64   `json:"materialId" v:"required#请选择材料"`
	VendorId     int64   `json:"vendorId" v:"required#请选择供应商"`
	PurchaseDate string  `json:"purchaseDate" v:"required#请选择采购日期"`
	Quantity     float64 `json:"quantity" v:"required|min:0.01#请输入采购数量|采购数量必须大于 0"`
	UnitPrice    float64 `json:"unitPrice" v:"min:0#采购单价不能为负数"`
	TaxMode      string  `json:"taxMode" v:"required|in:NO_INVOICE,TAX_INCLUDED,TAX_EXCLUDED#请选择税务方式|税务方式无效"`
	TaxRate      float64 `json:"taxRate" v:"min:0|max:100#税率不能超过 100%"`
	Remark       string  `json:"remark" v:"max-length:255#备注最多 255 个字符"`
}

type PurchaseOrderCreateRes struct {
	Id int64 `json:"id"`
}

// PurchaseOrderStockReq 采购单入库请求：选择仓库后直接完成入库。
type PurchaseOrderStockReq struct {
	g.Meta `path:"/project/purchaseOrder/stock" method:"post" tags:"库存" summary:"采购单入库"`

	Id          int64  `json:"id"          v:"required#采购单ID不能为空" dc:"采购单ID"`
	WarehouseId int64  `json:"warehouseId" v:"required|min:1#请选择仓库"   dc:"入库仓库ID"`
	Remark      string `json:"remark"      v:"max-length:255#备注最多 255 个字符" dc:"备注"`
}

// PurchaseOrderStockRes 采购单入库响应
type PurchaseOrderStockRes struct {
	Id       int64 `json:"id"       dc:"采购单ID"`
	RecordId int64 `json:"recordId" dc:"出入库流水ID"`
}

// ===================== 施工日志 / 现场照片 =====================

// SiteLogListReq 施工日志列表请求
type SiteLogListReq struct {
	g.Meta `path:"/project/log/list" method:"get" tags:"项目资料" summary:"施工日志列表"`

	ProjectId int64 `json:"projectId" v:"required#请选择项目" dc:"项目ID"`
	Current   int   `json:"current"   d:"1"  dc:"当前页码"`
	Size      int   `json:"size"      d:"10" dc:"每页条数"`
}

// SiteLogItem 施工日志行
type SiteLogItem struct {
	Id        int64    `json:"id"        dc:"日志ID"`
	ProjectId int64    `json:"projectId" dc:"项目ID"`
	Content   string   `json:"content"   dc:"日志内容"`
	Images    []string `json:"images"    dc:"图片相对路径（编辑时回传）"`
	ImageUrls []string `json:"imageUrls" dc:"图片绝对地址（预览用）"`
	LogDate   string   `json:"logDate"   dc:"日志业务日期（YYYY-MM-DD，日志实际发生日期）"`
	UserName  string   `json:"userName"  dc:"记录人"`
	CreatedAt string   `json:"createdAt" dc:"记录时间"`
}

// SiteLogListRes 施工日志列表响应
type SiteLogListRes struct {
	Current int            `json:"current" dc:"当前页码"`
	Size    int            `json:"size"    dc:"每页条数"`
	Total   int            `json:"total"   dc:"总条数"`
	Records []*SiteLogItem `json:"records" dc:"施工日志"`
}

// SiteLogCreateReq 新增施工日志请求
type SiteLogCreateReq struct {
	g.Meta `path:"/project/log/create" method:"post" tags:"项目资料" summary:"新增施工日志"`

	ProjectId int64    `json:"projectId" v:"required#请选择项目"                                    dc:"项目ID"`
	Content   string   `json:"content"   v:"required|max-length:500#请填写日志内容|日志内容最多 500 个字符" dc:"日志内容"`
	Images    []string `json:"images"    dc:"图片相对路径列表，最多 9 张"`
	LogDate   string   `json:"logDate"   dc:"日志业务日期（YYYY-MM-DD），为空默认当天；补录历史日志时传真实日期"`
}

// SiteLogCreateRes 新增施工日志响应
type SiteLogCreateRes struct {
	Id int64 `json:"id" dc:"日志ID"`
}

// SiteLogUpdateReq 更新施工日志请求
type SiteLogUpdateReq struct {
	g.Meta `path:"/project/log/update" method:"post" tags:"项目资料" summary:"更新施工日志"`

	Id      int64    `json:"id"      v:"required#日志ID不能为空"                                      dc:"日志ID"`
	Content string   `json:"content" v:"required|max-length:500#请填写日志内容|日志内容最多 500 个字符" dc:"日志内容"`
	Images  []string `json:"images"  dc:"图片相对路径列表，最多 9 张"`
}

// SiteLogUpdateRes 更新施工日志响应
type SiteLogUpdateRes struct {
	Id int64 `json:"id" dc:"日志ID"`
}

// SiteLogDeleteReq 删除施工日志请求
type SiteLogDeleteReq struct {
	g.Meta `path:"/project/log/delete" method:"post" tags:"项目资料" summary:"删除施工日志"`
	Id     int64 `json:"id" v:"required#日志ID不能为空" dc:"日志ID"`
}

// SiteLogDeleteRes 删除施工日志响应
type SiteLogDeleteRes struct{}

// ===================== 变更记录（增项 / 减项） =====================

// ChangeListReq 变更记录分页列表请求
type ChangeListReq struct {
	g.Meta `path:"/project/change/list" method:"get" tags:"项目资料" summary:"变更记录列表"`

	ProjectId  int64  `json:"projectId"  v:"required#请选择项目" dc:"项目ID"`
	ChangeType string `json:"changeType" dc:"变更类型 ADD/REDUCE，空为全部"`
	Keyword    string `json:"keyword"    dc:"变更事项 / 原因关键字"`
	Current    int    `json:"current"    d:"1"  dc:"当前页码"`
	Size       int    `json:"size"       d:"10" dc:"每页条数"`
}

// ChangeItem 变更记录行
type ChangeItem struct {
	Id         int64    `json:"id"         dc:"变更ID"`
	ProjectId  int64    `json:"projectId"  dc:"项目ID"`
	ChangeType string   `json:"changeType" dc:"变更类型 ADD 增项 / REDUCE 减项"`
	Title      string   `json:"title"      dc:"变更事项"`
	Amount     float64  `json:"amount"     dc:"变更金额"`
	ChangeDate string   `json:"changeDate" dc:"变更日期"`
	Reason     string   `json:"reason"     dc:"变更原因"`
	Images     []string `json:"images"     dc:"图片相对路径"`
	ImageUrls  []string `json:"imageUrls"  dc:"图片绝对地址"`
	Remark     string   `json:"remark"     dc:"备注"`
	CreatedAt  string   `json:"createdAt"  dc:"创建时间"`
}

// ChangeListRes 变更记录列表响应
type ChangeListRes struct {
	Current     int           `json:"current"     dc:"当前页码"`
	Size        int           `json:"size"        dc:"每页条数"`
	Total       int           `json:"total"       dc:"总条数"`
	AddTotal    float64       `json:"addTotal"    dc:"增项合计"`
	ReduceTotal float64       `json:"reduceTotal" dc:"减项合计"`
	Records     []*ChangeItem `json:"records"     dc:"变更记录"`
}

// ChangeForm 变更记录表单字段
type ChangeForm struct {
	ProjectId  int64    `json:"projectId"  v:"required#请选择项目"                             dc:"项目ID"`
	ChangeType string   `json:"changeType" v:"required|in:ADD,REDUCE#请选择变更类型|变更类型不正确" dc:"变更类型：ADD 增项 / REDUCE 减项"`
	Title      string   `json:"title"      v:"required|max-length:100#请填写变更事项|变更事项最多 100 个字符" dc:"变更事项"`
	Amount     float64  `json:"amount"     v:"min:0#变更金额不能为负数"                        dc:"变更金额"`
	ChangeDate string   `json:"changeDate" v:"required#请选择变更日期"                          dc:"变更日期（YYYY-MM-DD）"`
	Reason     string   `json:"reason"     v:"max-length:500#变更原因最多 500 个字符"           dc:"变更原因"`
	Images     []string `json:"images"     dc:"附件图片相对路径，最多 9 张"`
	Remark     string   `json:"remark"     v:"max-length:255#备注最多 255 个字符"              dc:"备注"`
}

// ChangeCreateReq 新增变更记录请求
type ChangeCreateReq struct {
	g.Meta `path:"/project/change/create" method:"post" tags:"项目资料" summary:"新增变更记录"`
	ChangeForm
}

// ChangeCreateRes 新增变更记录响应
type ChangeCreateRes struct {
	Id int64 `json:"id" dc:"变更ID"`
}

// ChangeUpdateReq 更新变更记录请求
type ChangeUpdateReq struct {
	g.Meta `path:"/project/change/update" method:"post" tags:"项目资料" summary:"更新变更记录"`
	Id     int64 `json:"id" v:"required#变更ID不能为空" dc:"变更ID"`
	ChangeForm
}

// ChangeUpdateRes 更新变更记录响应
type ChangeUpdateRes struct {
	Id int64 `json:"id" dc:"变更ID"`
}

// ChangeDeleteReq 删除变更记录请求
type ChangeDeleteReq struct {
	g.Meta `path:"/project/change/delete" method:"post" tags:"项目资料" summary:"删除变更记录"`
	Id     int64 `json:"id" v:"required#变更ID不能为空" dc:"变更ID"`
}

// ChangeDeleteRes 删除变更记录响应
type ChangeDeleteRes struct{}

// ===================== 项目合同 =====================

type ContractAttachment struct {
	Path string `json:"path" dc:"附件相对路径"`
	Name string `json:"name" dc:"附件原始名称"`
	Url  string `json:"url,omitempty" dc:"附件访问地址"`
}

type ContractListReq struct {
	g.Meta    `path:"/project/contract/list" method:"get" tags:"项目资料" summary:"项目合同列表"`
	ProjectId int64 `json:"projectId" v:"required#请选择项目" dc:"项目ID"`
	Current   int   `json:"current" d:"1" dc:"当前页码"`
	Size      int   `json:"size" d:"10" dc:"每页条数"`
}

type ContractItem struct {
	Id          int64                 `json:"id" dc:"合同ID"`
	ProjectId   int64                 `json:"projectId" dc:"项目ID"`
	Name        string                `json:"name" dc:"合同名称"`
	ContractNo  string                `json:"contractNo" dc:"合同编号"`
	Amount      float64               `json:"amount" dc:"合同金额"`
	SignedAt    string                `json:"signedAt" dc:"签订日期"`
	Remark      string                `json:"remark" dc:"备注"`
	Attachments []*ContractAttachment `json:"attachments" dc:"合同附件"`
	CreatedAt   string                `json:"createdAt" dc:"创建时间"`
}

type ContractListRes struct {
	Current int             `json:"current" dc:"当前页码"`
	Size    int             `json:"size" dc:"每页条数"`
	Total   int             `json:"total" dc:"总条数"`
	Records []*ContractItem `json:"records" dc:"合同列表"`
}

type ContractForm struct {
	ProjectId   int64                 `json:"projectId" v:"required#请选择项目" dc:"项目ID"`
	Name        string                `json:"name" v:"required|max-length:100#请填写合同名称|合同名称最多100个字符" dc:"合同名称"`
	ContractNo  string                `json:"contractNo" v:"max-length:100#合同编号最多100个字符" dc:"合同编号"`
	Amount      float64               `json:"amount" v:"min:0#合同金额不能为负数" dc:"合同金额"`
	SignedAt    string                `json:"signedAt" v:"required#请选择签订日期" dc:"签订日期"`
	Remark      string                `json:"remark" v:"max-length:500#备注最多500个字符" dc:"备注"`
	Attachments []*ContractAttachment `json:"attachments" dc:"合同附件，最多9个"`
}

type ContractCreateReq struct {
	g.Meta `path:"/project/contract/create" method:"post" tags:"项目资料" summary:"新增项目合同"`
	ContractForm
}

type ContractCreateRes struct {
	Id int64 `json:"id" dc:"合同ID"`
}

type ContractUpdateReq struct {
	g.Meta `path:"/project/contract/update" method:"post" tags:"项目资料" summary:"更新项目合同"`
	Id     int64 `json:"id" v:"required#合同ID不能为空" dc:"合同ID"`
	ContractForm
}

type ContractUpdateRes struct {
	Id int64 `json:"id" dc:"合同ID"`
}

type ContractDeleteReq struct {
	g.Meta `path:"/project/contract/delete" method:"post" tags:"项目资料" summary:"删除项目合同"`
	Id     int64 `json:"id" v:"required#合同ID不能为空" dc:"合同ID"`
}

type ContractDeleteRes struct{}
