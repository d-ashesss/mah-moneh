package currencies

type Rate struct {
	Base      string `gorm:"primaryKey"`
	Target    string `gorm:"primaryKey"`
	YearMonth string `gorm:"primaryKey"`
	Rate      float64
}
