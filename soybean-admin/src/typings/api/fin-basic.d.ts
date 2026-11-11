declare namespace Api {
  /**
   * namespace FinBasic
   *
   * backend api module: "fin_basic"（财务基础资料：会计科目 / 资金账户）
   */
  namespace FinBasic {
    /** 科目类别 */
    type CategoryCode = 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'EXPENSE';

    /** 余额方向 */
    type DirectionCode = 'DEBIT' | 'CREDIT';

    /** 会计科目 */
    interface Account {
      id: number;
      /** 科目编码 */
      code: string;
      /** 科目名称 */
      name: string;
      /** 科目类别 */
      categoryCode: CategoryCode;
      /** 余额方向，由类别自动决定 */
      directionCode: DirectionCode;
      /** 是否启用 */
      enabled: boolean;
      /** 是否系统预置 */
      isSystem: boolean;
      /** 排序号 */
      sortOrder: number;
      /** 是否已被凭证、余额或资金账户引用 */
      used: boolean;
    }

    /** 会计科目统计 */
    interface AccountSummary {
      total: number;
      enabled: number;
      disabled: number;
      system: number;
      fundAccountCount: number;
    }

    /** 会计科目列表结果 */
    interface AccountListResult {
      records: Account[];
      summary: AccountSummary;
    }

    /** 会计科目查询参数 */
    interface AccountSearchParams {
      keyword?: string | null;
      categoryCode?: CategoryCode | null;
      enabled?: boolean | null;
    }

    /** 会计科目保存参数 */
    interface AccountFormData {
      id?: number;
      code: string;
      name: string;
      categoryCode: CategoryCode;
      enabled: boolean;
    }

    /** 资金账户类型 */
    type FundAccountType = 'CASH' | 'BANK';

    /** 资金账户 */
    interface FundAccount {
      id: number;
      type: FundAccountType;
      /** 账户显示名称 */
      name: string;
      /** 开户银行 */
      bankName: string;
      /** 卡号末四位 */
      accountLast4: string;
      /** 关联会计科目ID */
      ledgerAccountId: number;
      /** 关联会计科目编码 */
      ledgerAccountCode: string;
      /** 关联会计科目名称 */
      ledgerAccountName: string;
      enabled: boolean;
    }

    /** 资金账户列表结果 */
    interface FundAccountListResult {
      records: FundAccount[];
    }

    /** 资金账户查询参数 */
    interface FundAccountSearchParams {
      keyword?: string | null;
      type?: FundAccountType | null;
    }

    /** 资金账户保存参数 */
    interface FundAccountFormData {
      id?: number;
      type: FundAccountType;
      name: string;
      bankName: string;
      accountLast4: string;
      ledgerAccountCode: string;
      enabled: boolean;
    }
  }
}
