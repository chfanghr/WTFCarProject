package uuid

import (
	"database/sql/driver"
	"fmt"
)

func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		if len(src) == Size {
			return u.UnmarshalBinary(src)
		}
		return u.UnmarshalText(src)

	case string:
		return u.UnmarshalText([]byte(src))
	}

	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

type NullUUID struct {
	UUID  UUID
	Valid bool
}

func (u NullUUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}

	return u.UUID.Value()
}

func (u *NullUUID) Scan(src interface{}) error {
	if src == nil {
		u.UUID, u.Valid = Nil, false
		return nil
	}
	u.Valid = true
	return u.UUID.Scan(src)
}
