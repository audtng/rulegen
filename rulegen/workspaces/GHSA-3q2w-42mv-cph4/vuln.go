package main

package runner

import (
	"os/exec"

	"github.com/filebrowser/filebrowser/v2/settings"
)

// ParseCommand parses the command taking in account if the current
// instance uses a shell to run the commands or just calls the binary
// directly.
func ParseCommand(s *settings.Settings, raw string) ([]string, error) {
	var command []string

	if len(s.Shell) == 0 || s.Shell[0] == "" {
		cmd, args, err := SplitCommandAndArgs(raw)
		if err != nil {
			return nil, err
		}

		_, err = exec.LookPath(cmd)
		if err != nil {
			return nil, err
		}

		command = append(command, cmd)
		command = append(command, args...)
	} else {
		command = append(s.Shell, raw) //nolint:gocritic
	}

	return command, nil
}
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

		}
	}

	command, err := runner.ParseCommand(d.settings, raw)
	if err != nil {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil { //nolint:govet
			wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	if !d.server.EnableExec || !d.user.CanExecute(command[0]) {
		if err := conn.WriteMessage(websocket.TextMessage, cmdNotAllowed); err != nil { //nolint:govet
			wsErr(conn, r, http.StatusInternalServerError, err)
		}

import (
	"path/filepath"
	"slices"

	"github.com/spf13/afero"

func (u *User) FullPath(path string) string {
	return afero.FullBaseFsPath(u.Fs.(*afero.BasePathFs), path)
}

// CanExecute checks if an user can execute a specific command.
func (u *User) CanExecute(command string) bool {
	if !u.Perm.Execute {
		return false
	}

	return slices.Contains(u.Commands, command)
}
