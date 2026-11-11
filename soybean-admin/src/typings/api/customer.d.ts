declare namespace Api {
  /**
   * namespace Customer
   *
   * backend api module: "customer"（对应后端 customer 表 / project 表）
   */
  namespace Customer {
    /** 客户状态：1=跟进中 2=已成交 3=已流失 */
    type Status = 1 | 2 | 3;

    /** 装修项目状态：1=待开工 2=施工中 3=已完工 4=已结算 */
    type ProjectStatus = 1 | 2 | 3 | 4;

    /** 客户名下的装修项目（详情展示用） */
    interface RelatedProject {
      id: number;
      /** 项目名称 */
      name: string;
      /** 项目状态：1=待开工 2=施工中 3=已完工 4=已结算 */
      status: ProjectStatus;
      /** 合同金额 */
      amount: number;
      /** 项目负责人 */
      managerName: string;
      /** 开工日期 */
      startDate: string | null;
      /** 预计完工日期 */
      endDate: string | null;
      /** 最近一条进度记录 */
      latestProgress?: string;
      /** 最近进度记录时间 */
      latestProgressAt?: string;
    }

    /** 客户信息 */
    interface Item {
      id: number;
      /** 客户姓名 */
      name: string;
      /** 手机号 */
      phone: string;
      /** 房屋地址 */
      address: string;
      /** 客户来源：转介绍/线上/门店等 */
      source: string;
      /** 客户状态：1=跟进中 2=已成交 3=已流失 */
      status: Status;
      /** 备注 */
      remark: string;
      /** 创建时间 */
      createdAt: string;
      /** 更新时间 */
      updatedAt: string;
      /**
       * 关联装修项目数量
       *
       * - 列表接口返回（后端按 project 子表实时统计）
       * - 详情接口可以不返回
       */
      projectCount?: number;
      /**
       * 关联的装修项目
       *
       * - 列表接口不返回；详情接口返回
       * - 对应后端 project 子表（customer_id 关联），进度来自 project_log 子表
       */
      projects?: RelatedProject[];
    }

    /** 客户列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：匹配客户姓名 / 手机号 / 房屋地址 */
      keyword?: string | null;
      /** 客户状态 */
      status?: Status | null;
      /** 客户来源 */
      source?: string | null;
    }

    /** 客户分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑客户提交的数据 */
    interface FormData {
      name: string;
      phone: string;
      address: string;
      source: string;
      status: Status;
      remark: string;
    }

    /** 客户概览统计 */
    interface Stats {
      /** 客户总数 */
      total: number;
      /** 跟进中数量 */
      following: number;
      /** 已成交数量 */
      deal: number;
      /** 已流失数量 */
      lost: number;
      /** 已成交客户的合同总金额 */
      dealAmount: number;
    }
  }
}
