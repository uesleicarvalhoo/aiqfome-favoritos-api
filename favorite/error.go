package favorite

import (
	"fmt"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type ErrFavoriteNotFound struct {
	ClientID  uuid.ID
	ProductID int
}

func (e *ErrFavoriteNotFound) Error() string {
	return fmt.Sprintf("client '%s' don't have the product with id '%d' on their favorites", e.ClientID.String(), e.ProductID)
}
