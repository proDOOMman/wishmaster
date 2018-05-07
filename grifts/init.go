package grifts

import (
	"github.com/gobuffalo/buffalo"
	"gitlab.com/SML-482HD/wishmaster/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
