# uuid

`github.com/muimihq/uuid` generates UUIDv7 identifiers and stores them in a
compact, unpadded, URL-safe Base64 representation. The encoded value is 22
characters long. When used with `database/sql`, it is converted to the
canonical 36-character UUID string for storage and converted back when scanned.

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
its byte slice form. Both database conversion methods reject malformed values
and UUIDs that are not version 7 with `uuid.ErrUUIDError`.
