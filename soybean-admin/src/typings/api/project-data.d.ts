declare namespace Api {
  /**
   * namespace ProjectData
   *
   * backend api module: "project_data"（项目管理 - 项目资料）
   */
  namespace ProjectData {
    /** 材料清单 */
    interface Material {
      id: number;
      catalogId: number;
      projectId: number;
      warehouseId: number;
      warehouseName: string;
      name: string;
      spec: string;
      unit: string;
      quantity: number;
      unitPrice: number;
      taxMode: 'NO_INVOICE' | 'TAX_INCLUDED' | 'TAX_EXCLUDED';
      taxRate: number;
      amount: number;
      usedQuantity: number;
      remainQuantity: number;
      supplier: string;
      remark: string;
      createdAt: string;
    }

    /** 材料清单保存参数 */
    interface MaterialForm {
      projectId: number | null;
      catalogId: number | null;
      warehouseId: number | null;
      name: string;
      spec: string;
      unit: string;
      quantity: number | null;
      unitPrice: number | null;
      taxMode: 'NO_INVOICE' | 'TAX_INCLUDED' | 'TAX_EXCLUDED';
      taxRate: number | null;
      supplier: string;
      remark: string;
    }

    /** 材料清单查询参数 */
    interface MaterialSearchParams {
      projectId: number;
      keyword?: string | null;
      current?: number;
      size?: number;
    }

    /** 材料领用记录 */
    interface Usage {
      id: number;
      projectId: number;
      materialId: number;
      materialName: string;
      /** 出库仓库ID，0 表示沿用材料所属仓库 */
      warehouseId: number;
      warehouseName: string;
      quantity: number;
      usedAt: string;
      receiver: string;
      remark: string;
      /** 出库时间，空表示未出入库 */
      stockedAt: string;
      /** 创建人ID，0 表示历史数据 */
      createdBy: number;
      createdAt: string;
    }

    /** 材料采购单 */
    interface PurchaseOrder {
      id: number;
      orderNo: string;
      projectId: number;
      projectName: string;
      materialId: number;
      vendorId: number;
      purchaseDate: string;
      materialName: string;
      spec: string;
      unit: string;
      quantity: number;
      unitPrice: number;
      taxMode: 'NO_INVOICE' | 'TAX_INCLUDED' | 'TAX_EXCLUDED';
      taxRate: number;
      netAmount: number;
      taxAmount: number;
      amount: number;
      supplier: string;
      remark: string;
      status: 'DRAFT' | 'CONFIRMED' | 'PAID';
      payoutId: number;
      /** 出入库仓库ID，0 表示沿用材料所属仓库 */
      warehouseId: number;
      warehouseName: string;
      /** 入库时间，空表示未出入库 */
      stockedAt: string;
      /** 创建人ID，0 表示历史数据 */
      createdBy: number;
      createdAt: string;
    }

    interface PurchaseOrderForm {
      projectId: number;
      materialId: number;
      /** 供应商ID，未选择时为 null（0 会被下拉框当成有效值显示成「0」） */
      vendorId: number | null;
      purchaseDate: string;
      quantity: number;
      unitPrice: number;
      taxMode: 'NO_INVOICE' | 'TAX_INCLUDED' | 'TAX_EXCLUDED';
      taxRate: number;
      remark: string;
    }

    interface PurchaseOrderList {
      current: number;
      size: number;
      total: number;
      records: PurchaseOrder[];
    }

    /** 领用记录保存参数 */
    interface UsageForm {
      projectId: number | null;
      materialId: number | null;
      quantity: number | null;
      usedAt: string | null;
      receiver: string;
      remark: string;
    }

    /** 领用记录查询参数 */
    interface UsageSearchParams {
      projectId: number;
      materialId?: number | null;
      current?: number;
      size?: number;
    }

    /** 施工日志 / 现场照片 */
    interface SiteLog {
      id: number;
      projectId: number;
      content: string;
      images: string[];
      imageUrls: string[];
      /** 日志业务日期（YYYY-MM-DD，日志实际发生日期） */
      logDate: string;
      userName: string;
      createdAt: string;
    }

    /** 施工日志保存参数 */
    interface SiteLogForm {
      projectId: number | null;
      content: string;
      images: string[];
      /** 日志业务日期（YYYY-MM-DD），为空默认当天 */
      logDate?: string;
    }

    /** 施工日志查询参数 */
    interface SiteLogSearchParams {
      projectId: number;
      current?: number;
      size?: number;
    }

    interface ContractAttachment {
      path: string;
      name: string;
      url?: string;
    }

    interface Contract {
      id: number;
      projectId: number;
      name: string;
      contractNo: string;
      amount: number;
      signedAt: string;
      remark: string;
      attachments: ContractAttachment[];
      createdAt: string;
    }

    interface ContractForm {
      projectId: number | null;
      name: string;
      contractNo: string;
      amount: number | null;
      signedAt: string | null;
      remark: string;
      attachments: ContractAttachment[];
    }

    interface ContractListParams {
      projectId: number;
      current?: number;
      size?: number;
    }

    interface ContractList {
      current: number;
      size: number;
      total: number;
      records: Contract[];
    }

    /** 变更类型：增项 / 减项 */
    type ChangeType = 'ADD' | 'REDUCE';

    /** 变更记录 */
    interface Change {
      id: number;
      projectId: number;
      changeType: ChangeType;
      title: string;
      amount: number;
      changeDate: string;
      reason: string;
      images: string[];
      imageUrls: string[];
      remark: string;
      createdAt: string;
    }

    /** 变更记录保存参数 */
    interface ChangeForm {
      projectId: number | null;
      changeType: ChangeType;
      title: string;
      amount: number | null;
      changeDate: string | null;
      reason: string;
      images: string[];
      remark: string;
    }

    /** 变更记录查询参数 */
    interface ChangeSearchParams {
      projectId: number;
      changeType?: ChangeType | null;
      keyword?: string | null;
      current?: number;
      size?: number;
    }

    /** 变更记录列表结果 */
    interface ChangeList {
      current: number;
      size: number;
      total: number;
      addTotal: number;
      reduceTotal: number;
      records: Change[];
    }

    /** 分页列表结果 */
    interface MaterialList {
      current: number;
      size: number;
      total: number;
      records: Material[];
    }

    /** 领用记录列表结果 */
    interface UsageList {
      current: number;
      size: number;
      total: number;
      records: Usage[];
    }

    /** 施工日志列表结果 */
    interface SiteLogList {
      current: number;
      size: number;
      total: number;
      records: SiteLog[];
    }
  }
}
