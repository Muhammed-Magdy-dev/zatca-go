package invoice

import "math"

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

type lineTotals struct {
	LineTotal    float64
	TaxAmount    float64
	TotalWithVAT float64
}

func calcLine(l InvoiceLine) lineTotals {
	baseTotal := round2(l.Quantity * l.Price)

	var discountTotal float64
	for _, d := range l.ItemDiscounts {
		if !d.Indicator {
			discountTotal += d.Amount
		}
	}

	lineTotal := round2(baseTotal - discountTotal)
	tax := round2(lineTotal * l.VATRate / 100)

	return lineTotals{
		LineTotal:    lineTotal,
		TaxAmount:    tax,
		TotalWithVAT: round2(lineTotal + tax),
	}
}

type InvoiceTotals struct {
	LineExtensionAmount float64
	TaxExclusiveAmount  float64
	TaxInclusiveAmount  float64
	AllowanceTotal      float64
	ChargeTotal         float64
	TaxAmount           float64
}

func CalculateTotals(input *InvoiceInput) InvoiceTotals {
	var lineExtension, taxAmount float64

	for _, l := range input.Lines {
		t := calcLine(l)
		lineExtension += t.LineTotal
		taxAmount += t.TaxAmount
	}

	var allowanceTotal, chargeTotal float64
	for _, ac := range input.InvoiceLevelACs {
		if ac.Indicator {
			chargeTotal += ac.Amount
		} else {
			allowanceTotal += ac.Amount
		}

		acTax := round2(ac.Amount * ac.VATRate / 100)

		if ac.Indicator {
			taxAmount += acTax
		} else {
			taxAmount -= acTax
		}
	}

	taxExclusive := round2(lineExtension - allowanceTotal + chargeTotal)
	taxAmount = round2(taxAmount)
	taxInclusive := round2(taxExclusive + taxAmount)

	return InvoiceTotals{
		LineExtensionAmount: round2(lineExtension),
		TaxExclusiveAmount:  taxExclusive,
		TaxInclusiveAmount:  taxInclusive,
		AllowanceTotal:      round2(allowanceTotal),
		ChargeTotal:         round2(chargeTotal),
		TaxAmount:           taxAmount,
	}
}

func taxCategoryID(vatRate float64) string {
	if vatRate > 0 {
		return "S"
	}
	return "O"
}
