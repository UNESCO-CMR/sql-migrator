package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println("USAGE: extract all tables:")
		fmt.Println("  go run main.go DUMP_FILE")
		fmt.Println("extract one table:")
		fmt.Println("  go run main.go DUMP_FILE [TABLE]")
		return
	}

	dumpFile := os.Args[1]
	start := "/-- Table structure for table/"
	// or
	// start := "/DROP TABLE IF EXISTS/"

	if len(os.Args) >= 3 {
		// extract one table os.Args[2]
		cmd := exec.Command("csplit", "-s", "-ftable", dumpFile, start, fmt.Sprintf("-- Table structure for table `%s`", os.Args[2]), "-- Table structure for table/", "%40103 SET TIME_ZONE=@OLD_TIME_ZONE%1")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else {
		// extract all tables
		cmd := exec.Command("csplit", "-s", "-ftable", dumpFile, start, "{*}")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	files, err := filepath.Glob("table*")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No table files found")
		return
	}

	headFile := "head"
	err = os.Rename("table00", headFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	lastFile := files[len(files)-1]
	var footFile string

	if len(os.Args) >= 3 {
		// cut off all other tables
		footFile = lastFile
	} else {
		// cut off the end of each file
		cmd := exec.Command("csplit", "-b", "%d", "-s", "-f"+lastFile, lastFile, "40103 SET TIME_ZONE=@OLD_TIME_ZONE/", "{*}")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		footFile = lastFile + "1"
	}

	for _, file := range files {
		name := getTableName(file)
		err := mergeFiles(headFile, file, footFile, name+".sql")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	err = os.Remove(headFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = os.Remove(footFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	fmt.Println("Extraction complete!")
}

func getTableName(file string) string {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	firstLine := scanner.Text()

	// Extract table name between backticks (`table_name`)
	name := strings.SplitN(firstLine, "`", 3)
	if len(name) < 3 {
		fmt.Println("Invalid table name:", firstLine)
		return ""
	}

	return name[1]
}

func mergeFiles(headFile, bodyFile, footFile, outputFile string) error {
	headData, err := ioutil.ReadFile(headFile)
	if err != nil {
		return err
	}

	bodyData, err := ioutil.ReadFile(bodyFile)
	if err != nil {
		return err
	}

	footData, err := ioutil.ReadFile(footFile)
	if err != nil {
		return err
	}

	outputData := append(append(headData, bodyData...), footData...)

	err = ioutil.WriteFile(outputFile, outputData, 0644)
	if err != nil {
		return err
	}

	return nil
}
