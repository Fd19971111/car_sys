declare namespace Api {
  /**
   * namespace Payment
   *
   * backend api module: "payment"（对应后端 project_payment 表，收款管理）
   */
  namespace Payment {
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

    /** 收款记录（列表 / 详情共用） */
    interface Item {
      id: number;
      /** 所属项目ID */
      projectId: number;
      /** 项目名称 */
      projectName: string;
      /** 客户姓名 */
      customerName: string;
      /** 收款金额 */
      amount: number;
      /** 款项类型：定金 / 首期款 / 中期款 / 尾款 */
      paymentType: string;
      /** 收款日期 */
      paymentDate: string;
      /** 收款账户ID */
      fundAccountId: number;
      /** 收款账户（现金 / 银行） */
      fundAccountName: string;
      /** 自动生成的凭证号，空字符串表示未记账 */
      voucherNo: string;
      /** 备注 */
      remark: string;
      /** 创建时间 */
      createdAt: string;
    }

    /** 收款记录详情 */
    interface Detail extends Item {
      /** 所属项目合同金额 */
      projectAmount: number;
      /** 所属项目累计已收款 */
      paidAmount: number;
      /** 自动生成的凭证ID，0 表示未记账 */
      voucherId: number;
      /** 凭证借贷分录 */
      ledgerEntries: LedgerEntry[];
    }

    /** 收款记录列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：项目名称 / 客户姓名 / 款项类型 / 备注 */
      keyword?: string | null;
      /** 所属项目 */
      projectId?: number | null;
      /** 款项类型 */
      paymentType?: string | null;
    }

    /** 收款记录分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑收款记录提交的数据 */
    interface FormData {
      /** 所属项目ID */
      projectId: number | null;
      /** 收款金额 */
      amount: number | null;
      /** 款项类型 */
      paymentType: string;
      /** 收款日期（YYYY-MM-DD） */
      paymentDate: string | null;
      /** 收款账户ID（财务已初始化时必填） */
      fundAccountId: number | null;
      /** 备注 */
      remark: string;
    }

    /** 收款概览统计 */
    interface Stats {
      /** 收款笔数 */
      totalCount: number;
      /** 收款总额 */
      totalAmount: number;
      /** 本月收款金额 */
      monthAmount: number;
      /** 涉及项目数 */
      projectCount: number;
    }
  }
}
