declare namespace Api {
  /**
   * namespace Department
   *
   * backend api module: "department"（企业部门管理）
   */
  namespace Department {
    /** 部门状态：1=启用 2=停用 */
    type Status = 1 | 2;

    /** 部门信息 */
    interface Item {
      id: number;
      /** 部门名称 */
      name: string;
      /** 部门编码 */
      code: string;
      /** 负责人姓名 */
      leaderName: string;
      /** 负责人ID */
      leaderId: number | null;
      /** 成员数量 */
      memberCount: number;
      /** 状态：1=启用 2=停用 */
      status: Status;
      /** 备注 */
      remark: string;
      /** 排序 */
      sortOrder: number;
      /** 创建时间 */
      createdAt: string;
      /** 更新时间 */
      updatedAt: string;
    }

    /** 部门列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：匹配部门名称 / 编码 */
      keyword?: string | null;
      /** 部门状态 */
      status?: Status | null;
    }

    /** 部门分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑部门提交的数据 */
    interface FormData {
      name: string;
      code: string;
      leaderId: number | null;
      status: Status;
      sortOrder: number;
      remark: string;
    }
  }
}