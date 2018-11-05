package outputs

import "github.com/EMCECS/influx/telegraf/plugins"

type baseOutput int

func (b baseOutput) Type() plugins.Type {
	return plugins.Output
}
