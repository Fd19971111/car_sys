declare namespace Api {
  /**
   * namespace Project
   *
   * backend api module: "project_list"（对应后端 project 表）
   */
  namespace Project {
    /** 项目状态：1=待开工 2=施工中 3=已完工 4=已结算 */
    type Status = 1 | 2 | 3 | 4;

    /** 下拉选项 */
    interface Option {
      /** 选项值（主键ID） */
      value: number;
      /** 选项文案 */
      label: string;
    }

    /** 项目信息（列表 / 详情共用） */
    interface Item {
      id: number;
      /** 所属客户ID */
      customerId: number;
      /** 所属客户姓名 */
      customerName: string;
      /** 项目名称 */
      name: string;
      /** 施工地址 */
      address: string;
      /** 负责人ID（site_user.id） */
      managerId: number;
      /** 负责人姓名 */
      managerName: string;
      /** 项目状态：1=待开工 2=施工中 3=已完工 4=已结算 */
      status: Status;
      /** 合同金额 */
      amount: number;
      /** 开工日期 */
      startDate: string | null;
      /** 预计完工日期 */
      endDate: string | null;
      /** 备注 */
      remark: string;
      /** 最近一条进度记录（列表接口返回） */
      latestProgress?: string;
      /** 最近进度记录时间（列表接口返回） */
      latestProgressAt?: string;
      /** 创建时间 */
      createdAt: string;
      /** 更新时间 */
      updatedAt: string;
    }

    /** 施工进度记录（详情时间轴） */
    interface ProgressLog {
      id: number;
      /** 进度内容 */
      content: string;
      /** 进度图片URL集合 */
      images: string[];
      /** 记录人姓名 */
      userName: string;
      /** 进度业务日期（YYYY-MM-DD，进度实际发生日期） */
      logDate?: string;
      /** 记录时间 */
      createdAt: string;
    }

    /** 项目详情 */
    interface Detail extends Item {
      /** 客户手机号 */
      customerPhone: string;
      /** 客户房屋地址 */
      customerAddress: string;
      /** 施工进度记录（按时间倒序） */
      logs: ProgressLog[];
    }

    /** 项目列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：匹配项目名称 / 施工地址 / 客户姓名 */
      keyword?: string | null;
      /** 项目状态 */
      status?: Status | null;
      /** 所属客户 */
      customerId?: number | null;
    }

    /** 项目分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑项目提交的数据 */
    interface FormData {
      /** 所属客户ID */
      customerId: number | null;
      /** 项目名称 */
      name: string;
      /** 施工地址 */
      address: string;
      /** 负责人ID */
      managerId: number | null;
      /** 项目状态 */
      status: Status;
      /** 合同金额 */
      amount: number | null;
      /** 开工日期（YYYY-MM-DD） */
      startDate: string | null;
      /** 预计完工日期（YYYY-MM-DD） */
      endDate: string | null;
      /** 备注 */
      remark: string;
    }

    /** 项目概览统计 */
    interface Stats {
      /** 项目总数 */
      total: number;
      /** 待开工数量 */
      pending: number;
      /** 施工中数量 */
      building: number;
      /** 已完工数量 */
      finished: number;
      /** 已结算数量 */
      settled: number;
      /** 合同总金额 */
      totalAmount: number;
    }

    /** 项目表单下拉选项（客户 / 负责人） */
    interface Options {
      customers: Option[];
      managers: Option[];
    }

    /**
     * 项目下拉选项
     *
     * - 供报价 / 收款等模块选择所属项目
     * - label 由后端拼装为「项目名称 · 客户姓名」
     */
    interface SelectOption {
      value: number;
      label: string;
      /** 合同金额（表单联动展示用） */
      amount: number;
      /** 项目状态 */
      status: Status;
    }
  }
}
