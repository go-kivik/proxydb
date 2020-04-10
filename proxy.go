package proxydb

import (
	"context"
	"net/http"

	"github.com/go-kivik/kivik/v3"
	"github.com/go-kivik/kivik/v3/driver"
	"github.com/go-kivik/kivik/v3/errors"
)

var notYetImplemented = errors.Status(http.StatusNotImplemented, "kivik: not yet implemented in proxy driver")

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
	d := c.Client.DB(ctx, name, options)
	return &db{d}, nil
}

type db struct {
	*kivik.DB
}

var _ driver.DB = &db{}

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

type atts struct {
	*kivik.AttachmentsIterator
}

var _ driver.Attachments = &atts{}

func (a *atts) Close() error { return nil }
func (a *atts) Next(att *driver.Attachment) error {
	next, err := a.AttachmentsIterator.Next()
	if err != nil {
		return err
	}
	*att = driver.Attachment(*next)
	return nil
}

func (d *db) Get(ctx context.Context, id string, opts map[string]interface{}) (*driver.Document, error) {
	row := d.DB.Get(ctx, id, opts)
	return &driver.Document{
		ContentLength: row.ContentLength,
		Rev:           row.Rev,
		Body:          row.Body,
		Attachments:   &atts{row.Attachments},
	}, nil
}

func (d *db) Stats(ctx context.Context) (*driver.DBStats, error) {
	i, err := d.DB.Stats(ctx)
	if err != nil {
		return nil, err
	}
	var cluster *driver.ClusterStats
	if i.Cluster != nil {
		c := driver.ClusterStats(*i.Cluster)
		cluster = &c
	}
	return &driver.DBStats{
		Name:           i.Name,
		CompactRunning: i.CompactRunning,
		DocCount:       i.DocCount,
		DeletedCount:   i.DeletedCount,
		UpdateSeq:      i.UpdateSeq,
		DiskSize:       i.DiskSize,
		ActiveSize:     i.ActiveSize,
		ExternalSize:   i.ExternalSize,
		Cluster:        cluster,
		RawResponse:    i.RawResponse,
	}, nil
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

func (d *db) PutAttachment(_ context.Context, _, _ string, _ *driver.Attachment, _ map[string]interface{}) (string, error) {
	panic("PutAttachment should never be called")
}

func (d *db) GetAttachment(ctx context.Context, docID, filename string, _ map[string]interface{}) (*driver.Attachment, error) {
	panic("GetAttachment should never be called")
}

func (d *db) GetAttachmentMeta(ctx context.Context, docID, rev, filename string, opts map[string]interface{}) (*driver.Attachment, error) {
	// FIXME: Unimplemented
	return nil, notYetImplemented
}

func (d *db) CreateDoc(_ context.Context, _ interface{}, _ map[string]interface{}) (string, string, error) {
	panic("CreateDoc should never be called")
}

func (d *db) Delete(_ context.Context, _, _ string, _ map[string]interface{}) (string, error) {
	panic("Delete should never be called")
}

func (d *db) DeleteAttachment(_ context.Context, _, _, _ string, _ map[string]interface{}) (string, error) {
	panic("DeleteAttachment should never be called")
}

func (d *db) Put(_ context.Context, _ string, _ interface{}, _ map[string]interface{}) (string, error) {
	panic("Put should never be called")
}
