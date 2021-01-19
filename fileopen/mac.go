// + build darwin

package fileopen

import "os/exec"

//Open opens the file in the correct directory
func Open(path, name string) error {
	cmd := exec.Command("open", path+"/"+name)
	return cmd.Run()
}
