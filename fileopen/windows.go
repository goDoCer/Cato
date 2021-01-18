// + build windows

package fileopen

import (
	"os/exec"
)

//Open opens the file in the correct directory
func Open(path, name string) error {
	cmd := exec.Command("explorer.exe", name)
	cmd.Dir = path
	return cmd.Run()
}
