package types

import (
	"fmt"
)

type Position struct {
	Ticker string  `json:"ticker"`
	Buy    float64 `json:"buy"`
	Loss   float64 `json:"loss"`
	Profit float64 `json:"profit"`
}

func (p Position) String() string {
	if p.Buy == 0 {
		return fmt.Sprintf("\nBuy Position: %s\n", p.Ticker)
	}

	str := fmt.Sprintf("\nPosition for %s:\n", p.Ticker)
	str += fmt.Sprintf("  - Buy price: %.2f\n", p.Buy)

	if p.Loss != 0 {
		str += fmt.Sprintf("  - Loss sell price: %.2f\n", p.Loss)
	} else {
		str += "  - Loss sell price: Not set\n"
	}

	if p.Profit != 0 {
		str += fmt.Sprintf("  - Profit sell price: %.2f\n", p.Profit)
	} else {
		str += "  - Profit sell price: Not set\n"
	}

	return str
}
