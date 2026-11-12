declare namespace Api {
  /**
   * namespace Order
   *
   * backend api module: "order"（订单中心：全部订单 / 行程监控）
   */
  namespace Order {
    /** 订单状态 */
    type Status =
      | 'quoting'
      | 'awaiting-confirmation'
      | 'dispatching'
      | 'driver-en-route'
      | 'in-progress'
      | 'completed'
      | 'cancelled';

    /** 客户反馈 */
    type CustomerResponse = 'pending' | 'accepted' | 'rejected';

    /** 车型 */
    type VehicleCategory = 'economy' | 'comfort' | 'business' | 'six-seat';

    /** 订单 */
    interface Item {
      id: number;
      orderNumber: string;
      customerId: number;
      customerName: string;
      customerPhone: string;
      status: Status;
      passengerCount: number;
      vehicleCategory: VehicleCategory;
      pickupAddress: string;
      destinationAddress: string;
      estimatedDistanceKm: number;
      estimatedDurationMin: number;
      quotedFare: number;
      finalFare: number;
      driverFare: number;
      driverFare: number;
      customerResponse: CustomerResponse;
      scheduledAt: string;
      driverId: number;
      driverName: string;
      driverPhone: string;
      vehicleId: number;
      vehiclePlate: string;
      completedAt: string;
      internalNotes: string;
      createdAt: string;
      updatedAt: string;
    }

    /** 订单查询参数 */
    interface SearchParams extends Api.Common.CommonSearchParams {
      keyword?: string | null;
      status?: Status | null;
      customerResponse?: CustomerResponse | null;
      dateFrom?: string | null;
      dateTo?: string | null;
    }

    /** 订单分页列表 */
    type List = Api.Common.PaginatingQueryRecord<Item>;

    /** 新增 / 编辑订单提交的数据 */
    interface FormData {
      customerId: number;
      status: Status;
      passengerCount: number;
      vehicleCategory: VehicleCategory;
      pickupAddress: string;
      destinationAddress: string;
      estimatedDistanceKm: number;
      estimatedDurationMin: number;
      quotedFare: number;
      finalFare: number;
      driverFare: number;
      customerResponse: CustomerResponse;
      scheduledAt: string;
      driverId: number;
      vehicleId: number;
      completedAt: string | null;
      internalNotes: string;
    }

    /** 状态数量 */
    interface StatusCount {
      status: string;
      count: number;
    }

    /** 订单统计 */
    interface Stats {
      total: number;
      today: number;
      active: number;
      completed: number;
      cancelled: number;
      revenue: number;
      statusCounts: StatusCount[];
    }

    /** 行程监控 */
    interface Monitor {
      activeCount: number;
      dispatching: number;
      driverEnRoute: number;
      inProgress: number;
      records: Item[];
      updatedAt: string;
    }
  }
}
