declare namespace Api {
  /**
   * namespace Staff
   *
   * backend api module: "staff"（企业人员管理）
   */
  namespace Staff {
    /** 人员状态：1=在职 2=离职 */
    type Status = 1 | 2;

    /** 人员信息 */
    interface Item {
      id: number;
      /** 姓名 */
      name: string;
      /** 手机号 */
      phone: string;
      /** 性别：男/女 */
      gender: string;
      /** 职位/岗位 */
      position: string;
      /** 所属部门ID */
      departmentId: number;
      /** 所属部门名称 */
      departmentName: string;
      /** 入职日期 */
      joinDate: string;
      /** 状态：1=在职 2=离职 */
      status: Status;
      /** 备注 */
      remark: string;
      /** 创建时间 */
      createdAt: string;
      /** 更新时间 */
      updatedAt: string;
    }

    /** 人员列表查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：匹配姓名 / 手机号 */
      keyword?: string | null;
      /** 部门ID */
      departmentId?: number | null;
      /** 人员状态 */
      status?: Status | null;
    }

    /** 人员分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑人员提交的数据 */
    interface FormData {
      name: string;
      phone: string;
      gender: string;
      position: string;
      departmentId: number;
      joinDate: string | null;
      status: Status;
      remark: string;
    }

    /** 部门下拉选项 */
    interface DepartmentOption {
      value: number;
      label: string;
    }
  }
}