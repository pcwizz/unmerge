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
	"bytes"
	"strings"
	"testing"
)

type test struct {
	Title           string
	BeginColumn     uint
	EndColumn       uint
	InputDelimiter  string
	OutputDelimiter string
	Verbose         bool
	Input           string
	Output          string
}

var TESTS = []test{
	{
		Title:           "Single column column defined delimiter \"\"",
		BeginColumn:     0,
		EndColumn:       1,
		InputDelimiter:  "",
		OutputDelimiter: "",
		Verbose:         false,
		Input: `Beer



Vodka
Vodka
Beer


Whisky




Vodka

Vodka


Beer`,
		Output: `Beer
Beer
Beer
Beer
Vodka
Vodka
Beer
Beer
Beer
Whisky
Whisky
Whisky
Whisky
Whisky
Vodka
Vodka
Vodka
Vodka
Vodka
Beer`,
	},
	{
		Title:           "Double column column defined delimiter \"\\s*:\\s*\"",
		BeginColumn:     0,
		EndColumn:       2,
		InputDelimiter:  `\s*:\s*`,
		OutputDelimiter: "\t:\t",
		Verbose:         false,
		Input: `Beer	:	Beer
	:	
	:	
	:	
Vodka	:	Spirit
Vodka	:	
Beer	:	Beer
	:	
	:	
Whisky	:	Spirit
	:	
	: 
	:	
	: 
Vodka	:
	:	
Vodka	: 
	:	
	: 
Beer	:	Beer`,
		Output: `Beer	:	Beer
Beer	:	Beer
Beer	:	Beer
Beer	:	Beer
Vodka	:	Spirit
Vodka	:	Spirit
Beer	:	Beer
Beer	:	Beer
Beer	:	Beer
Whisky	:	Spirit
Whisky	:	Spirit
Whisky	:	Spirit
Whisky	:	Spirit
Whisky	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Beer	:	Beer`,
	},
	{
		Title:           "Double column column defined delimiter \"\\s*:\\s*\" undefined starts and ends",
		BeginColumn:     0,
		EndColumn:       0,
		InputDelimiter:  `\s*:\s*`,
		OutputDelimiter: "\t:\t",
		Verbose:         false,
		Input: `Beer	:	Beer
	:	
	:	
	:	
Vodka	:	Spirit
Vodka	:	
Beer	:	Beer
	:	
	:	
Whisky	:	Spirit
	:	
	: 
	:	
	: 
Vodka	:
	:	
Vodka	: 
	:	
	: 
Beer	:	Beer`,
		Output: `Beer	:	Beer
Beer	:	Beer
Beer	:	Beer
Beer	:	Beer
Vodka	:	Spirit
Vodka	:	Spirit
Beer	:	Beer
Beer	:	Beer
Beer	:	Beer
Whisky	:	Spirit
Whisky	:	Spirit
Whisky	:	Spirit
Whisky	:	Spirit
Whisky	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Vodka	:	Spirit
Beer	:	Beer`,
	},
	{
		Title:           "Double column column defined delimiter \"\\s*:\\s*\" defined start =1 and end = 2",
		BeginColumn:     1,
		EndColumn:       2,
		InputDelimiter:  `\s*:\s*`,
		OutputDelimiter: "\t:\t",
		Verbose:         false,
		Input: `Beer	:	Beer
	:	
	:	
	:	
Vodka	:	Spirit
Vodka	:	
Beer	:	Beer
	:	
	:	
Whisky	:	Spirit
	:	
	: 
	:	
	: 
Vodka	:
	:	
Vodka	: 
	:	
	: 
Beer	:	Beer`,
		Output: `Beer
Beer
Beer
Beer
Spirit
Spirit
Beer
Beer
Beer
Spirit
Spirit
Spirit
Spirit
Spirit
Spirit
Spirit
Spirit
Spirit
Spirit
Beer`,
	},
	{
		Title:           "Double column column defined delimiter \"\\s*:\\s*\" defined start =0 and end = 1",
		BeginColumn:     0,
		EndColumn:       1,
		InputDelimiter:  `\s*:\s*`,
		OutputDelimiter: "\t:\t",
		Verbose:         false,
		Input: `Beer	:	Beer
	:	
	:	
	:	
Vodka	:	Spirit
Vodka	:	
Beer	:	Beer
	:	
	:	
Whisky	:	Spirit
	:	
	: 
	:	
	: 
Vodka	:
	:	
Vodka	: 
	:	
	: 
Beer	:	Beer`,
		Output: `Beer
Beer
Beer
Beer
Vodka
Vodka
Beer
Beer
Beer
Whisky
Whisky
Whisky
Whisky
Whisky
Vodka
Vodka
Vodka
Vodka
Vodka
Beer`,
	},
}

func TestWorker(T *testing.T) {
	for _, test := range TESTS {
		ibuf := bufio.NewReader(strings.NewReader(test.Input))
		obuf := bytes.NewBufferString("")
		err := worker(ibuf, test.BeginColumn, test.EndColumn, test.InputDelimiter, test.OutputDelimiter, test.Verbose, obuf)
		if err != nil {
			T.Error(err)
		}
		out := strings.TrimSpace(obuf.String())
		if test.Output != out {
			T.Errorf("Faild test %s\nInput\n%s\nExpected\n%s\nOutput\n%s\n", test.Title, test.Input, test.Output, out)
		}
	}

}

func TestValidateColumnMarkers001(T *testing.T) {
	err := validateColumnMarkers(0, 0)
	if err != nil {
		T.Fatal("Start & finish maker both 0 was not allowed")
	}
}

func TestValidateColumnMarkers002(T *testing.T) {
	err := validateColumnMarkers(1, 0)
	if err == nil {
		T.Fatal("Start marker shouldn't be alowed to be larger than finish marker")
	}
}

func TestValidateColumnMarkers003(T *testing.T) {
	err := validateColumnMarkers(5, 10)
	if err != nil {
		T.Fatal("The start value is allowed to be less than the end value")
	}
}
