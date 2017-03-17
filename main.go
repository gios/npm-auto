package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	fmt.Println(npm.version)
	packageMap["version"] = string(npm.version)

	packageJSON, _ := json.Marshal(packageMap)
	json.Indent(&outFormatted, packageJSON, "", "\t")
	errWriteFile := ioutil.WriteFile(pwd+"/package.json", outFormatted.Bytes(), 0644)
	errorCheck(errWriteFile)
}

func (npm *npmWriter) WriteToFiles() {
	npm.IncreaseVersion()
	npm.AddChangelog()
	npm.WriteToPackage()
}

func main() {
	npm := npmWriter{}
	npm.WriteToFiles()
	// fmt.Print("Press 'Enter' to continue...")
	// bufio.NewReader(os.Stdin).ReadBytes('\n')
}
