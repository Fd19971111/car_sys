declare namespace Api {
  /**
   * namespace FinVoucher
   *
   * backend api module: "fin_voucher"（记账凭证 / 月结）
   */
  namespace FinVoucher {
    /** 凭证状态 */
    type VoucherStatus = 'DRAFT' | 'POSTED' | 'VOID';

    /** 凭证类型 */
    type VoucherType = 'GENERAL' | 'RECEIPT' | 'PAYMENT' | 'TRANSFER';

    /** 会计期间状态 */
    type PeriodStatus = 'OPEN' | 'CLOSED' | 'FUTURE';

    /** 凭证列表行 */
    interface Item {
      id: number;
      voucherNo: string;
      voucherDate: string;
      periodCode: string;
      voucherTypeCode: VoucherType;
      summary: string;
      statusCode: VoucherStatus;
      sourceType: string;
      sourceId: number;
      debitTotal: number;
      creditTotal: number;
      createdAt: string;
    }

    /** 凭证分录 */
    interface Entry {
      lineNo: number;
      accountCode: string;
      accountName: string;
      summary: string;
      directionCode: 'DEBIT' | 'CREDIT';
      amount: number;
    }

    /** 凭证详情 */
    interface Detail extends Item {
      entries: Entry[];
    }

    /** 凭证列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      periodCode?: string | null;
      statusCode?: VoucherStatus | null;
      keyword?: string | null;
    }

    /** 凭证分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 会计期间 */
    interface Period {
      periodCode: string;
      periodYear: number;
      periodMonth: number;
      statusCode: PeriodStatus;
      voucherCount: number;
      debitTotal: number;
      creditTotal: number;
      canClose: boolean;
      canReopen: boolean;
      beforeStart: boolean;
    }

    /** 会计期间列表 */
    interface PeriodList {
      initialized: boolean;
      currentPeriod: string;
      startPeriod: string;
      records: Period[];
    }

    /** 凭证导出条件 */
    interface ExportParams {
      periodCode?: string | null;
      statusCode?: VoucherStatus | null;
      keyword?: string | null;
    }

    /** 凭证导出行（一条分录一行） */
    interface ExportRow {
      voucherNo: string;
      voucherDate: string;
      periodCode: string;
      voucherTypeCode: VoucherType;
      voucherSummary: string;
      statusCode: VoucherStatus;
      sourceType: string;
      sourceId: number;
      lineNo: number;
      accountCode: string;
      accountName: string;
      entrySummary: string;
      directionCode: 'DEBIT' | 'CREDIT';
      amount: number;
    }

    /** 凭证导出结果 */
    interface ExportResult {
      records: ExportRow[];
    }
  }
}
