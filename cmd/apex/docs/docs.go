// Package docs outputs Wiki documentation.
package docs

import (
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
	doc "github.com/apex/apex/docs"
	"github.com/apex/apex/stats"
)

// example output.
const example = `  Output documentation
  $ apex docs`

// Command config.
var Command = &cobra.Command{
	Use:              "docs",
	Short:            "Output documentation",
	Example:          example,
	PersistentPreRun: root.PreRunNoop,
	RunE:             run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// Run command.
func run(c *cobra.Command, args []string) (err error) {
	stats.Track("Docs", nil)

	var w io.WriteCloser = os.Stdout

	if isatty.IsTerminal(os.Stdout.Fd()) {
		cmd := exec.Command("less", "-R")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		defer cmd.Wait()

		w, err = cmd.StdinPipe()
		if err != nil {
			return err
		}
		defer w.Close()

		if err := cmd.Start(); err != nil {
			return err
		}
	}

	_, err = io.Copy(w, doc.Reader())
	return err
}
