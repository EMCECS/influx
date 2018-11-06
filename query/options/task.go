// Package options implements flux options.
package options

import (
	"github.com/EMCECS/flux"
	"github.com/EMCECS/flux/values"
)

func init() {
	flux.RegisterBuiltInOption("task", taskObject())
}

func taskObject() values.Object {
	obj := values.NewObject()

	obj.Set("name", values.NewString(""))
	obj.Set("cron", values.NewString(""))
	obj.Set("every", values.NewDuration(0))
	obj.Set("delay", values.NewDuration(0))
	obj.Set("concurrency", values.NewInt(0))
	obj.Set("retry", values.NewInt(0))
	return obj
}
