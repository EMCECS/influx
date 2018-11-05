package gather

import (
	"context"

	"github.com/EMCECS/influx"
)

// Scraper gathers metrics from a scraper target.
type Scraper interface {
	Gather(ctx context.Context, target platform.ScraperTarget) (ms []Metrics, err error)
}
