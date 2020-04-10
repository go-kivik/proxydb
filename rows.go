package proxydb

import (
	"encoding/json"

	"github.com/go-kivik/kivik/v3"
	"github.com/go-kivik/kivik/v3/driver"
)

type rows struct {
	*kivik.Rows
}

var _ driver.Rows = &rows{}

func (r *rows) Next(row *driver.Row) error {
	if !r.Rows.Next() {
		return r.Rows.Err()
	}
	var value json.RawMessage
	if err := r.Rows.ScanValue(&value); err != nil {
		return err
	}
	var doc json.RawMessage
	if err := r.Rows.ScanDoc(&doc); err != nil {
		return err
	}
	row.ID = r.Rows.ID()
	row.Key = json.RawMessage(r.Rows.Key())
	row.Value = value
	row.Doc = doc
	return nil
}
