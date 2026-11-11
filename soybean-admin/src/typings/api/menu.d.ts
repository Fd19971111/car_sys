declare namespace Api {
  namespace Menu {
    interface Item {
      id: number;
      parentId: number;
      title: string;
      routeName: string;
      routePath: string;
      component: string;
      icon: string;
      sortOrder: number;
      enabled: boolean;
      hideInMenu: boolean;
    }

    interface List {
      items: Item[];
    }

    interface Department {
      id: number;
      name: string;
    }

    interface DepartmentList {
      items: Department[];
    }

    interface DepartmentMenuList {
      menuIds: number[];
    }

    interface DepartmentUser {
      id: number;
      username: string;
      nickname: string;
      role: number;
      departmentId: number;
    }

    interface UserList {
      items: DepartmentUser[];
    }

    type Form = Omit<Item, 'id'> & { id?: number };
  }
}