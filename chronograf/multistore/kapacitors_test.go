package multistore

import (
	"testing"

	"github.com/EMCECS/influx/chronograf"
)

func TestInterfaceImplementation(t *testing.T) {
	var _ chronograf.ServersStore = &KapacitorStore{}
}
