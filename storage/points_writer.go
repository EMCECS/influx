package storage

import (
	"github.com/EMCECS/influx/models"
)

// PointsWriter describes the ability to write points into a storage engine.
type PointsWriter interface {
	WritePoints([]models.Point) error
}
