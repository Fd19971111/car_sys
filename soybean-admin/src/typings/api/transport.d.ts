declare namespace Api {
  /**
   * namespace Transport
   *
   * backend api module: "transport"（运力管理：司机 / 车辆）
   */
  namespace Transport {
    /** 司机接单状态 */
    type DriverServiceStatus = 'available' | 'busy' | 'offline' | 'suspended';

    /** 司机档案 */
    interface Driver {
      id: number;
      /** 司机姓名 */
      name: string;
      /** 手机号 */
      phone: string;
      /** 驾驶证号 */
      licenseNumber: string;
      /** 准驾车型 */
      licenseType: string;
      /** 接单状态 */
      serviceStatus: DriverServiceStatus;
      /** 服务评分（0-5） */
      rating: number;
      /** 完成订单数 */
      completedTrips: number;
      /** 加入日期 YYYY-MM-DD */
      joinedAt: string;
      /** 备注 */
      remark: string;
      /** 创建时间 */
      createdAt: string;
      /** 更新时间 */
      updatedAt: string;
    }

    /** 司机列表查询参数 */
    interface DriverSearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：姓名 / 手机号 / 驾驶证号 */
      keyword?: string | null;
      /** 接单状态 */
      serviceStatus?: DriverServiceStatus | null;
    }

    /** 司机分页列表 */
    type DriverList = Api.Common.PaginatingQueryRecord<Driver>;

    /** 新增 / 编辑司机提交的数据 */
    interface DriverFormData {
      name: string;
      phone: string;
      licenseNumber: string;
      licenseType: string;
      serviceStatus: DriverServiceStatus;
      rating: number;
      completedTrips: number;
      joinedAt: string | null;
      remark: string;
    }

    /** 司机下拉选项 */
    interface DriverOption {
      id: number;
      name: string;
      phone: string;
      serviceStatus: DriverServiceStatus;
    }

    /** 车辆服务车型 */
    type VehicleCategory = 'economy' | 'comfort' | 'business' | 'six-seat';

    /** 车辆状态 */
    type VehicleStatus = 'active' | 'maintenance' | 'inactive';

    /** 车辆档案 */
    interface Vehicle {
      id: number;
      /** 车牌号 */
      plateNumber: string;
      /** 品牌车型 */
      brand: string;
      /** 服务车型 */
      category: VehicleCategory;
      /** 可乘人数 */
      seatCount: number;
      /** 车辆状态 */
      status: VehicleStatus;
      /** 当前司机ID */
      driverId: number;
      /** 当前司机姓名 */
      driverName: string;
      /** 备注 */
      remark: string;
      /** 创建时间 */
      createdAt: string;
      /** 更新时间 */
      updatedAt: string;
    }

    /** 车辆列表查询参数 */
    interface VehicleSearchParams extends Api.Common.CommonSearchParams {
      /** 关键词：车牌号 / 品牌 / 司机姓名 */
      keyword?: string | null;
      /** 车辆状态 */
      status?: VehicleStatus | null;
    }

    /** 车辆分页列表 */
    type VehicleList = Api.Common.PaginatingQueryRecord<Vehicle>;

    /** 新增 / 编辑车辆提交的数据 */
    interface VehicleFormData {
      plateNumber: string;
      brand: string;
      category: VehicleCategory;
      seatCount: number;
      status: VehicleStatus;
      driverId: number;
      remark: string;
    }
  }
}
