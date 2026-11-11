declare namespace Api {
  /**
   * namespace OpLog
   *
   * backend api module: "oplog"（数据报表 → 操作日志）
   */
  namespace OpLog {
    /** 操作动作 */
    type Action = 'CREATE' | 'UPDATE' | 'DELETE';

    /** 操作日志行 */
    interface Item {
      /** 日志ID */
      id: number;
      /** 操作时间 YYYY-MM-DD HH:mm:ss */
      createdAt: string;
      /** 操作人ID */
      userId: number;
      /** 操作人 */
      userName: string;
      /** 业务模块：收款 / 付款 / 项目 / 库存 */
      module: string;
      /** 动作 */
      action: Action;
      /** 动作中文名 */
      actionName: string;
      /** 对象类型，如 payment / payout */
      targetType: string;
      /** 对象ID */
      targetId: number;
      /** 操作摘要 */
      title: string;
      /** 涉及金额 */
      amount: number;
      /** 变更详情 JSON */
      detail: string;
      /** 操作来源IP */
      ip: string;
    }

    /** 操作日志查询参数 */
    interface SearchParams {
      /** 业务模块，空为全部 */
      module?: string | null;
      /** 动作，空为全部 */
      action?: Action | '' | null;
      /** 摘要 / 操作人 / 详情关键字 */
      keyword?: string;
      /** 开始日期 YYYY-MM-DD */
      startDate?: string | null;
      /** 结束日期 YYYY-MM-DD */
      endDate?: string | null;
      /** 页码 */
      current?: number;
      /** 每页条数 */
      size?: number;
    }

    /** 操作日志列表 */
    interface List {
      records: Item[];
      total: number;
      /** 已出现过的业务模块，供筛选下拉 */
      modules: string[];
    }
  }
}
