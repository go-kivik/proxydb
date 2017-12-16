package proxydb

import (
	"context"
	"encoding/json"
	"io"

	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kivik/driver"
	"github.com/go-kivik/kivik/errors"
)

var notYetImplemented = errors.Status(kivik.StatusNotImplemented, "kivik: not yet implemented in proxy driver")

// CompleteClient is a composite of all compulsory and optional driver.* client
// interfaces.
type CompleteClient interface {
	driver.Client
	driver.Authenticator
}

// NewClient wraps an existing *kivik.Client connection, allowing it to be used
// as a driver.Client
func NewClient(c *kivik.Client) driver.Client {
	return &client{c}
}

type client struct {
	*kivik.Client
}

var _ CompleteClient = &client{}

func (c *client) AllDBs(ctx context.Context, options map[string]interface{}) ([]string, error) {
	return c.Client.AllDBs(ctx, options)
}

func (c *client) CreateDB(ctx context.Context, dbname string, options map[string]interface{}) error {
	return c.Client.CreateDB(ctx, dbname, options)
}

func (c *client) DBExists(ctx context.Context, dbname string, options map[string]interface{}) (bool, error) {
	return c.Client.DBExists(ctx, dbname, options)
}

func (c *client) DestroyDB(ctx context.Context, dbname string, options map[string]interface{}) error {
	return c.Client.DestroyDB(ctx, dbname, options)
}

func (c *client) Version(ctx context.Context) (*driver.Version, error) {
	ver, err := c.Client.Version(ctx)
	if err != nil {
		return nil, err
	}
	return &driver.Version{
		Version:     ver.Version,
		Vendor:      ver.Vendor,
		RawResponse: ver.RawResponse,
	}, nil
}

func (c *client) DB(ctx context.Context, name string, options map[string]interface{}) (driver.DB, error) {
	d, err := c.Client.DB(ctx, name, options)
	return &db{d}, err
}

type db struct {
	*kivik.DB
}

var _ driver.DB = &db{}
var _ driver.DBOpts = &db{}

func (d *db) AllDocs(ctx context.Context, opts map[string]interface{}) (driver.Rows, error) {
	kivikRows, err := d.DB.AllDocs(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &rows{kivikRows}, nil
}

func (d *db) Query(ctx context.Context, ddoc, view string, opts map[string]interface{}) (driver.Rows, error) {
	kivikRows, err := d.DB.Query(ctx, ddoc, view, opts)
	if err != nil {
		return nil, err
	}
	return &rows{kivikRows}, nil
}

func (d *db) Get(ctx context.Context, id string, opts map[string]interface{}) (json.RawMessage, error) {
	row, err := d.DB.Get(ctx, id, opts)
	if err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err = row.ScanDoc(&raw)
	return raw, err
}

func (d *db) Stats(ctx context.Context) (*driver.DBStats, error) {
	i, err := d.DB.Stats(ctx)
	stats := driver.DBStats(*i)
	return &stats, err
}

func (d *db) Security(ctx context.Context) (*driver.Security, error) {
	s, err := d.DB.Security(ctx)
	if err != nil {
		return nil, err
	}
	sec := driver.Security{
		Admins:  driver.Members(s.Admins),
		Members: driver.Members(s.Members),
	}
	return &sec, err
}

func (d *db) SetSecurity(ctx context.Context, security *driver.Security) error {
	sec := &kivik.Security{
		Admins:  kivik.Members(security.Admins),
		Members: kivik.Members(security.Members),
	}
	return d.DB.SetSecurity(ctx, sec)
}

func (d *db) Changes(ctx context.Context, opts map[string]interface{}) (driver.Changes, error) {
	return nil, notYetImplemented
}

func (d *db) BulkDocs(_ context.Context, _ []interface{}) (driver.BulkResults, error) {
	// FIXME: Unimplemented
	return nil, notYetImplemented
}

func (d *db) PutAttachment(_ context.Context, _, _, _, _ string, _ io.Reader) (string, error) {
	panic("PutAttachment should never be called")
}

func (d *db) PutAttachmentOpts(_ context.Context, _, _, _, _ string, _ io.Reader, opts map[string]interface{}) (string, error) {
	// FIXME: Unimplemented
	return "", notYetImplemented
}

func (d *db) GetAttachment(ctx context.Context, docID, rev, filename string) (contentType string, md5sum driver.MD5sum, body io.ReadCloser, err error) {
	panic("GetAttachment should never be called")
}

func (d *db) GetAttachmentOpts(ctx context.Context, docID, rev, filename string, opts map[string]interface{}) (contentType string, md5sum driver.MD5sum, body io.ReadCloser, err error) {
	// FIXME: Unimplemented
	return "", [16]byte{}, nil, notYetImplemented
}

func (d *db) CreateDoc(_ context.Context, _ interface{}) (string, string, error) {
	panic("CreateDoc should never be called")
}

func (d *db) CreateDocOpts(ctx context.Context, doc interface{}, opts map[string]interface{}) (string, string, error) {
	return d.DB.CreateDoc(ctx, doc, opts)
}

func (d *db) Delete(_ context.Context, _, _ string) (string, error) {
	panic("Delete should never be called")
}

func (d *db) DeleteOpts(ctx context.Context, id, rev string, opts map[string]interface{}) (string, error) {
	return d.DB.Delete(ctx, id, rev, opts)
}

func (d *db) DeleteAttachment(_ context.Context, _, _, _ string) (string, error) {
	panic("DeleteAttachment should never be called")
}

func (d *db) DeleteAttachmentOpts(ctx context.Context, id, rev, filename string, opts map[string]interface{}) (string, error) {
	return d.DB.DeleteAttachment(ctx, id, rev, filename, opts)
}

func (d *db) Put(_ context.Context, _ string, _ interface{}) (string, error) {
	panic("Put should never be called")
}

func (d *db) PutOpts(ctx context.Context, id string, doc interface{}, opts map[string]interface{}) (string, error) {
	return d.DB.Put(ctx, id, doc, opts)
}
