declare namespace Api {
  /**
   * namespace Inventory
   *
   * backend api module: "inventory"（库存管理）
   */
  namespace Inventory {
    /** 公司级采购/销售材料基础档案 */
    interface MaterialCatalog {
      id: number;
      name: string;
      spec: string;
      unit: string;
      unitPrice: number;
      taxMode: 'NO_INVOICE' | 'TAX_INCLUDED' | 'TAX_EXCLUDED';
      taxRate: number;
      remark: string;
      createdAt: string;
      updatedAt: string;
    }

    interface MaterialCatalogForm {
      name: string;
      spec: string;
      unit: string;
      unitPrice: number | null;
      taxMode: 'NO_INVOICE' | 'TAX_INCLUDED' | 'TAX_EXCLUDED';
      taxRate: number | null;
      remark: string;
    }

    interface MaterialCatalogSearchParams {
      keyword?: string | null;
      current?: number;
      size?: number;
    }

    type MaterialCatalogList = Api.Common.PaginatingQueryRecord<MaterialCatalog>;

    /** 库存行（来源：project_material + 领用汇总） */
    interface Stock {
      id: number;
      projectId: number;
      projectName: string;
      warehouseId: number;
      warehouseName: string;
      name: string;
      spec: string;
      unit: string;
      /** 库存数量（累计入库） */
      quantity: number;
      unitPrice: number;
      amount: number;
      /** 已领用数量（出库） */
      usedQuantity: number;
      /** 当前结存数量 */
      remainQuantity: number;
      supplier: string;
      remark: string;
      /** 最后变动时间 */
      updatedAt: string;
    }

    /** 库存汇总 */
    interface StockSummary {
      /** 库存种类数 */
      kinds: number;
      /** 库存总金额 */
      amount: number;
      /** 结存为 0 的种类数 */
      shortage: number;
    }

    /** 库存查询参数 */
    interface StockSearchParams {
      /** 0 表示全部项目 */
      projectId?: number;
      /** 0 表示全部仓库 */
      warehouseId?: number;
      keyword?: string | null;
      current?: number;
      size?: number;
    }

    interface StockList {
      current: number;
      size: number;
      total: number;
      summary: StockSummary;
      records: Stock[];
    }

    /** 库存保存参数 */
    interface StockForm {
      projectId: number | null;
      warehouseId: number | null;
      name: string;
      spec: string;
      unit: string;
      quantity: number | null;
      unitPrice: number | null;
      supplier: string;
      remark: string;
    }

    /** 仓库行 */
    interface Warehouse {
      id: number;
      name: string;
      code: string;
      address: string;
      remark: string;
      /** 是否默认仓库 */
      isDefault: boolean;
      /** 排序（越小越靠前） */
      sort: number;
      createdAt: string;
      updatedAt: string;
    }

    /** 仓库查询参数 */
    interface WarehouseSearchParams extends Api.Common.CommonSearchParams {
      keyword?: string | null;
    }

    type WarehouseList = Api.Common.PaginatingQueryRecord<Warehouse>;

    /** 仓库保存参数 */
    interface WarehouseForm {
      name: string;
      code: string;
      address: string;
      remark: string;
      isDefault: boolean;
      sort: number | null;
    }

    /** 材料入库参数（可选联动生成付款单） */
    interface StockInForm {
      projectId: number | null;
      materialId: number | null;
      quantity: number | null;
      /** 0 表示使用材料清单单价 */
      unitPrice: number | null;
      generateFinance: boolean;
      vendorId: number | null;
      fundAccountId: number | null;
      bizDate: string;
      remark: string;
    }

    /** 材料出库参数（项目领用：只扣库存与计入项目材料成本，不生成财务单据） */
    interface StockOutForm {
      projectId: number | null;
      materialId: number | null;
      quantity: number | null;
      /** 出库单价（领用成本单价），0 表示使用材料清单单价 */
      unitPrice: number | null;
      receiver: string;
      bizDate: string;
      remark: string;
    }

    /** 库存补录参数（不生成财务单据） */
    interface StockAdjustForm {
      recordType: 'IN' | 'OUT';
      projectId: number | null;
      materialId: number | null;
      quantity: number | null;
      unitPrice: number | null;
      receiver: string;
      bizDate: string;
      remark: string;
    }

    /** 出入库流水行 */
    interface StockRecord {
      id: number;
      /** IN=入库 OUT=出库 */
      recordType: 'IN' | 'OUT';
      /** NORMAL=业务单据 ADJUST=补录 */
      bizType: 'NORMAL' | 'ADJUST';
      projectName: string;
      materialName: string;
      spec: string;
      unit: string;
      warehouseName: string;
      quantity: number;
      unitPrice: number;
      amount: number;
      vendorName: string;
      fundAccountName: string;
      payoutId: number;
      paymentId: number;
      payoutVoucherNo: string;
      paymentVoucherNo: string;
      bizDate: string;
      remark: string;
      createdAt: string;
    }

    /** 出入库流水查询参数 */
    interface StockRecordSearchParams {
      projectId?: number;
      materialId?: number;
      recordType?: string | null;
      current?: number;
      size?: number;
    }

    type StockRecordList = Api.Common.PaginatingQueryRecord<StockRecord>;
  }
}
