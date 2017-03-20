package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type npmWriter struct {
	version   string
	changelog string
}

func errorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func (npm *npmWriter) IncreaseVersion() {
	versionReader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new version: ")
	version, _, _ := versionReader.ReadLine()
	npm.version = string(version)
}

func (npm *npmWriter) AddChangelog() {
	changelogReader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter changelog message: ")
	message, _, _ := changelogReader.ReadLine()
	npm.changelog = string(message)
}

func (npm *npmWriter) WriteToPackage() bool {
	var packageMap map[string]interface{}
	var outFormatted bytes.Buffer

	pwd, _ := os.Getwd()
	packageData, errReadFile := ioutil.ReadFile(pwd + "/package.json")
	errorCheck(errReadFile)

	errJSONDecode := json.Unmarshal(packageData, &packageMap)
	errorCheck(errJSONDecode)
	packageMap["version"] = string(npm.version)

	packageJSON, _ := json.Marshal(packageMap)
	json.Indent(&outFormatted, packageJSON, "", "\t")
	errWriteFile := ioutil.WriteFile(pwd+"/package.json", outFormatted.Bytes(), 0644)
	errorCheck(errWriteFile)
	fmt.Println("package.json has been updated!")
	return true
}

func (npm *npmWriter) WriteToReadme() bool {
	var newVersionStr string

	pwd, _ := os.Getwd()
	readmeData, errReadFile := ioutil.ReadFile(pwd + "/README.md")
	errorCheck(errReadFile)

	readmeHeader := strings.Split(string(readmeData), "\n")[0]
	regexpVersion := regexp.MustCompile("(v(\\d+\\.)?(\\d+\\.)?(\\*|\\d+))")
	oldVersion := regexpVersion.FindString(readmeHeader)

	newVersionStr += npm.version
	outReadme := strings.Replace(string(readmeData), oldVersion, "v"+newVersionStr, 1)
	errWriteFile := ioutil.WriteFile(pwd+"/README.md", []byte(outReadme), 0644)
	errorCheck(errWriteFile)
	fmt.Println("README.md has been updated!")
	return true
}

func (npm *npmWriter) WriteToChangelog() bool {
	pwd, _ := os.Getwd()
	changelogData, errReadFile := ioutil.ReadFile(pwd + "/CHANGELOG.md")
	errorCheck(errReadFile)

	changelogSplitted := strings.Split(string(changelogData), "\n")
	changelogAdder := make([]string, 4)
	changelogAdder[1] = "## " + npm.version
	changelogAdder[3] = "- " + npm.changelog

	changelogSplitted = append(changelogSplitted[:1], append(changelogAdder, changelogSplitted[1:]...)...)
	changelogJoined := strings.Join(changelogSplitted, "\n")
	errWriteFile := ioutil.WriteFile(pwd+"/CHANGELOG.md", []byte(changelogJoined), 0644)
	errorCheck(errWriteFile)
	fmt.Println("CHANGELOG.md has been updated!")
	return true
}

func (npm *npmWriter) addTag() {
	cmd := exec.Command("git", "tag", "v"+npm.version, "-f")
	err := cmd.Run()
	errorCheck(err)
	fmt.Println("-----------------------------")
	fmt.Printf("Tag has been added")
}

func (npm *npmWriter) Finish(loaderPackage, loaderReadme, loaderChangelog bool) {
	if loaderPackage && loaderReadme && loaderChangelog {
		fmt.Println("-----------------------------")
		fmt.Printf("Version successfully updated to: %v\n", npm.version)
	} else {
		fmt.Println("Oops something has gone wrong, please try again!")
	}
}

func (npm *npmWriter) WriteToFiles() {
	npm.IncreaseVersion()
	npm.AddChangelog()
	fmt.Println("-----------------------------")
	npm.Finish(npm.WriteToPackage(), npm.WriteToReadme(), npm.WriteToChangelog())
	npm.addTag()
}

func main() {
	npm := npmWriter{}
	npm.WriteToFiles()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
