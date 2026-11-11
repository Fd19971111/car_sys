declare namespace Api {
  /**
   * namespace Payout
   *
   * backend api module: "payout"（对应后端 project_payout 表，付款管理）
   */
  namespace Payout {
    /** 凭证分录：详情页展示「借什么、贷什么」 */
    interface LedgerEntry {
      /** 借贷方向：DEBIT 借 / CREDIT 贷 */
      direction: 'DEBIT' | 'CREDIT';
      /** 会计科目名称 */
      accountName: string;
      /** 分录金额 */
      amount: number;
      /** 分录摘要 */
      summary: string;
    }

    /** 付款记录（列表 / 详情共用） */
    interface Item {
      id: number;
      /** 所属项目ID，0 表示公司公共支出 */
      projectId: number;
      /** 项目名称 */
      projectName: string;
      /** 客户姓名 */
      customerName: string;
      /** 付款金额 */
      amount: number;
      /** 费用类型：材料款 / 人工费 / 施工费 / 设计费 / 管理费 / 销售费 / 其他 */
      payoutType: string;
      /** 付款日期 */
      payoutDate: string;
      /** 供应商ID */
      vendorId: number | null;
      /** 供应商名称 */
      vendorName: string;
      /** 关联采购单ID，0表示手工付款 */
      purchaseOrderId: number;
      /** 付款账户ID */
      fundAccountId: number;
      /** 付款账户（现金 / 银行） */
      fundAccountName: string;
      /** 自动生成的凭证号，空字符串表示未记账 */
      voucherNo: string;
      /** 备注 */
      remark: string;
      /** 创建时间 */
      createdAt: string;
    }

    /** 付款记录详情 */
    interface Detail extends Item {
      /** 所属项目合同金额 */
      projectAmount: number;
      /** 所属项目累计已付款 */
      projectPaid: number;
      /** 自动生成的凭证ID，0 表示未记账 */
      voucherId: number;
      /** 凭证借贷分录 */
      ledgerEntries: LedgerEntry[];
    }

    /** 付款记录列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：项目名称 / 客户姓名 / 收款方 / 费用类型 / 备注 */
      keyword?: string | null;
      /** 所属项目 */
      projectId?: number | null;
      /** 费用类型 */
      payoutType?: string | null;
      /** 付款来源：purchase=采购付款 manual=其他付款 */
      source?: 'purchase' | 'manual' | null;
    }

    /** 付款记录分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑付款记录提交的数据 */
    interface FormData {
      /** 所属项目ID，null / 0 表示公司公共支出 */
      projectId: number | null;
      /** 付款金额 */
      amount: number | null;
      /** 费用类型 */
      payoutType: string;
      /** 付款日期（YYYY-MM-DD） */
      payoutDate: string | null;
      /** 供应商ID */
      vendorId: number | null;
      /** 关联采购单ID，手工付款不传 */
      purchaseOrderId?: number | null;
      /** 付款账户ID（财务已初始化时必填） */
      fundAccountId: number | null;
      /** 备注 */
      remark: string;
    }

    /** 付款概览统计 */
    interface Stats {
      /** 付款笔数 */
      totalCount: number;
      /** 付款总额 */
      totalAmount: number;
      /** 本月付款金额 */
      monthAmount: number;
      /** 涉及项目数 */
      projectCount: number;
    }
  }
}
