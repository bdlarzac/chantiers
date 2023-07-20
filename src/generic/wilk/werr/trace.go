// Copyright (c) 2023 William Dode. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package werr

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Allow trace or disable for the project
var WithTrace = true

// Show full path instead of basename
var WithFullPath = false

// Show full name instead of [1:]
var WithFullName = false

// Wrapf returns formated error with trace and : %w
func Wrapf(err error, s string, vals ...any) error {
	if err == nil {
		return nil
	}
	vals = append(vals, err)
	return tracef(2, s+": %w", vals...)
}

// Wrap returns error with trace and %w
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return tracef(2, "%w", err)
}

// Errorf is like fmt.Errorf with trace
// wrap if %w
// not wrap if %v
func Errorf(s string, vals ...any) error {
	return tracef(2, s, vals...)
}

// New is like errors.New with trace
func New(s string) error {
	return tracef(2, s)
}

// tracef add trace before calling fmt.Errorf
func tracef(skip int, s string, vals ...any) error {
	pc, file, line, ok := runtime.Caller(skip)
	if ok && WithTrace {
		name := runtime.FuncForPC(pc).Name()
		if !WithFullName {
			splt := strings.Split(runtime.FuncForPC(pc).Name(), "/")
			name = strings.Join(splt[1:], "/")
		}
		path := file
		if !WithFullPath {
			path = filepath.Base(file)
		}
		info := fmt.Sprintf("\n> %s() %s:%d\n",
			name, path, line)
		s = info + strings.TrimSpace(s)
	}
	return fmt.Errorf(s, vals...)
}

// ====== additions Thierry Graff ======

func Print(err error) string {
    return err.Error()
}

// SprintHTML returns traceback as string formatted for html display
func SprintHTML(err error) string {
	res := ""
	res += `
<style>
    table.werr{
        border-collapse:collapse;
        margin:.5rem;
    }
    table.werr td{
        padding:2px 4px;
        border:1px solid #a2a9b1;
    }
    table.werr tr:nth-child(even){
        background-color:#f2e6d9;
    }
    table.werr tr:nth-child(odd){
        background-color:lightblue;
    }
    div.werr{
        font-weight:bold;
    }
</style>
    `
	res += "<table class=\"werr\">\n"
	lines := strings.Split(err.Error(), "\n")
	// Each error takes 2 lines
	// - first line contains '> ' followed by function, filename and line number
	// - second line contains the error message
	msg := "<div class=\"werr\">"
	for i, line := range lines {
	    if line == "" {
	        continue
	    }
	    if i%2 == 1 {
			res += "    <tr>\n"
			line, _ = strings.CutPrefix(line, "> ")
			parts := strings.Split(line, " ")
			res += "        <td>" + parts[0] + "</td>\n"
			res += "        <td>" + parts[1] + "</td>\n"
	    } else {
	        line = strings.TrimSpace(line)
			line, _ = strings.CutSuffix(line, ":")
			res += "        <td>" + line + "</td>\n"
			res += "    </tr>\n"
	    }
	}
	res += "</table>\n"
	msg += "</div>\n"
	res += msg
	return res
}

// ====== end additions Thierry Graff ======

