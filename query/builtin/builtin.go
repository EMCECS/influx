// Package builtin ensures all packages related to Flux built-ins are imported and initialized.
// This should only be imported from main or test packages.
// It is a mistake to import it from any other package.
package builtin

import (
	"github.com/EMCECS/flux"

	_ "github.com/EMCECS/flux/functions" // Import the built-in functions
	_ "github.com/EMCECS/flux/functions/inputs"
	_ "github.com/EMCECS/flux/functions/outputs"
	_ "github.com/EMCECS/flux/functions/transformations"
	_ "github.com/EMCECS/flux/options"           // Import the built-in options
	_ "github.com/EMCECS/influx/query/functions" // Import the built-in functions
	_ "github.com/EMCECS/influx/query/functions/inputs"
	_ "github.com/EMCECS/influx/query/functions/outputs"
	_ "github.com/EMCECS/influx/query/options" // Import the built-in options
)

func init() {
	flux.FinalizeBuiltIns()
}
