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
      /** 合同总额 */
      contractAmount: number;
      /** 累计回款 */
      receivedAmount: number;
      /** 所选年份回款 */
      yearReceived: number;
      /** 待回款 */
      unpaidAmount: number;
      /** 回款率（%） */
      receivedRate: number;
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

    /** 项目成本构成项 */
    interface ProjectProfitCostItem {
      /** 成本类型，如：材料款 / 人工费 */
      type: string;
      amount: number;
      /** 占成本合计百分比 */
      percent: number;
    }

    /** 项目毛利行 */
    interface ProjectProfitItem {
      projectId: number;
      projectName: string;
      customerName: string;
      /** 项目状态：1=待开工 2=施工中 3=已完工 4=已结算 */
      status: number;
      startDate: string;
      endDate: string;
      /** 合同金额 */
      contractAmount: number;
      /** 变更净额（增项-减项） */
      changeNet: number;
      /** 结算金额 = 合同 + 变更净额 */
      settlementAmount: number;
      /** 已收款 */
      receivedAmount: number;
      /** 待收款 */
      unpaidAmount: number;
      /** 回款率（%） */
      receivedRate: number;
      /** 成本合计（按付款归集） */
      costTotal: number;
      /** 成本构成 */
      costItems: ProjectProfitCostItem[];
      /** 材料款付款 */
      materialPaid: number;
      /** 材料耗用（领用金额，参考） */
      materialUsed: number;
      /** 预计毛利 = 结算金额 - 成本合计 */
      grossProfit: number;
      /** 毛利率（%） */
      grossRate: number;
      /** 现金结余 = 已收款 - 成本合计 */
      cashProfit: number;
    }

    /** 项目毛利汇总 */
    interface ProjectProfitSummary {
      projectCount: number;
      settlementTotal: number;
      receivedTotal: number;
      unpaidTotal: number;
      costTotal: number;
      grossProfit: number;
      /** 整体毛利率（%） */
      grossRate: number;
      /** 亏损项目数 */
      lossCount: number;
      /** 公司公共支出（未关联项目，未摊入成本） */
      companyCost: number;
      /** 材料耗用合计（参考） */
      materialUsedTotal: number;
    }

    /** 项目毛利列表 */
    interface ProjectProfitList {
      records: ProjectProfitItem[];
      summary: ProjectProfitSummary;
    }

    /** 项目毛利查询参数 */
    interface ProjectProfitSearchParams {
      keyword?: string | null;
      status?: number | null;
      onlyLoss?: boolean | null;
    }

    /** 项目收款明细 */
    interface ProjectProfitReceipt {
      id: number;
      paymentType: string;
      amount: number;
      paymentDate: string;
      remark: string;
    }

    /** 项目付款明细（成本） */
    interface ProjectProfitPayout {
      id: number;
      payoutType: string;
      amount: number;
      payoutDate: string;
      vendorName: string;
      remark: string;
    }

    /** 项目材料领用明细（参考） */
    interface ProjectProfitMaterial {
      id: number;
      materialName: string;
      unit: string;
      quantity: number;
      unitPrice: number;
      amount: number;
      usedAt: string;
      receiver: string;
    }

    /** 项目变更明细 */
    interface ProjectProfitChange {
      id: number;
      changeType: 'ADD' | 'REDUCE';
      title: string;
      amount: number;
      changeDate: string;
      reason: string;
    }

    /** 项目毛利明细 */
    interface ProjectProfitDetail extends ProjectProfitItem {
      receipts: ProjectProfitReceipt[];
      payouts: ProjectProfitPayout[];
      materials: ProjectProfitMaterial[];
      changes: ProjectProfitChange[];
    }
  }
}
