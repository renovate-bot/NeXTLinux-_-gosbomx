package options

import (
	"github.com/nextlinux/gosbom/internal"
	"github.com/nextlinux/gosbom/internal/config"
	"github.com/nextlinux/gosbom/internal/log"
)

func IsVerbose(app *config.Application) (result bool) {
	isPipedInput, err := internal.IsPipedInput()
	if err != nil {
		// since we can't tell if there was piped input we assume that there could be to disable the ETUI
		log.Warnf("unable to determine if there is piped input: %+v", err)
		return true
	}
	// verbosity should consider if there is piped input (in which case we should not show the ETUI)
	return app.Verbosity > 0 || isPipedInput
}
