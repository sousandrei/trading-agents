package types

import (
	"fmt"
	"regexp"
	"strings"
)

type Action struct {
	Ticker string `json:"ticker"`
	Action string `json:"action"`
	Price  string `json:"price,omitempty"`
	Loss   string `json:"loss,omitempty"`
	Profit string `json:"profit,omitempty"`
}

func (a Action) String() string {
	switch a.Action {
	case "BUY":
		return fmt.Sprintf("Action: %s, Price: %s", a.Action, a.Price)
	case "SELL":
		return fmt.Sprintf("Action: %s, Price: %s", a.Action, a.Price)
	case "UPDATE_STOPS":
		return fmt.Sprintf("Action: %s, Loss: %s, Profit: %s", a.Action, a.Loss, a.Profit)
	case "HOLD":
		return fmt.Sprintf("Action: %s", a.Action)
	default:
		return fmt.Sprintf("Unknown action: %s", a.Action)
	}
}

func priceRegexp(label string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`%s:\s([0-9]{0,4}.[0-9]{0,2})`, label))
}

func ParseOutput(ticker, output string) (*Action, error) {
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
		priceMatches := priceRegexp("BUY PRICE").FindStringSubmatch(parts[1])
		price := strings.TrimSpace(priceMatches[1])

		return &Action{
			Ticker: ticker,
			Action: "BUY",
			Price:  price,
		}, nil

	case "SELL":
		priceMatches := priceRegexp("SELL PRICE").FindStringSubmatch(parts[1])
		price := strings.TrimSpace(priceMatches[1])

		return &Action{
			Ticker: ticker,
			Action: "SELL",
			Price:  price,
		}, nil

	case "UPDATE_STOPS":
		lossStr := priceRegexp("LOSS").FindStringSubmatch(parts[1])
		loss := strings.TrimSpace(lossStr[1])

		profitStr := priceRegexp("PROFIT").FindStringSubmatch(parts[2])
		profit := strings.TrimSpace(profitStr[1])

		return &Action{
			Ticker: ticker,
			Action: "UPDATE_STOPS",
			Loss:   loss,
			Profit: profit,
		}, nil

	case "HOLD":
		return &Action{
			Ticker: ticker,
			Action: "HOLD",
		}, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}

}
