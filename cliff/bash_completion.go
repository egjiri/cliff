package cliff

import (
	"fmt"
	"os/user"
	"strings"
)

// GenerateBashCompletionFile creates the bash completion files,
// saves it in the specified path, and prints instructions on how to use it
func GenerateBashCompletionFile(path string) error {
	if path == "" {
		path = fmt.Sprintf("~/.%v-completion", rootCmd.Use)
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	expandedPath := strings.Replace(path, "~", usr.HomeDir, 1)
	if err := rootCmd.GenBashCompletionFile(expandedPath); err != nil {
		return err
	}
	snippet := fmt.Sprintf("if [ -f %v ]; then . %v; fi\n", path, path)
	fmt.Printf("Bash completion script generated!\nAdd the following line to your .bash_profile:\n\n%v", snippet)
	return nil
}
