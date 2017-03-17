package main

import (
	"bufio"
	"fmt"
	"os"
)

type npmWriter struct {
	version   string
	changelog string
}

func (npm *npmWriter) IncreaseVersion() {
	versionReader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new version: ")
	version, _ := versionReader.ReadString('\n')
	npm.version = version
}

func (npm *npmWriter) AddChangelog() {
	changelogReader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter changelog message: ")
	message, _ := changelogReader.ReadString('\n')
	npm.changelog = message
}

func main() {
	npm := npmWriter{}
	npm.IncreaseVersion()
	npm.AddChangelog()
	fmt.Println(npm)
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
