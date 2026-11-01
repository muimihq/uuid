package uuid

import (
	"encoding/base64"
	"errors"
	"testing"

	googleuuid "github.com/google/uuid"
)

func TestV7(t *testing.T) {
	got, err := V7()
	if err != nil {
		t.Fatalf("V7() error = %v", err)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(string(got))
	if err != nil {
		t.Fatalf("V7() returned invalid base64 %q: %v", got, err)
	}
	if len(decoded) != 16 {
		t.Fatalf("decoded UUID length = %d, want 16", len(decoded))
	}
	parsed, err := googleuuid.FromBytes(decoded)
	if err != nil {
		t.Fatalf("decode UUID bytes: %v", err)
	}
	if parsed.Version() != 7 {
		t.Errorf("UUID version = %d, want 7", parsed.Version())
	}
	if len(got) != 22 {
		t.Errorf("encoded UUID length = %d, want 22", len(got))
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
	if got != want {
		t.Errorf("Scan(Value(%q)) = %q, want %q", want, got, want)
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
	if got != u {
		t.Errorf("Scan([]byte) = %q, want %q", got, u)
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
		{name: "malformed base64", value: UUID("not-valid")},
		{name: "wrong byte length", value: UUID(base64.RawURLEncoding.EncodeToString([]byte{1, 2}))},
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
