package capital

import (
	"github.com/d-ashesss/mah-moneh/model"
	"github.com/d-ashesss/mah-moneh/model/account"
)

func Get(u *model.User) (*model.Capital, error) {
	c := &model.Capital{}
	accs, err := account.GetAll(u)
	if err != nil {
		return nil, err
	}
	for _, acc := range accs {
		amount, err := account.GetAmount(acc)
		if err != nil {
			continue
		}
		c.Amount += amount.Amount
	}
	return c, nil
}
