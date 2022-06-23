package repo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"typewriter/config"
)

func DownloadContents() {
	err := os.Mkdir("./contents", 0755)
	if err != nil {
		panic(err)
	}

	// Configure git
	log.Println("Initializing git configuration")
	out, err := exec.Command("git", "config", "--global", "user.email", fmt.Sprintf("\"%v\"", config.Config.Git.Email)).CombinedOutput()
	if err != nil {
		log.Printf("Git config email: %v", string(out))
		panic(err)
	}
	out, err = exec.Command("git", "config", "--global", "user.name", fmt.Sprintf("\"%v\"", config.Config.Git.Name)).CombinedOutput()
	if err != nil {
		log.Printf("Git config name: %v", string(out))
		panic(err)
	}

	prefix := "https://"
	if strings.HasPrefix(config.Config.Git.URL, "http://") {
		prefix = "http://"
	} else {
		if !strings.HasPrefix(config.Config.Git.URL, "https://") {
			config.Config.Git.URL = "https://" + config.Config.Git.URL
		}
	}
	domain := config.Config.Git.URL[len(prefix):]
	out, err = exec.Command("git", "clone", "--depth", "1", "-b", config.Config.Git.Branch, fmt.Sprintf("%v%v:%v@%v", prefix, config.Config.Git.Username, config.Config.Git.Password, domain), "./contents").CombinedOutput()
	if err != nil {
		log.Printf("Git clone: %v", string(out))
		panic(err)
	}
}
