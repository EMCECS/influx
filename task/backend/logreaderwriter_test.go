package backend_test

import (
	"testing"

	"github.com/EMCECS/influx/task/backend"
	"github.com/EMCECS/influx/task/backend/storetest"
)

func TestInMemRunStore(t *testing.T) {
	storetest.NewRunStoreTest(
		"inmem",
		func(t *testing.T) (backend.LogWriter, backend.LogReader) {
			rw := backend.NewInMemRunReaderWriter()
			return rw, rw
		},
		func(t *testing.T, w backend.LogWriter, r backend.LogReader) {})(t)
}
