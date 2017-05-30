package main

import "os"
import "fmt"
import "io"
import "io/ioutil"
import "path/filepath"
import "strings"

const _helpUsage = `./file_operations <operation>
<operation>:
	cp fileA fileB
	mv locationA/fileA locationB/fileB
	mkdir path
	ls path
	del file
	search string path`

//command ids
const (
	cmdCp = iota
	cmdMv
	cmdMkdir
	cmdLs
	cmdDel
	cmdSearch
)

//command structure
type command struct {
	cmd    int
	params []string
}

func main() {
	progargs := os.Args[1:]
	if len(progargs) < 1 {
		fmt.Println(_helpUsage)
		return
	}
	// fmt.Println("Arguments", progargs)
	cmd, success := parseArgs(progargs)
	if !success {
		fmt.Println("Error parsing the arguments")
	} else {
		// fmt.Println("cmd:", cmd.cmd, cmd.params)

		switch cmd.cmd {
		case cmdCp:
			ret := copyFile(cmd.params[0], cmd.params[1])
			if !ret {
				fmt.Println("Failed to copy", cmd.params[0], "to", cmd.params[1])
			}
		case cmdMv:
			ret := moveFile(cmd.params[0], cmd.params[1])
			if !ret {
				fmt.Println("Failed to move", cmd.params[0], "to", cmd.params[1])
			}
		case cmdMkdir:
			ret := makeDir(cmd.params[0])
			if !ret {
				fmt.Println("Failed to create dir", cmd.params[0])
			}
		case cmdDel:
			ret := deleteFile(cmd.params[0])
			if !ret {
				fmt.Println("Failed to delete ", cmd.params[0])
			}
		case cmdLs:
			ret := listFiles(cmd.params[0])
			if !ret {
				fmt.Println("Failed to list ", cmd.params[0])
			}
		case cmdSearch:
			s := new(searchObj)
			ret := s.searchFileName(cmd.params[0], cmd.params[1])
			if !ret {
				fmt.Println("Failed to search ", cmd.params[0], " in ", cmd.params[1])
			}
		default:
			fmt.Println("Command not supported")
		}
	}
}

func parseArgs(args []string) (result command, success bool) {
	cmd := command{}
	success = true
	switch args[0] {
	case "cp":
		cmd.cmd = cmdCp
		if len(args) < 3 {
			return cmd, false
		}
	case "mv":
		cmd.cmd = cmdMv
		if len(args) < 3 {
			return cmd, false
		}
	case "mkdir":
		cmd.cmd = cmdMkdir
	case "ls":
		cmd.cmd = cmdLs
	case "del":
		cmd.cmd = cmdDel
	case "search":
		if len(args) < 3 {
			return cmd, false
		}
		cmd.cmd = cmdSearch
	default:
		success = false
	}
	if len(args) > 1 {
		cmd.params = args[1:]
	}

	return cmd, success
}

func copyFile(src string, dest string) (success bool) {
	in, err := os.Open(src)
	if err != nil {
		fmt.Println("Source file ", src, "could not be opened")
		return false
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		fmt.Println("Destination file ", dest, "could not be created")
		return false
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Println("File could not be copied")
		return false
	}
	err = out.Close()
	if err != nil {
		fmt.Println("Destination file could not be closed")
		return false
	}

	return true
}

func moveFile(src string, dest string) (success bool) {
	err := os.Rename(src, dest)
	if err != nil {
		fmt.Println("File", src, "could not be moved")
		return false
	}
	return true
}

func makeDir(dest string) (success bool) {
	err := os.Mkdir(dest, os.FileMode(777))
	if err != nil {
		return false
	}
	return true
}

func deleteFile(src string) (success bool) {
	err := os.Remove(src)
	if err != nil {
		return false
	}
	return true
}

func listFiles(src string) (success bool) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return false
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}

	return true
}

type searchObj struct {
	strToSearch string
}

func (s *searchObj) visit(path string, f os.FileInfo, err error) error {
	// fmt.Printf("Visited: %s\n", path)
	if strings.Contains(f.Name(), s.strToSearch) {
		fmt.Println("found match: ", path)
	}
	return nil
}

func (s *searchObj) searchFileName(toSearch string, path string) (success bool) {
	s.strToSearch = toSearch
	filepath.Walk(path, s.visit)
	return true
}
