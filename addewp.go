/*The command line utility adds a file into IAR's EWP file

Interactive usage:

	$ ./addewp.exe

Quick usage with flags:

	$ ./addewp.exe -ewp MyEwpFile.ewp -file main.cpp

Flags may be omitted for interactivity. For example, the following invocation
prompts for an EWP file, since the new file is provided as an argument:

	$ ./addewp.exe -file main.cpp
*/
package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/IQ-Inc/iarewp"
)

const (
	defaultPath string = ""
)

var (
	ewpFile = flag.String("ewp", defaultPath, "the EWP file")
	newFile = flag.String("file", defaultPath,
		"the new file to include in the project")
)

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	flag.Parse()

	stdin := bufio.NewReader(os.Stdin)

	// Handle default EWP / interactive EWP
	if *ewpFile == defaultPath {
		fmt.Print("Specify EWP file path: ")
		path, err := stdin.ReadString('\n')

		if err != nil {
			fail("Error obtaining EWP path")
		}

		*ewpFile = strings.TrimSpace(path)
	}

	// Check for EWP file modes
	ewpinfo, err := os.Stat(*ewpFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fail("Error locating EWP file")
	}

	if int(ewpinfo.Mode().Perm())&os.O_RDWR == 0 {
		fail("Error: EWP file is not read/write")
	}

	// Handle default new file / interactive new file
	if *newFile == defaultPath {
		fmt.Print("Specify new file to include in EWP: ")

		path, err := stdin.ReadString('\n')

		if err != nil {
			fail("Error obtaining file path")
		}

		*newFile = strings.TrimSpace(path)
	}

	// Check if new file exists
	if _, err := os.Stat(*newFile); os.IsNotExist(err) {
		fmt.Print("Warning: file ", *newFile, " not found. Continue to add? (y): ")
		answer, err := stdin.ReadString('\n')

		if err != nil {
			fail("Error obtaining user confirmation")
		} else if ans := strings.TrimSpace(answer); ans != "y" && ans != "Y" {
			fmt.Println("Not adding", *newFile)
			return
		}
	}

	// Unmarshal EWP and insert file
	ewpf, err := ioutil.ReadFile(*ewpFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fail("Error opening ewp file: " + *ewpFile)
	}

	var ewp iarewp.Ewp
	if err := xml.Unmarshal(ewpf, &ewp); err != nil {
		fail("Error unmarshalling EWP file")
	}

	f := iarewp.MakeFile(*newFile)

	if ewp.Contains(f) {
		fmt.Printf("%s already contains %s. Not adding %s.", *ewpFile, *newFile, *newFile)
		return
	}

	ewp.InsertFile(f)

	bs, err := xml.MarshalIndent(ewp, "", "    ")
	if err != nil {
		fail("Error reconstructing EWP. Original EWP is not modified.")
	}

	header := []byte(xml.Header)
	final := append(header, bs[:]...)

	ioutil.WriteFile(*ewpFile, final, 0644)
}
