declare namespace Api {
  namespace User {
    interface Item {
      id: number;
      username: string;
      nickname: string;
      phone: string;
      role: number;
      roleName: string;
      status: number;
      departmentId: number;
      departmentName: string;
      lastLoginAt: string;
      createdAt: string;
    }

    interface SearchParams extends Api.Common.CommonSearchParams {
      keyword?: string | null;
      status?: number | null;
      role?: number | null;
    }

    type List = Api.Common.PaginatingQueryRecord<Item>;

    interface CreateAccountForm {
      staffId: number;
      username: string;
      password: string;
    }

    interface StaffOption {
      id: number;
      name: string;
      phone: string;
    }

    interface UpdateAccountForm {
      id: number;
      username: string;
      role: number;
    }
  }
}