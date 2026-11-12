declare namespace Api {
  namespace PricingRule {
    interface Item {
      id: number;
      /** 规则名称 */
      name: string;
      /** 适用车型 */
      carType: string;
      /** 起步价（元） */
      basePrice: number;
      /** 超里程单价（元/公里） */
      perKmPrice: number;
      /** 是否启用 */
      enabled: boolean;
      createdAt: string;
      updatedAt: string;
    }

    interface SearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：规则名称 / 适用车型 */
      keyword?: string | null;
    }

    type List = Api.Common.PaginatingQueryRecord<Item>;

    interface FormData {
      name: string;
      carType: string;
      basePrice: number;
      perKmPrice: number;
      enabled: boolean;
    }
  }
}
