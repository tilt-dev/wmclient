package analytics

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/windmilleng/wmclient/pkg/dirs"
)

func Init(appName string) (Analytics, *cobra.Command, error) {
	c, err := initCLI()

	d, err := dirs.UseWindmillDir()
	if err != nil {
		return nil, nil, err
	}
	appDir := filepath.Join("analytics", appName)
	if err := d.MkdirAll(appDir); err != nil {
		return nil, nil, err
	}

	dbName, err := d.Abs(filepath.Join(appDir, "db"))
	if err != nil {
		return nil, nil, err
	}

	b, err := newBoltByteEventStore(dbName)
	if err != nil {
		return nil, nil, err
	}
	m, err := NewMarshalingEventStore(b)
	if err != nil {
		return nil, nil, err
	}
	s := newMemoryEventStore()
	comp := newCompositeEventWriter(s, m)
	return newStoreAnalytics(comp), c, err
}

type UserCollection int

const (
	CollectDefault UserCollection = iota
	CollectNotEvenLocal
	CollectLocal
	CollectUpload
)

var choices = map[UserCollection]string{
	CollectDefault:      "default",
	CollectNotEvenLocal: "not-even-local",
	CollectLocal:        "local",
	CollectUpload:       "upload",
}

type Analytics interface {
	Register(name string, agg Aggregation) (AnyWriter, error)
	Flush() error
}

type AnyWriter interface {
	Write(v interface{})
}

type StringWriter interface {
	Write(v string)
}

type ErrorWriter interface {
	Write(v error)
}

type NoopAnyWriter struct {
}

func NewNoopAnyWriter() AnyWriter {
	return &NoopAnyWriter{}
}

func (w *NoopAnyWriter) Write(v interface{}) {}

type StoreAnyWriter struct {
	store EventWriter
	key   string
}

func (w *StoreAnyWriter) Write(v interface{}) {
	w.store.Write(Event{Key: w.key, Value: v, T: time.Now().UTC()})
}

type DelegatingStringWriter struct {
	del AnyWriter
}

func NewStringWriter(w AnyWriter) StringWriter {
	return &DelegatingStringWriter{del: w}
}

func (w *DelegatingStringWriter) Write(v string) {
	w.del.Write(v)
}

type DelegatingErrorWriter struct {
	del AnyWriter
}

func NewErrorWriter(w AnyWriter) ErrorWriter {
	return &DelegatingErrorWriter{del: w}
}

func (w *DelegatingErrorWriter) Write(v error) {
	w.del.Write(v)
}

type StoreAnalytics struct {
	registered map[string]bool
	store      EventWriter
}

func newStoreAnalytics(store EventWriter) *StoreAnalytics {
	return &StoreAnalytics{
		registered: make(map[string]bool),
		store:      store,
	}
}

func (a *StoreAnalytics) Register(name string, agg Aggregation) (AnyWriter, error) {
	if a.registered[name] {
		return nil, fmt.Errorf("duplicate key %v", name)
	}

	return &StoreAnyWriter{store: a.store, key: name}, nil
}

func (a *StoreAnalytics) Flush() error {
	return a.store.Flush()
}