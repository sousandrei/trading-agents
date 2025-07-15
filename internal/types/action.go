package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Action struct {
	Action string
	Price  float64
	Loss   float64
	Profit float64
}

func priceRegexp(label string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`%s:\s([0-9]{0,4}.[0-9]{0,2})`, label))
}

func ParseOutput(output string) (*Action, error) {
	parts := strings.Split(output, "FINAL TRANSACTION PROPOSAL: ")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid output format: %s", output)
	}

	parts = strings.Split(parts[1], "\n")

	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid output format: %s", output)
	}

	action := strings.TrimSpace(parts[0])

	switch action {
	case "BUY":
		priceStr := priceRegexp("BUY PRICE").FindStringSubmatch(parts[1])

		price, err := strconv.ParseFloat(priceStr[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid BUY PRICE format: %s", priceStr[1])
		}

		return &Action{
			Action: "BUY",
			Price:  price,
		}, nil

	case "SELL":
		priceStr := priceRegexp("SELL PRICE").FindStringSubmatch(parts[1])

		price, err := strconv.ParseFloat(priceStr[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid SELL PRICE format: %s", priceStr[1])
		}

		return &Action{
			Action: "SELL",
			Price:  price,
		}, nil

	case "UPDATE_STOPS":
		lossStr := priceRegexp("LOSS").FindStringSubmatch(parts[1])

		loss, err := strconv.ParseFloat(lossStr[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid LOSS format: %s", lossStr[1])
		}

		profitStr := priceRegexp("PROFIT").FindStringSubmatch(parts[2])
		profit, err := strconv.ParseFloat(profitStr[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid PROFIT format: %s", profitStr[1])
		}

		return &Action{
			Action: "UPDATE_STOPS",
			Loss:   loss,
			Profit: profit,
		}, nil

	case "HOLD":
		return &Action{
			Action: "HOLD",
		}, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}

}
