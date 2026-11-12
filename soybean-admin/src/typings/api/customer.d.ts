declare namespace Api {
  /**
   * namespace Customer
   *
   * backend api module: "customer"（客户管理）
   */
  namespace Customer {
    /** 客户状态：1=跟进中 2=已成交 3=已流失 */
    type Status = 1 | 2 | 3;

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
