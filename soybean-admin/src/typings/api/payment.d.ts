declare namespace Api {
  /**
   * namespace Payment
   *
   * backend api module: "payment"（收款管理）
   */
  namespace Payment {
    /** 凭证分录：详情页展示「借什么、贷什么」 */
    interface LedgerEntry {
      direction: 'DEBIT' | 'CREDIT';
      accountName: string;
      amount: number;
      summary: string;
    }

    /** 收款记录（列表 / 详情共用） */
    interface Item {
      id: number;
      rideOrderId: number;
      amount: number;
      /** 款项类型：定金 / 首期款 / 中期款 / 尾款 */
      paymentType: string;
      paymentDate: string;
      fundAccountId: number;
      fundAccountName: string;
      voucherNo: string;
      remark: string;
      createdAt: string;
    }

    /** 收款记录详情 */
    interface Detail extends Item {
      voucherId: number;
      ledgerEntries: LedgerEntry[];
    }

    /** 收款记录列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：款项类型 / 备注 */
      keyword?: string | null;
      paymentType?: string | null;
    }

    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑收款记录提交的数据 */
    interface FormData {
      rideOrderId: number | null;
      amount: number | null;
      paymentType: string;
      paymentDate: string | null;
      fundAccountId: number | null;
      remark: string;
    }

    /** 收款概览统计 */
    interface Stats {
      totalCount: number;
      totalAmount: number;
      monthAmount: number;
      monthCount: number;
    }
  }
}
