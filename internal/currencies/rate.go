package currencies

import "github.com/d-ashesss/mah-moneh/internal/accounts"

type Rate struct {
	Base      accounts.Currency `gorm:"primaryKey"`
	Target    accounts.Currency `gorm:"primaryKey"`
	YearMonth string            `gorm:"primaryKey"`
	Rate      float64
}
