package uuid

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
)

type ID uuid.UUID

var Nil ID

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func Parse(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return Nil, domainerror.New(domainerror.InvalidParams, "invalid id", map[string]any{"id": s, "error": err.Error()})
	}

	return ID(id), nil
}

func MustParse(s string) ID {
	return ID(uuid.MustParse(s))
}

func (id ID) IsZero() bool {
	return id == Nil
}

func (id *ID) Scan(value any) error {
	if value == nil {
		*id = Nil

		return nil
	}

	switch v := value.(type) {
	case []byte:
		parsedID, err := uuid.ParseBytes(v)
		if err != nil {
			return err
		}

		*id = ID(parsedID)
	case string:
		parsedID, err := uuid.Parse(v)
		if err != nil {
			return err
		}

		*id = ID(parsedID)
	default:
		return fmt.Errorf("unsupported Scan source: %T", value)
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*id = ID(parsedID)

	return nil
}

// UnmarshalText implementa encoding.TextUnmarshaler.
// Isso permite que uuid.ID seja convertido de query params automaticamente.
func (id *ID) UnmarshalText(text []byte) error {
	parsed, err := uuid.Parse(string(text))
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}
	*id = ID(parsed)
	return nil
}

func NextID() ID {
	return idGenerator.NextID()
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return Nil, err
	}

	return ID(id), nil
}
