package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	// "regexp"
)

// inDir is dir to start. must be full path
const inDir = "/Users/xah/xx_manual/"

// ext is file extension, with the dot. only these are searched
const ext = ".html"
const backupSuffix = "~~"
const writeToFile = false
const doBackup = false

var dirsToSkip = []string{".git"}

type frPair struct {
	fs string // find string
	rs string // replace string
}

var frPairs = []frPair{

	frPair{
		fs: `haskell`,
		rs: `ppppp`,
	},
}

// pass return false if x equals any of y
func pass(x string, y []string) bool {
	for _, v := range y {
		if x == v {
			return false
		}
	}
	return true
}

func doFile(path string) error {
	contentBytes, er := ioutil.ReadFile(path)
	if er != nil {
		panic(er)
	}

	var content = string(contentBytes)

	var changed = false
	for _, pair := range frPairs {
		var found = strings.Index(content, pair.fs)
		if found != -1 {
			content = strings.Replace(content, pair.fs, pair.rs, -1)
			changed = true
		}
	}

	if changed {
		fmt.Printf("changed: %v\n", path)
		if doBackup {
			err := os.Rename(path, path+backupSuffix)
			if err != nil {
				panic(err)
			}
		}
		if writeToFile {
			err2 := ioutil.WriteFile(path, []byte(content), 0644)
			if err2 != nil {
				panic("write file problem")
			}
		}
	}
	return nil
}

func main() {
	// need to print date, find string, rep string, and root dir, extension

	var pWalker = func(pathX string, infoX os.FileInfo, errX error) error {
		if errX != nil {
			fmt.Printf("error 「%v」 at a path 「%q」\n", errX, pathX)
			return errX
		}

		if infoX.IsDir() {
			if !pass(filepath.Base(pathX), dirsToSkip) {
				return filepath.SkipDir
			}
		} else {
			if filepath.Ext(pathX) == ext {
				doFile(pathX)
			}
		}

		return nil
	}

	err := filepath.Walk(inDir, pWalker)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", inDir, err)
	}

	fmt.Printf("%v\n", "Done.")
}
