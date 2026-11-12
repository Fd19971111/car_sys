declare namespace Api {
  namespace DriverFinance {
    interface DriverBalance {
      driverId: number;
      driverName: string;
      phone: string;
      tripCount: number;
      earnedTotal: number;
      settledTotal: number;
      payable: number;
    }

    interface Settlement {
      id: number;
      settlementNo: string;
      driverId: number;
      driverName: string;
      settlementDate: string;
      amount: number;
      fundAccountName: string;
      voucherNo: string;
      remark: string;
    }

    type DriverList = { records: DriverBalance[] };
    type SettlementList = { records: Settlement[]; total: number };
    type SettlementSearchParams = { driverId?: number; current: number; size: number };
    type SettlementForm = {
      driverId: number;
      amount: number;
      settlementDate: string;
      fundAccountId: number;
      remark: string;
    };
    type SettlementResult = { id: number; settlementNo: string; voucherNo: string };
  }
}