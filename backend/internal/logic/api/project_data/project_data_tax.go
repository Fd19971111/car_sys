package project_data

import "github.com/gogf/gf/v2/errors/gerror"

const (
	taxModeNoInvoice = "NO_INVOICE"
	taxModeIncluded  = "TAX_INCLUDED"
	taxModeExcluded  = "TAX_EXCLUDED"
)

// PurchaseTaxAmounts separates the pre-tax amount, tax amount, and payable total.
type PurchaseTaxAmounts struct {
	NetAmount   float64
	TaxAmount   float64
	TotalAmount float64
}

// validateTaxSelection checks the independently selected tax mode and rate.
func ValidateTaxSelection(mode string, rate float64) error {
	switch mode {
	case taxModeNoInvoice:
		if rate != 0 {
			return gerror.New("无票方式的税率必须为 0")
		}
	case taxModeIncluded, taxModeExcluded:
		if rate < 0 || rate > 100 {
			return gerror.New("税率必须在 0% 到 100% 之间")
		}
	default:
		return gerror.New("请选择有效的税务方式")
	}
	return nil
}

// calculatePurchaseTax calculates tax amounts without mixing tax rules into procurement flow.
// UnitPrice is gross for TAX_INCLUDED and net for TAX_EXCLUDED / NO_INVOICE.
func calculatePurchaseTax(quantity, unitPrice float64, mode string, rate float64) (*PurchaseTaxAmounts, error) {
	if quantity < 0 || unitPrice < 0 {
		return nil, gerror.New("采购数量和单价不能为负数")
	}
	if err := ValidateTaxSelection(mode, rate); err != nil {
		return nil, err
	}

	baseAmount := round2(quantity * unitPrice)
	amounts := &PurchaseTaxAmounts{NetAmount: baseAmount, TotalAmount: baseAmount}
	switch mode {
	case taxModeIncluded:
		if rate > 0 {
			amounts.NetAmount = round2(baseAmount / (1 + rate/100))
		}
		amounts.TaxAmount = round2(baseAmount - amounts.NetAmount)
	case taxModeExcluded:
		amounts.TaxAmount = round2(baseAmount * rate / 100)
		amounts.TotalAmount = round2(baseAmount + amounts.TaxAmount)
	}
	return amounts, nil
}
