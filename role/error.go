package role

import (
	"fmt"
)

type ErrNotFound struct {
	Name string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("role '%s' não encontrada", e.Name)
}
