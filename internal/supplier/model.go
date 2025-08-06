package supplier

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
	Name          string `json:"name"`
	ContactPerson string `json:"contactPerson"`
	Email         string `json:"email" gorm:"unique"`
	Phone         string `json:"phone"`
}
