package cli

import (
	"fmt"
	"log"
	"os"
	"os/user"
)

// GenerateBashCompletionFile ...
func GenerateBashCompletionFile() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 {
		log.Fatal("CLI name not specified!")
	}
	cliName := rootCmd.Use
	path := fmt.Sprintf("%v/.%v-completion", usr.HomeDir, cliName)
	if err := rootCmd.GenBashCompletionFile(path); err != nil {
		log.Fatal(err)
	}
	path = fmt.Sprintf("~/.%v-completion", cliName)
	snippet := fmt.Sprintf("if [ -f %v ]; then . %v; fi\n", path, path)
	fmt.Printf("Bash completion script generated!\nAdd the following line to your .bash_profile:\n\n%v", snippet)
}
