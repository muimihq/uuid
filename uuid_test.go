package uuid

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	googleuuid "github.com/google/uuid"
)

func TestV7(t *testing.T) {
	got, err := V7()
	if err != nil {
		t.Fatalf("V7() error = %v", err)
	}

	parsed, err := got.ToUUID()
	if err != nil {
		t.Fatalf("ToUUID() error = %v", err)
	}
	if parsed.Version() != 7 {
		t.Errorf("UUID version = %d, want 7", parsed.Version())
	}
	if len(got.String()) != 22 {
		t.Errorf("encoded UUID length = %d, want 22", len(got.String()))
	}
}

func TestValueAndScanRoundTrip(t *testing.T) {
	want, err := V7()
	if err != nil {
		t.Fatalf("V7() error = %v", err)
	}

	value, err := want.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	canonical, ok := value.(string)
	if !ok {
		t.Fatalf("Value() type = %T, want string", value)
	}
	if len(canonical) != 36 {
		t.Errorf("database UUID length = %d, want 36", len(canonical))
	}

	var got UUID
	if err := got.Scan(canonical); err != nil {
		t.Fatalf("Scan(string) error = %v", err)
	}
	if got.String() != want.String() {
		t.Errorf("Scan(Value(%q)) = %q, want %q", want.String(), got.String(), want.String())
	}
}

func TestScanByteSlice(t *testing.T) {
	u, err := V7()
	if err != nil {
		t.Fatalf("V7() error = %v", err)
	}
	value, err := u.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	canonical, ok := value.(string)
	if !ok {
		t.Fatalf("Value() type = %T, want string", value)
	}

	var got UUID
	if err := got.Scan([]byte(canonical)); err != nil {
		t.Fatalf("Scan([]byte) error = %v", err)
	}
	if got.String() != u.String() {
		t.Errorf("Scan([]byte) = %q, want %q", got.String(), u.String())
	}
}

func TestScanRejectsInvalidValues(t *testing.T) {
	nonV7, err := googleuuid.NewRandom()
	if err != nil {
		t.Fatalf("NewRandom() error = %v", err)
	}

	tests := []struct {
		name  string
		input any
	}{
		{name: "nil", input: nil},
		{name: "unsupported type", input: int64(1)},
		{name: "malformed UUID", input: "not-a-uuid"},
		{name: "non-v7 UUID", input: nonV7.String()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got UUID
			if err := got.Scan(tt.input); !errors.Is(err, ErrUUIDError) {
				t.Errorf("Scan(%#v) error = %v, want %v", tt.input, err, ErrUUIDError)
			}
		})
	}
}

func TestScanNilReceiver(t *testing.T) {
	var got *UUID
	if err := got.Scan("019535d8-9d4c-7a13-8c3d-1f3bb0a45f7e"); !errors.Is(err, ErrUUIDError) {
		t.Errorf("Scan() error = %v, want %v", err, ErrUUIDError)
	}
}

func TestValueRejectsInvalidValues(t *testing.T) {
	nonV7, err := googleuuid.NewRandom()
	if err != nil {
		t.Fatalf("NewRandom() error = %v", err)
	}
	nonV7Encoded := v7toUuid(nonV7)

	tests := []struct {
		name  string
		value UUID
	}{
		{name: "malformed base64", value: UUID{v: "not-valid"}},
		{name: "wrong byte length", value: UUID{v: base64.RawURLEncoding.EncodeToString([]byte{1, 2})}},
		{name: "non-v7 UUID", value: nonV7Encoded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.value.Value(); !errors.Is(err, ErrUUIDError) {
				t.Errorf("Value() error = %v, want %v", err, ErrUUIDError)
			}
		})
	}
}

func TestParseAndFromUUID(t *testing.T) {
	want, err := V7()
	if err != nil {
		t.Fatalf("V7() error = %v", err)
	}

	got, err := Parse(want.String())
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if got.String() != want.String() {
		t.Errorf("Parse() = %q, want %q", got.String(), want.String())
	}

	googleValue, err := want.ToUUID()
	if err != nil {
		t.Fatalf("ToUUID() error = %v", err)
	}
	fromUUID, err := FromUUID(googleValue)
	if err != nil {
		t.Fatalf("FromUUID() error = %v", err)
	}
	if fromUUID.String() != want.String() {
		t.Errorf("FromUUID() = %q, want %q", fromUUID.String(), want.String())
	}
}

func TestParseRejectsInvalidValues(t *testing.T) {
	nonV7, err := googleuuid.NewRandom()
	if err != nil {
		t.Fatalf("NewRandom() error = %v", err)
	}

	tests := []struct {
		name  string
		value string
	}{
		{name: "malformed base64", value: "not-valid"},
		{name: "wrong byte length", value: base64.RawURLEncoding.EncodeToString([]byte{1, 2})},
		{name: "non-v7 UUID", value: base64.RawURLEncoding.EncodeToString(nonV7[:])},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.value); !errors.Is(err, ErrUUIDError) {
				t.Errorf("Parse(%q) error = %v, want %v", tt.value, err, ErrUUIDError)
			}
		})
	}
}

func TestFromUUIDRejectsNonV7(t *testing.T) {
	nonV7, err := googleuuid.NewRandom()
	if err != nil {
		t.Fatalf("NewRandom() error = %v", err)
	}
	if _, err := FromUUID(nonV7); !errors.Is(err, ErrUUIDError) {
		t.Errorf("FromUUID() error = %v, want %v", err, ErrUUIDError)
	}
}

func TestJSONMarshaler(t *testing.T) {
	want, err := V7()
	if err != nil {
		t.Fatalf("V7() error = %v", err)
	}

	encoded, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	expected, err := json.Marshal(want.String())
	if err != nil {
		t.Fatalf("json.Marshal(string) error = %v", err)
	}
	if !bytes.Equal(encoded, expected) {
		t.Errorf("json.Marshal() = %s, want %s", encoded, expected)
	}

	var got UUID
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got != want {
		t.Errorf("json.Unmarshal() = %q, want %q", got.String(), want.String())
	}
}

func TestTextMarshaler(t *testing.T) {
	want, err := V7()
	if err != nil {
		t.Fatalf("V7() error = %v", err)
	}

	encoded, err := want.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() error = %v", err)
	}
	if string(encoded) != want.String() {
		t.Errorf("MarshalText() = %q, want %q", encoded, want.String())
	}

	var got UUID
	if err := got.UnmarshalText(encoded); err != nil {
		t.Fatalf("UnmarshalText() error = %v", err)
	}
	if got != want {
		t.Errorf("UnmarshalText() = %q, want %q", got.String(), want.String())
	}
}
