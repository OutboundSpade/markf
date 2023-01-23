package macro

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/OutboundSpade/markf/logger"
)

var VariableStore = map[string]string{}

var whitespace_sel = regexp.MustCompile(`[\ \t]+`)
var newln_sel = regexp.MustCompile(`\n`)
var list_del = regexp.MustCompile(`\|`)
var DefaultMacros = []*Macro{
	{
		MacroName: "var",
		MacroFunc: func(args *[]string) string {
			if len(*args) == 2 && VariableStore[(*args)[1]] != "" {
				return VariableStore[(*args)[1]]
			} else if len(*args) >= 3 && (*args)[2] == "=" {
				VariableStore[(*args)[1]] = strings.Join((*args)[3:], " ")
				return ""
			} else if VariableStore[(*args)[1]] == "" {
				return fmt.Sprintf("var %s is not defined", (*args)[1])
			}
			return "Usage: var \\<varname> = \\<value>"
		},
	},
	{
		MacroName: "list",
		MacroFunc: func(args *[]string) string {
			if len(*args) <= 1 {
				return "Usage: list \\<items...>"
			}
			if len(newln_sel.Split(strings.Join((*args)[1:], " "), -1)) > 1 {
				return strings.Join(newln_sel.Split(strings.Join((*args)[1:], " "), -1), "|")
			}
			return strings.Join(whitespace_sel.Split(strings.Join((*args)[1:], " "), -1), "|")
		},
	},
	{
		MacroName: "trim",
		MacroFunc: func(args *[]string) string {
			if len(*args) <= 1 {
				return "Usage: trim \\<from> \\<to> \\<list>"
			}
			from, err := strconv.Atoi((*args)[1])
			if err != nil {
				return fmt.Sprintf("invalid '\\<from>': %s", (*args)[1])
			}
			to, err := strconv.Atoi((*args)[2])
			if err != nil {
				return fmt.Sprintf("invalid '\\<to>': %s", (*args)[2])
			}
			list := strings.Split((*args)[3], "|")
			if from < 0 {
				from = 0
			}
			if to < 0 {
				to = len(list)
			}
			if from > len(list) {
				from = len(list) - 1
			}
			if to > len(list) {
				to = len(list)
			}

			var output []string
			for i, item := range list {
				if i < from || i >= to {
					continue
				}
				output = append(output, item)
			}
			return strings.Join(output, "|")
		},
	},
	{
		MacroName: "foreach",
		MacroFunc: func(args *[]string) string {
			if len(*args) < 5 {
				return "Usage: foreach \\<varname> in \\<list> \\<body>"
			}
			varname := (*args)[1]
			list := (*args)[3]
			body := (*args)[4]
			var output string
			for _, item := range list_del.Split(list, -1) {
				if item == "" {
					continue
				}
				VariableStore[varname] = item
				output += fmt.Sprintf("#!(var %s = %s)", varname, item)
				// fmt.Println(item)
				output += body
			}
			logger.Printf("foreach out: %s\n", output)
			return output
		},
	},
	{
		MacroName: "if",
		MacroFunc: func(args *[]string) string {
			if len(*args) < 5 {
				return "Usage: if <data> \\<condition> <data> \\<body> else \\<body>"
			}
			numberify := func(str string) (int, bool) {
				num, err := strconv.Atoi(str)
				if err != nil {
					return 0, false
				}
				return num, true
			}
			data1 := (*args)[1]
			condition := (*args)[2]
			data2 := (*args)[3]
			body := (*args)[4]
			isTrue := false
			num1, num1Valid := numberify(data1)
			num2, num2Valid := numberify(data2)
			requiresValidNums := false
			switch condition {
			case "=", "==":
				if data1 == data2 {
					isTrue = true
				}
			case "!=":
				if data1 != data2 {
					isTrue = true
				}
			case ">":
				requiresValidNums = true
				if num1 > num2 {
					isTrue = true
				}
			case "<":
				requiresValidNums = true
				if num1 < num2 {
					isTrue = true
				}
			case ">=":
				requiresValidNums = true
				if num1 >= num2 {
					isTrue = true
				}
			case "<=":
				requiresValidNums = true
				if num1 <= num2 {
					isTrue = true
				}
			default:
				return fmt.Sprintf("invalid condition: %s", condition)
			}
			if requiresValidNums {
				if !num1Valid {
					return fmt.Sprintf("invalid number: %s", data1)
				}
				if !num2Valid {
					return fmt.Sprintf("invalid number: %s", data2)
				}
			}
			if isTrue {
				return body
			} else if len(*args) > 5 && (*args)[5] == "else" {
				return (*args)[6]
			}
			return ""
		},
	},
}
