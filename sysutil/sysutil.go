// Package sysutil provide some system util functions. eg: sysenv, exec, user, process
package sysutil

import (
	"errors"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
)

// Workdir get
func Workdir() string {
	dir, _ := os.Getwd()
	return dir
}

// BinDir get
func BinDir() string {
	binFile := os.Args[0]
	return path.Dir(binFile)
}

// BinFile get
func BinFile() string {
	return os.Args[0]
}

// GoVersion get go runtime version. eg: "1.18.2"
func GoVersion() string {
	return runtime.Version()[2:]
}

// GoInfo define
//
// On os by:
//
//	go env GOVERSION GOOS GOARCH
//	go version // "go version go1.19 darwin/amd64"
type GoInfo struct {
	Version string
	GoOS    string
	Arch    string
}

// match "go version go1.19 darwin/amd64"
var goVerRegex = regexp.MustCompile(`\sgo([\d.]+)\s(\w+)/(\w+)`)

// ParseGoVersion get info by parse `go version` results.
//
// Examples:
//
//		line, err := sysutil.ExecLine("go version")
//		if err != nil {
//			return err
//		}
//
//		info, err := sysutil.ParseGoVersion()
//	 	dump.P(info)
func ParseGoVersion(line string) (*GoInfo, error) {
	// eg: [" go1.19 darwin/amd64", "1.19", "darwin", "amd64"]
	lines := goVerRegex.FindStringSubmatch(line)
	if len(lines) != 4 {
		return nil, errors.New("returns go info is not full")
	}

	info := &GoInfo{}
	info.Version = strings.TrimPrefix(lines[1], "go")
	info.GoOS = lines[2]
	info.Arch = lines[3]

	return info, nil
}

// OsGoInfo fetch and parse
func OsGoInfo() (*GoInfo, error) {
	cmdArgs := []string{"env", "GOVERSION", "GOOS", "GOARCH"}
	line, err := ExecCmd("go", cmdArgs)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(line), "\n")

	if len(lines) != len(cmdArgs)-1 {
		return nil, errors.New("returns go info is not full")
	}

	info := &GoInfo{}
	info.Version = strings.TrimPrefix(lines[0], "go")
	info.GoOS = lines[1]
	info.Arch = lines[2]

	return info, nil
}
