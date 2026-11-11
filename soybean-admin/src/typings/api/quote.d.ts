declare namespace Api {
  /**
   * namespace Quote
   *
   * backend api module: "quote"（对应后端 project_quote 表）
   */
  namespace Quote {
    /** 报价明细（列表 / 详情共用） */
    interface Item {
      id: number;
      /** 所属项目ID */
      projectId: number;
      /** 项目名称 */
      projectName: string;
      /** 客户姓名 */
      customerName: string;
      /** 报价项名称，如：水电改造 */
      itemName: string;
      /** 计量单位，如：平米 */
      unit: string;
      /** 单价 */
      price: number;
      /** 数量 */
      quantity: number;
      /** 小计 = 单价 × 数量（后端计算） */
      total: number;
      /** 备注 */
      remark: string;
      /** 创建时间 */
      createdAt: string;
      /** 更新时间 */
      updatedAt: string;
    }

    /** 报价明细详情 */
    interface Detail extends Item {
      /** 所属项目合同金额 */
      projectAmount: number;
      /** 所属项目报价合计 */
      quoteAmount: number;
    }

    /** 报价明细列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：报价项名称 / 项目名称 / 客户姓名 */
      keyword?: string | null;
      /** 所属项目 */
      projectId?: number | null;
    }

    /** 报价明细分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑报价明细提交的数据 */
    interface FormData {
      /** 所属项目ID */
      projectId: number | null;
      /** 报价项名称 */
      itemName: string;
      /** 计量单位 */
      unit: string;
      /** 单价 */
      price: number | null;
      /** 数量 */
      quantity: number | null;
      /** 备注 */
      remark: string;
    }

    /** 报价概览统计 */
    interface Stats {
      /** 报价条目数 */
      totalCount: number;
      /** 报价总额 */
      totalAmount: number;
      /** 涉及项目数 */
      projectCount: number;
      /** 平均单条金额 */
      avgAmount: number;
    }
  }
}
