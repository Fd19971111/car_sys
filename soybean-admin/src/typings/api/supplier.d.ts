declare namespace Api {
  namespace Supplier {
    interface Item {
      id: number;
      name: string;
      contactName: string;
      phone: string;
      address: string;
      remark: string;
      createdAt: string;
      updatedAt: string;
    }

    interface SearchParams extends Api.Common.CommonSearchParams {
      keyword?: string | null;
    }

    type List = Api.Common.PaginatingQueryRecord<Item>;

    interface FormData {
      name: string;
      contactName: string;
      phone: string;
      address: string;
      remark: string;
    }
  }
}