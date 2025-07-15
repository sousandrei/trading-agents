package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Position struct {
	Ticker          string
	BuyPrice        float64
	LossSellPrice   float64
	ProfitSellPrice float64
}

func (p Position) String() string {
	if p.BuyPrice == 0 {
		return fmt.Sprintf("\nBuy Position to be considered: %s\n", p.Ticker)
	}

	str := fmt.Sprintf("\nPosition to be analyzed for %s:\n", p.Ticker)
	str += fmt.Sprintf("  - Buy price: %.2f\n", p.BuyPrice)

	if p.LossSellPrice != 0 {
		str += fmt.Sprintf("  - Loss sell price: %.2f\n", p.LossSellPrice)
	} else {
		str += "  - Loss sell price: Not set\n"
	}

	if p.ProfitSellPrice != 0 {
		str += fmt.Sprintf("  - Profit sell price: %.2f\n", p.ProfitSellPrice)
	} else {
		str += "  - Profit sell price: Not set\n"
	}

	return str
}

func ParsePositions(str string) ([]Position, error) {
	positions := []Position{}

	positionStrings := strings.Split(str, ";")
	for _, ps := range positionStrings {
		parts := strings.Split(ps, ":")

		if len(parts) < 1 {
			return nil, fmt.Errorf("invalid position format: %s", ps)
		}

		ticker := parts[0]

		if len(parts) < 2 {
			positions = append(positions, Position{
				Ticker: ticker,
			})
			continue
		}

		buyPrice, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid buy price for %s: %w", ticker, err)
		}

		if len(parts) < 3 {
			positions = append(positions, Position{
				Ticker:   ticker,
				BuyPrice: buyPrice,
				// default zero values for loss and profit sell prices
			})
			continue
		}

		lossSellPrice, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid loss sell price for %s: %w", ticker, err)
		}

		if len(parts) < 4 {
			positions = append(positions, Position{
				Ticker:        ticker,
				BuyPrice:      buyPrice,
				LossSellPrice: lossSellPrice,
				// default zero value for profit sell price
			})
			continue
		}

		profitSellPrice, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid profit sell price for %s: %w", ticker, err)
		}

		positions = append(positions, Position{
			Ticker:          ticker,
			BuyPrice:        buyPrice,
			LossSellPrice:   lossSellPrice,
			ProfitSellPrice: profitSellPrice,
		})

		if len(parts) > 4 {
			return nil, fmt.Errorf("too many fields for position %s: %s", ticker, ps)
		}
	}

	return positions, nil
}
