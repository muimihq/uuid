package uuid

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
)

var ErrUUIDError = errors.New("invalid UUID")

// UUID is a UUIDv7 represented as unpadded URL-safe base64.
// Its database value is the canonical UUID string.
type UUID string

var (
	_ sql.Scanner   = (*UUID)(nil)
	_ driver.Valuer = UUID("")
)

func (u *UUID) Scan(v any) error {
	if u == nil {
		return ErrUUIDError
	}

	var value string
	switch v := v.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	default:
		return ErrUUIDError
	}

	vU, err := parseV7(value)
	if err != nil {
		return ErrUUIDError
	}
	*u = v7toUuid(vU)
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	vU, err := uuidtoV7(u)
	if err != nil {
		return nil, ErrUUIDError
	}
	return vU.String(), nil
}

func (u UUID) String() string {
	return string(u)
}

// V7 generates a UUIDv7 represented as unpadded URL-safe base64.
func V7() (UUID, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return v7toUuid(u), nil
}

func v7toUuid(u uuid.UUID) UUID {
	return UUID(base64.RawURLEncoding.EncodeToString(u[:]))
}

func uuidtoV7(v UUID) (uuid.UUID, error) {
	uBytes, err := base64.RawURLEncoding.DecodeString(string(v))
	if err != nil {
		return uuid.UUID{}, ErrUUIDError
	}
	u, err := uuid.FromBytes(uBytes)
	if err != nil || u.Version() != 7 {
		return uuid.UUID{}, ErrUUIDError
	}
	return u, nil
}

func parseV7(value string) (uuid.UUID, error) {
	u, err := uuid.Parse(value)
	if err != nil || u.Version() != 7 {
		return uuid.UUID{}, ErrUUIDError
	}
	return u, nil
}
