package base

import (
	"github.com/okcoin-okex/open-api-v3-sdk/okex-go-sdk-api"
)

type Wallet struct {
	client    *okex.Client
	MCurrency map[string]okex.AccountWalletResult
}

func NewWallet(c *okex.Client) *Wallet {
	return &Wallet{
		client:    c,
		MCurrency: make(map[string]okex.AccountWalletResult),
	}
}

func (w *Wallet) LoadCurrency() error {
	results, err := w.client.GetWallet()
	if err != nil {
		return err
	}
	for _, result := range results {
		w.MCurrency[result.Currency] = result
	}
	return nil
}

func (w *Wallet) GetCurrency(currency string) (okex.AccountWalletResult, error) {
	return w.client.GetWalletByCurrency(currency)
}

func (w *Wallet) GetCurrencyNames() []string {
	var strs []string
	for currency, _ := range w.MCurrency {
		strs = append(strs, currency)
	}
	return strs
}
