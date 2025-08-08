package favorite

import (
	"time"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

type Favorite struct {
	ClientID    uuid.ID   `json:"clientId"`
	ProductID   int       `json:"productId"`
	RegistredAt time.Time `json:"registredAt"`
}

func (f Favorite) validate() error {
	v := validator.New()

	if f.ClientID.IsZero() {
		v.AddError("clientId", "campo obrigatório")
	}

	if f.ProductID == 0 {
		v.AddError("productId", "campo obrigatório")
	}

	return v.Validate()
}

func New(clientID uuid.ID, productID int) (Favorite, error) {
	f := Favorite{
		ClientID:    clientID,
		ProductID:   productID,
		RegistredAt: time.Now(),
	}

	if err := f.validate(); err != nil {
		return Favorite{}, err
	}

	return f, nil
}
