/*
	Unmerge is a program for 'dragging down' cells in tables.
	Copyright (C) 2015  Morgan Hill <morgan@pcwizzltd.com>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Parse arguments
	verbose := flag.Bool("v", false, "Print verbose output to stdout")
	startColumn := flag.Uint("b", 0, "Fist column from the left of the table to unmerge")
	endColumn := flag.Uint("e", 0, "Last column from left of table to unmerge")
	inputDelimiter := flag.String("id", `\s*:\s*`, "Feild input delimiter (Golang regexp)")
	outputDelimiter := flag.String("od", "\t:\t", "Feild output delimiter")
	flag.Parse()
	// Open stdin
	ibuf := bufio.NewReader(os.Stdin)
	err := worker(ibuf, *startColumn, *endColumn, *inputDelimiter, *outputDelimiter, *verbose, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

}

func worker(ibuf *bufio.Reader, startColumn uint, endColumn uint, inputDelimiter string, outputDelimiter string, verbose bool, r io.Writer) error {
	err := validateColumnMarkers(startColumn, endColumn)
	if err != nil {
		return err
	}
	if inputDelimiter == "" {
		inputDelimiter = `\s*:\s*`
	}
	if outputDelimiter == "" {
		outputDelimiter = "\t:\t"
	}
	// Initiate array to store the last non-null value of each column
	var lastLine []string
	if endColumn == 0 {
		lastLine = make([]string, 0, 25)
	} else {
		lastLine = make([]string, endColumn-startColumn)
	}
	delim := regexp.MustCompile(inputDelimiter)
	EOF := false
	for !EOF {
		if verbose {
			fmt.Fprintf(r, "Last Line %s\n", lastLine)
		}
		line, err := ibuf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if verbose {
					fmt.Fprintln(r, "Reached end of file")
				}
				EOF = true
			} else {
				return err
			}
		}
		if verbose {
			fmt.Fprintf(r, "Line %s\n", line)
		}
		line = strings.TrimSpace(line)
		fields := delim.Split(line, -1)
		if verbose {
			fmt.Fprintf(r, "fields = %s\n", fields)
		}
		if endColumn == 0 {
			endColumn = uint(len(fields))
			lastLine = lastLine[:endColumn]
		}
		if verbose {
			fmt.Fprintf(r, "endColumn = %v\n", fields)
		}
		for i := startColumn; i < endColumn && i < uint(len(fields)); i++ {
			lastLinePos := i - startColumn
			if match, err := regexp.MatchString(`\s*`, fields[i]); fields[i] != "" || !match || err != nil {
				if err != nil {
					return err
				}
				if verbose {
					fmt.Fprintf(r, "\nLast line[%v] set to %s\n", i, fields[i])
				}
				lastLine[lastLinePos] = fields[i]
			}

		}
		if err := outputTableLine(r, lastLine, outputDelimiter); err != nil {
			return err
		}

	}
	return nil

}

func validateColumnMarkers(startColumn uint, endColumn uint) error {
	if startColumn > endColumn {
		return errors.New("Start column can't be greater than end column")
	}
	return nil
}

func outputTableLine(r io.Writer, row []string, delim string) error {
	if len(row) < 1 {
		return errors.New("Row lentgh less than one")
	}
	var line string
	for i, d := range row {
		line += d
		if i < len(row)-1 {
			line += delim
		}
	}
	fmt.Fprintln(r, line)
	return nil
}
