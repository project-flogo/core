package sqltx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"sync"
	"sync/atomic"
)

// A minimal in-memory database/sql driver. No real database is available, and adding a mock
// dependency to core is not acceptable, so the tests drive the genuine database/sql machinery -
// pooling, connection pinning, Tx.PrepareContext, Commit/Rollback - against this.
//
// It deliberately implements driver.ConnBeginTx and driver.ConnPrepareContext, because those are
// exactly the interfaces the real connectors' drivers implement (verified for mssql, mysql,
// godror and lib/pq) and they are what makes prepare-on-transaction work.

type fakeDriver struct {
	mu        sync.Mutex
	opened    int32 // connections handed out, cumulative
	live      int32 // connections currently open
	prepares  int32 // PrepareContext calls, cumulative
	commits   int32
	rollbacks int32
}

func (d *fakeDriver) Open(string) (driver.Conn, error) {
	atomic.AddInt32(&d.opened, 1)
	atomic.AddInt32(&d.live, 1)
	return &fakeConn{drv: d}, nil
}

func (d *fakeDriver) counts() (opened, live, prepares, commits, rollbacks int32) {
	return atomic.LoadInt32(&d.opened), atomic.LoadInt32(&d.live),
		atomic.LoadInt32(&d.prepares), atomic.LoadInt32(&d.commits),
		atomic.LoadInt32(&d.rollbacks)
}

type fakeConn struct {
	drv   *fakeDriver
	inTx  bool
	badTx bool
}

var (
	_ driver.Conn               = (*fakeConn)(nil)
	_ driver.ConnBeginTx        = (*fakeConn)(nil)
	_ driver.ConnPrepareContext = (*fakeConn)(nil)
)

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

func (c *fakeConn) PrepareContext(_ context.Context, query string) (driver.Stmt, error) {
	atomic.AddInt32(&c.drv.prepares, 1)
	return &fakeStmt{conn: c, query: query}, nil
}

func (c *fakeConn) Close() error {
	atomic.AddInt32(&c.drv.live, -1)
	return nil
}

func (c *fakeConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *fakeConn) BeginTx(_ context.Context, _ driver.TxOptions) (driver.Tx, error) {
	if c.inTx {
		// Mirrors godror conn.go:337-341, which rejects a nested begin on the same connection.
		return nil, driver.ErrBadConn
	}
	c.inTx = true
	return &fakeTx{conn: c}, nil
}

type fakeTx struct{ conn *fakeConn }

func (t *fakeTx) Commit() error {
	t.conn.inTx = false
	atomic.AddInt32(&t.conn.drv.commits, 1)
	if t.conn.badTx {
		return driver.ErrBadConn
	}
	return nil
}

func (t *fakeTx) Rollback() error {
	t.conn.inTx = false
	atomic.AddInt32(&t.conn.drv.rollbacks, 1)
	return nil
}

type fakeStmt struct {
	conn  *fakeConn
	query string
}

func (s *fakeStmt) Close() error  { return nil }
func (s *fakeStmt) NumInput() int { return -1 }

func (s *fakeStmt) Exec([]driver.Value) (driver.Result, error) {
	return driver.RowsAffected(1), nil
}

func (s *fakeStmt) Query([]driver.Value) (driver.Rows, error) {
	return &fakeRows{cols: []string{"c"}, rows: [][]driver.Value{{int64(1)}}}, nil
}

type fakeRows struct {
	cols []string
	rows [][]driver.Value
	i    int
}

func (r *fakeRows) Columns() []string { return r.cols }
func (r *fakeRows) Close() error      { return nil }

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.i >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.i])
	r.i++
	return nil
}

// newFakeDB registers a uniquely-named driver and opens a pool on it. The unique name matters:
// sql.Register panics on a duplicate, and these tests each want isolated counters.
var fakeSeq int32

func newFakeDB(maxOpen int) (*sql.DB, *fakeDriver) {
	d := &fakeDriver{}
	name := "sqltx-fake-" + itoa(int(atomic.AddInt32(&fakeSeq, 1)))
	sql.Register(name, d)

	db, err := sql.Open(name, "")
	if err != nil {
		panic(err)
	}
	if maxOpen > 0 {
		db.SetMaxOpenConns(maxOpen)
	}
	return db, d
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
