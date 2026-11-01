package uuid

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
)

// ErrUUIDError is returned when a value is malformed or is not a UUIDv7.
var ErrUUIDError = errors.New("invalid UUID")

// UUID is a UUIDv7 represented as unpadded URL-safe base64.
// Its database value is the canonical UUID string.
type UUID struct {
	v string
}

var (
	// Nil is the zero value and is not a valid UUID.
	Nil UUID

	_ sql.Scanner   = (*UUID)(nil)
	_ driver.Valuer = Nil
)

// Scan implements sql.Scanner.
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

	vU, err := uuid.Parse(value)
	if err != nil {
		return ErrUUIDError
	}
	parsed, err := FromUUID(vU)
	if err != nil {
		return ErrUUIDError
	}
	*u = parsed
	return nil
}

// Value implements driver.Valuer.
func (u UUID) Value() (driver.Value, error) {
	vU, err := u.ToUUID()
	if err != nil {
		return nil, ErrUUIDError
	}
	return vU.String(), nil
}

// String returns the unpadded URL-safe base64 representation of u.
func (u UUID) String() string {
	return u.v
}

// ToUUID converts u to a github.com/google/uuid UUID.
func (u UUID) ToUUID() (uuid.UUID, error) {
	uBytes, err := base64.RawURLEncoding.DecodeString(u.v)
	if err != nil {
		return uuid.UUID{}, ErrUUIDError
	}
	parsed, err := uuid.FromBytes(uBytes)
	if err != nil || parsed.Version() != 7 {
		return uuid.UUID{}, ErrUUIDError
	}
	return parsed, nil
}

// V7 generates a UUIDv7 represented as unpadded URL-safe base64.
func V7() (UUID, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return Nil, err
	}
	return v7toUuid(u), nil
}

// Parse parses an unpadded URL-safe base64 UUIDv7.
func Parse(v string) (UUID, error) {
	uBytes, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return Nil, ErrUUIDError
	}
	u, err := uuid.FromBytes(uBytes)
	if err != nil || u.Version() != 7 {
		return Nil, ErrUUIDError
	}

	return v7toUuid(u), nil
}

// FromUUID converts a github.com/google/uuid UUIDv7 to UUID.
func FromUUID(v uuid.UUID) (UUID, error) {
	if v.Version() != 7 {
		return Nil, ErrUUIDError
	}
	return v7toUuid(v), nil
}

func v7toUuid(u uuid.UUID) UUID {
	return UUID{base64.RawURLEncoding.EncodeToString(u[:])}
}
