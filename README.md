# uuid

`github.com/muimihq/uuid` generates UUIDv7 identifiers represented as a compact,
22-character unpadded URL-safe Base64 string. The `UUID` type hides this
representation behind methods. For `database/sql`, it stores the canonical
36-character UUID string and converts it back when scanned.

## Usage

```go
id, err := uuid.V7()
if err != nil {
    return err
}

_, err = db.ExecContext(ctx, `INSERT INTO records (id) VALUES (?)`, id)
if err != nil {
    return err
}

var loaded uuid.UUID
err = db.QueryRowContext(ctx, `SELECT id FROM records WHERE id = ?`, id).
    Scan(&loaded)
if err != nil {
    return err
}
```

`UUID` implements `database/sql/driver.Valuer` as a value and
`database/sql.Scanner` as a pointer. Scanning accepts a canonical UUID string or
its byte slice form. `Parse` accepts the package's URL-safe Base64
representation; `FromUUID` converts a `github.com/google/uuid.UUID`. These
conversions and the SQL methods reject malformed values and UUIDs that are not
version 7 with `uuid.ErrUUIDError`. The zero value and `uuid.Nil` are invalid
identifiers and cannot be stored.
