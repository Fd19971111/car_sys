declare namespace Api {
  /**
   * namespace FinReport
   *
   * backend api module: "fin_report"（财务报表：经营总览 / 资金流水 / 科目余额）
   */
  namespace FinReport {
    /** 资金方向 */
    type FlowDirection = 'INCOME' | 'EXPENSE';

    /** 科目类别 */
    type CategoryCode = 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'EXPENSE';

    /** 报表核心指标（金额单位：元） */
    interface Summary {
      /** 收入合计 */
      income: number;
      /** 支出合计 */
      expense: number;
      /** 结余 */
      net: number;
      /** 结余率（%） */
      netRate: number;
      /** 资金余额合计 */
      totalBalance: number;
      /** 现金余额 */
      cashBalance: number;
      /** 银行余额 */
      bankBalance: number;
      /** 应收账款余额 */
      receivable: number;
      /** 应付账款余额 */
      payable: number;
      /** 所选年份完单车费 */
      rideRevenue: number;
      /** 所选年份关联行程收款 */
      passengerPaid: number;
      /** 所选年份所有收款 */
      yearReceived: number;
      /** 完单待收车费余额 */
      unpaidFare: number;
      /** 司机待结余额 */
      driverPayable: number;
      /** 所选年份司机已结 */
      driverSettled: number;
      /** 所选年份司机应得 */
      driverExpense: number;
    }

    /** 月度收支 */
    interface MonthlyItem {
      month: number;
      label: string;
      income: number;
      expense: number;
      net: number;
    }

    /** 收入 / 支出按科目汇总 */
    interface AccountAmount {
      code: string;
      name: string;
      amount: number;
      /** 占同类合计百分比 */
      percent: number;
    }

    /** 资金账户余额 */
    interface FundBalance {
      id: number;
      type: 'CASH' | 'BANK';
      name: string;
      bankName: string;
      accountLast4: string;
      accountCode: string;
      accountName: string;
      balance: number;
    }

    /** 科目余额 */
    interface AccountBalance {
      code: string;
      name: string;
      categoryCode: CategoryCode;
      balance: number;
    }

    /** 资金流水 */
    interface FlowItem {
      id: number;
      /** 发生时间 YYYY-MM-DD HH:mm:ss */
      occurredAt: string;
      directionCode: FlowDirection;
      amount: number;
      remark: string;
      accountName: string;
      fundAccountName: string;
    }

    /** 报表总览 */
    interface Overview {
      /** 财务是否已初始化 */
      initialized: boolean;
      /** 当前会计期间 YYYY-MM */
      currentPeriod: string;
      /** 当前期间状态 OPEN / CLOSED / FUTURE */
      periodStatus: string;
      year: number;
      summary: Summary;
      monthly: MonthlyItem[];
      incomeItems: AccountAmount[];
      expenseItems: AccountAmount[];
      fundAccounts: FundBalance[];
      accountBalances: AccountBalance[];
      recentFlows: FlowItem[];
    }

    /** 总览查询参数 */
    interface OverviewParams {
      year?: number;
    }

    /** 资金流水查询参数 */
    interface FlowSearchParams {
      year?: number;
      /** 月份 1-12，null 表示全年 */
      month?: number | null;
      directionCode?: FlowDirection | null;
      fundAccountId?: number | null;
      keyword?: string | null;
      current?: number;
      size?: number;
    }

    /** 资金流水列表结果 */
    interface FlowListResult {
      records: FlowItem[];
      total: number;
      summary: {
        income: number;
        expense: number;
        net: number;
      };
    }

  }
}
