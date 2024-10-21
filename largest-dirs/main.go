package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

// Directory represents a directory with its size.
type Directory struct {
	Size int64
	Path string
}

// DirectoryList implements sort.Interface based on the Size field.
type DirectoryList []Directory

func (a DirectoryList) Len() int           { return len(a) }
func (a DirectoryList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DirectoryList) Less(i, j int) bool { return a[i].Size > a[j].Size }

func main() {

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usage: largest-dirs <root-dir> <limit>")
		return
	}

	rootDir := args[0]
	limit := args[1]

	// terminalCommand := fmt.Sprintf("sudo du -a %s | sort -n -r | head -n %s", initialDir, head)
	terminalCommand := fmt.Sprintf("sudo du -a %s", rootDir)
	commandInput := strings.Fields(terminalCommand)

	// Execute the `du` command to get directory sizes.
	cmd := exec.Command(commandInput[0], commandInput[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Println("Error executing du command:", err)
		return
	}

	// Parse the output.
	lines := strings.Split(out.String(), "\n")

	var directories DirectoryList

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		size, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing size:", err)
			return
		}

		path := fields[1]

		directories = append(directories, Directory{Size: size, Path: path})
	}

	sort.Sort(DirectoryList(directories))

	fmt.Println(fmt.Sprintf("Top %s largest directories:", limit))

	headConv, err := strconv.Atoi(limit)

	if err != nil {
		fmt.Println("Error converting head to int:", err)
		return
	}

	renderTable(headConv, directories)
}

func renderTable(head int, documents DirectoryList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Size", "Path"})
	for i, dir := range documents {
		if i >= head {
			break
		}

		t.AppendRow([]interface{}{i, formatSize(dir.Size), dir.Path})
	}
	t.SetAllowedColumnLengths([]int{5, 20, 50})
	t.SetStyle(table.StyleRounded)
	t.Render()
}

// formatSize converts a size in bytes to a human-readable string.
func formatSize(size int64) string {
	units := []string{"KB", "MB", "GB", "TB"}

	var i int

	for i = range units {
		if size < 1024 {
			break
		}

		size /= 1024
	}

	return fmt.Sprintf("%d %s", size, units[i])
}
