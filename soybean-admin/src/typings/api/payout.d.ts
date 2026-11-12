declare namespace Api {
  /**
   * namespace Payout
   *
   * backend api module: "payout"（付款管理）
   */
  namespace Payout {
    /** 凭证分录：详情页展示「借什么、贷什么」 */
    interface LedgerEntry {
      direction: 'DEBIT' | 'CREDIT';
      accountName: string;
      amount: number;
      summary: string;
    }

    /** 付款记录（列表 / 详情共用） */
    interface Item {
      id: number;
      amount: number;
      /** 费用类型：材料款 / 人工费 / 施工费 / 设计费 / 管理费 / 销售费 / 其他 */
      payoutType: string;
      payoutDate: string;
      vendorName: string;
      fundAccountId: number;
      fundAccountName: string;
      voucherNo: string;
      remark: string;
      createdAt: string;
    }

    /** 付款记录详情 */
    interface Detail extends Item {
      voucherId: number;
      ledgerEntries: LedgerEntry[];
    }

    /** 付款记录列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：收款方 / 费用类型 / 备注 */
      keyword?: string | null;
      payoutType?: string | null;
    }

    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑付款记录提交的数据 */
    interface FormData {
      amount: number | null;
      payoutType: string;
      payoutDate: string | null;
      vendorName: string;
      fundAccountId: number | null;
      remark: string;
    }

    /** 付款概览统计 */
    interface Stats {
      totalCount: number;
      totalAmount: number;
      monthAmount: number;
      monthCount: number;
    }
  }
}
