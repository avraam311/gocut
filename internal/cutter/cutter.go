package cutter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cutter struct {
}

func New() *Cutter {
	return &Cutter{}
}

func (c *Cutter) Cut(delim, flagF, filename string, flagS bool) {
	if flagF == "" {
		fmt.Println("flag -f is required")
		os.Exit(1)
	}
	columns := parseFields(flagF)
	scanner := getScanner(filename)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r\n")
		if line == "" {
			continue
		}
		fields := strings.Split(line, delim)
		outputColumns := c.extractColumns(fields, columns, flagS)
		if outputColumns != nil {
			fmt.Println(strings.Join(outputColumns, " "))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("error scanning: %v\n", err)
		os.Exit(1)
	}
}

func (c *Cutter) extractColumns(fields []string, columns []int, flagS bool) []string {
	if flagS && len(fields) <= 1 {
		return nil
	}
	output := make([]string, 0, len(columns))
	for _, col := range columns {
		if col > 0 && col <= len(fields) {
			output = append(output, fields[col-1])
		} else {
			output = append(output, "")
		}
	}
	return output
}

func parseFields(fieldStr string) []int {
	var fields []int
	for _, part := range strings.Split(fieldStr, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			parts := strings.Split(part, "-")
			if len(parts) != 2 {
				fmt.Println("invalid -f range format")
				os.Exit(1)
			}
			start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil || start < 1 {
				fmt.Println("-f start >= 1")
				os.Exit(1)
			}
			endStr := strings.TrimSpace(parts[1])
			end, err := strconv.Atoi(endStr)
			if err != nil || end < 1 || start > end {
				fmt.Println("-f end >= start >= 1")
				os.Exit(1)
			}
			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil || num < 1 {
				fmt.Println("-f positive integers only")
				os.Exit(1)
			}
			fields = append(fields, num)
		}
	}
	return fields
}

func getScanner(filename string) *bufio.Scanner {
	if filename == "" {
		return bufio.NewScanner(os.Stdin)
	}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("open %s: %v\n", filename, err)
		os.Exit(1)
	}
	return bufio.NewScanner(f)
}
