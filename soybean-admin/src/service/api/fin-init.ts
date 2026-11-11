import { request } from '../request';

export interface FinanceInitAccount {
  code: string;
  name: string;
  categoryCode: 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'EXPENSE';
  enabled: boolean;
}

export interface FinanceInitFundAccount {
  type: 'CASH' | 'BANK';
  name: string;
  bankName: string;
  accountLast4: string;
  ledgerAccountCode: string;
}

export interface FinanceInitOpeningBalance {
  accountCode: string;
  amount: string;
}

export interface FinanceInitData {
  year: number;
  startMonth: number;
  accounts: FinanceInitAccount[];
  fundAccounts: FinanceInitFundAccount[];
  openingBalances: FinanceInitOpeningBalance[];
  reconfigure: boolean;
}

export function fetchSaveFinanceInit(data: FinanceInitData) {
  return request<{ initialized: boolean }>({
    url: '/fin-init/save',
    method: 'post',
    data
  });
}