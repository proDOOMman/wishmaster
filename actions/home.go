package actions

import (
	"context"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

func RebootHandler(c buffalo.Context) error {
	exec.Command("reboot").Start()
	c.Flash().Add("success", "Приставка перезагружается...")
	return c.Redirect(http.StatusSeeOther, "/")
}

func ExecuteHandler(c buffalo.Context) error {
	cmd := c.Request().FormValue("command")
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	messageData, err := exec.CommandContext(ctx, head, parts...).CombinedOutput()

	if err == nil {
		c.Flash().Add("success", string(messageData))
	} else {
		c.Flash().Add("danger", err.Error()+" : "+string(messageData))
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
