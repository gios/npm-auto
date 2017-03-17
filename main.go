package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func (npm *npmWriter) WriteToPackage() {
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
}

func (npm *npmWriter) WriteToReadme() {
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
}

func (npm *npmWriter) WriteToFiles() {
	npm.IncreaseVersion()
	npm.AddChangelog()
	npm.WriteToPackage()
	npm.WriteToReadme()
}

func main() {
	npm := npmWriter{}
	npm.WriteToFiles()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
