package product

import (
	"fmt"
)

type ErrNotFound struct {
	ID int
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("product '%d' not found", e.ID)
}

type ErrProductsNotFound struct {
	IDs []int
}

func (e *ErrProductsNotFound) Error() string {
	return fmt.Sprintf("products not found: %+v", e.IDs)
}
