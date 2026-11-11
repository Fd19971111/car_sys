package project_data

import "testing"

func TestCalculatePurchaseTax(t *testing.T) {
	tests := []struct {
		name      string
		quantity  float64
		unitPrice float64
		mode      string
		rate      float64
		want      PurchaseTaxAmounts
	}{
		{name: "no invoice", quantity: 3, unitPrice: 12.5, mode: taxModeNoInvoice, want: PurchaseTaxAmounts{NetAmount: 37.5, TotalAmount: 37.5}},
		{name: "tax included", quantity: 2, unitPrice: 109, mode: taxModeIncluded, rate: 9, want: PurchaseTaxAmounts{NetAmount: 200, TaxAmount: 18, TotalAmount: 218}},
		{name: "tax excluded", quantity: 2, unitPrice: 100, mode: taxModeExcluded, rate: 13, want: PurchaseTaxAmounts{NetAmount: 200, TaxAmount: 26, TotalAmount: 226}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := calculatePurchaseTax(test.quantity, test.unitPrice, test.mode, test.rate)
			if err != nil {
				t.Fatalf("calculatePurchaseTax() error = %v", err)
			}
			if *got != test.want {
				t.Fatalf("calculatePurchaseTax() = %+v, want %+v", *got, test.want)
			}
		})
	}
}

func TestValidateTaxSelection(t *testing.T) {
	if err := ValidateTaxSelection(taxModeNoInvoice, 1); err == nil {
		t.Fatal("expected no-invoice mode with non-zero tax rate to fail")
	}
	if err := ValidateTaxSelection(taxModeIncluded, 100.01); err == nil {
		t.Fatal("expected tax rate above 100% to fail")
	}
}
