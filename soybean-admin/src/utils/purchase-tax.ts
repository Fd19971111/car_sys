export type PurchaseTaxMode = 'NO_INVOICE' | 'TAX_INCLUDED' | 'TAX_EXCLUDED';

export interface PurchaseTaxAmounts {
  netAmount: number;
  taxAmount: number;
  totalAmount: number;
}

/** Unit price is gross for TAX_INCLUDED and net for TAX_EXCLUDED / NO_INVOICE. */
export function calculatePurchaseTax(
  quantity: number,
  unitPrice: number,
  mode: PurchaseTaxMode,
  rate: number
): PurchaseTaxAmounts {
  const baseAmount = Number((Math.max(0, quantity) * Math.max(0, unitPrice)).toFixed(2));
  const amounts: PurchaseTaxAmounts = { netAmount: baseAmount, taxAmount: 0, totalAmount: baseAmount };

  if (mode === 'TAX_INCLUDED') {
    amounts.netAmount = rate > 0 ? Number((baseAmount / (1 + rate / 100)).toFixed(2)) : baseAmount;
    amounts.taxAmount = Number((baseAmount - amounts.netAmount).toFixed(2));
  } else if (mode === 'TAX_EXCLUDED') {
    amounts.taxAmount = Number((baseAmount * rate / 100).toFixed(2));
    amounts.totalAmount = Number((baseAmount + amounts.taxAmount).toFixed(2));
  }

  return amounts;
}
