package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Token struct {
	Command string
	Args    []string
}

func tokenizeLine(line string) []string {
	var tokens []string
	currentToken := ""
	var inQuote rune = 0
	inBracket := 0
	escape := false

	for _, char := range line {
		if escape {
			currentToken += string(char)
			escape = false
			continue
		}

		switch char {
		case '\\':
			escape = true
		case '\'', '"':
			if inQuote == 0 {
				inQuote = char
			} else if inQuote == char {
				inQuote = 0
			}
			currentToken += string(char)
		case '[':
			if inQuote == 0 {
				inBracket++
			}
			currentToken += string(char)
		case ']':
			if inQuote == 0 && inBracket > 0 {
				inBracket--
			}
			currentToken += string(char)
		case ' ', '\t':
			if inQuote == 0 && inBracket == 0 {
				if currentToken != "" {
					tokens = append(tokens, currentToken)
					currentToken = ""
				}
			} else {
				currentToken += string(char)
			}
		default:
			currentToken += string(char)
		}
	}

	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	return tokens
}

func Lex(input string) []Token {
	lines := strings.Split(input, "\n")
	tokens := []Token{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := tokenizeLine(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]
		tokens = append(tokens, Token{Command: command, Args: args})
	}

	return tokens
}

func isType(s string) bool {
	switch s {
	case "int", "float64", "string", "bool", "rune", "byte", "[]int", "[]string", "[]float64":
		return true
	default:
		return false
	}
}

func Compile(tokens []Token) string {
	imports := map[string]bool{"fmt": true}
	output := []string{"package main", "", "import ("}
	additionalCode := []string{}

	currentStruct := ""

	for _, token := range tokens {
		switch token.Command {
		case "print":
			args := token.Args
			output = append(output, fmt.Sprintf("\tfmt.Println(%s)", strings.Join(args, " ")))
		case "printf":
			if len(token.Args) < 1 {
				continue
			}
			format := token.Args[0]
			args := token.Args[1:]
			output = append(output, fmt.Sprintf("\tfmt.Printf(%s, %s)", format, strings.Join(args, ", ")))
		case "let":
			if len(token.Args) < 2 {
				continue
			}
			varName := token.Args[0]
			if isType(token.Args[1]) {
				if len(token.Args) < 3 {
					continue
				}
				typeName := token.Args[1]
				value := strings.Join(token.Args[2:], " ")
				output = append(output, fmt.Sprintf("\tvar %s %s = %s", varName, typeName, value))
			} else {
				value := strings.Join(token.Args[1:], " ")
				output = append(output, fmt.Sprintf("\t%s := %s", varName, value))
			}
		case "set":
			if len(token.Args) < 2 {
				continue
			}
			varName := token.Args[0]
			value := strings.Join(token.Args[1:], " ")
			output = append(output, fmt.Sprintf("\t%s = %s", varName, value))
		case "if":
			condition := strings.Join(token.Args, " ")
			output = append(output, fmt.Sprintf("\tif %s {", condition))
		case "else":
			output = append(output, "\t} else {")
		case "end":
			output = append(output, "\t}")
		case "while":
			condition := strings.Join(token.Args, " ")
			output = append(output, fmt.Sprintf("\tfor %s {", condition))
		case "for":
			parts := strings.Join(token.Args, " ")
			output = append(output, fmt.Sprintf("\tfor %s {", parts))
		case "func":
			if len(token.Args) < 1 {
				continue
			}
            funcParts := token.Args
            receiver := ""
            funcName := ""
            returnType := ""

			if strings.HasPrefix(funcParts[0], "(") {
				receiver = funcParts[0]
				funcName = funcParts[1]
				funcParts = funcParts[2:]
			} else {
				funcName = funcParts[0]
				funcParts = funcParts[1:]
			}

			for i, part := range funcParts {
				if part == "->" && i+1 < len(funcParts) {
					returnType = funcParts[i+1]
					funcParts = funcParts[:i]
					break
				}
			}

			params := []string{}
			i := 0
			for i < len(funcParts) {
				if i+1 >= len(funcParts) {
					break
				}
				paramName := funcParts[i]
				paramType := funcParts[i+1]
				params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
				i += 2
			}
			paramStr := strings.Join(params, ", ")

			funcDef := ""
			if receiver != "" {
				funcDef = fmt.Sprintf("func %s %s(%s)", receiver, funcName, paramStr)
			} else {
				funcDef = fmt.Sprintf("func %s(%s)", funcName, paramStr)
			}

			if returnType != "" {
				funcDef += " " + returnType
			}

			additionalCode = append(additionalCode, fmt.Sprintf("%s {", funcDef))
		case "endfunc":
			additionalCode = append(additionalCode, "}")
		case "call":
			if len(token.Args) < 1 {
				continue
			}
			funcName := token.Args[0]
			args := token.Args[1:]
			output = append(output, fmt.Sprintf("\t%s(%s)", funcName, strings.Join(args, ", ")))
		case "return":
			if len(token.Args) == 0 {
				additionalCode = append(additionalCode, "\treturn")
			} else {
				additionalCode = append(additionalCode, fmt.Sprintf("\treturn %s", strings.Join(token.Args, " ")))
			}
		case "input":
			if len(token.Args) == 2 && isType(token.Args[1]) {
				varName := token.Args[0]
				typeName := token.Args[1]
				output = append(output, fmt.Sprintf("\tvar %s %s", varName, typeName))
				output = append(output, fmt.Sprintf("\tfmt.Scan(&%s)", varName))
			} else if len(token.Args) == 1 {
				varName := token.Args[0]
				output = append(output, fmt.Sprintf("\tfmt.Scan(&%s)", varName))
			}
		case "switch":
			expr := strings.Join(token.Args, " ")
			output = append(output, fmt.Sprintf("\tswitch %s {", expr))
		case "case":
			expr := strings.Join(token.Args, " ")
			output = append(output, fmt.Sprintf("\tcase %s:", expr))
		case "default":
			output = append(output, "\tdefault:")
		case "endswitch":
			output = append(output, "\t}")
		case "try":
			output = append(output, "\tfunc() {")
			output = append(output, "\t\tdefer func() {")
			output = append(output, "\t\t\tif err := recover(); err != nil {")
		case "catch":
			if len(token.Args) >= 1 {
				errVar := token.Args[0]
				output = append(output, fmt.Sprintf("\t\t\t\t%s := err", errVar))
			}
			output = append(output, "\t\t\t\t// Handle error")
		case "endtry":
			output = append(output, "\t\t\t}")
			output = append(output, "\t\t}()")
		case "struct":
			if len(token.Args) < 1 {
				continue
			}
			currentStruct = token.Args[0]
			additionalCode = append(additionalCode, fmt.Sprintf("type %s struct {", currentStruct))
		case "field":
			if currentStruct == "" || len(token.Args) < 2 {
				continue
			}
			fieldName := token.Args[0]
			fieldType := strings.Join(token.Args[1:], " ")
			additionalCode = append(additionalCode, fmt.Sprintf("\t%s %s", fieldName, fieldType))
		case "endstruct":
			if currentStruct != "" {
				additionalCode = append(additionalCode, "}")
				currentStruct = ""
			}
		case "inc":
			if len(token.Args) == 1 {
				output = append(output, fmt.Sprintf("\t%s++", token.Args[0]))
			}
		case "dec":
			if len(token.Args) == 1 {
				output = append(output, fmt.Sprintf("\t%s--", token.Args[0]))
			}
		case "const":
			if len(token.Args) < 2 {
				continue
			}
			constName := token.Args[0]
			if isType(token.Args[1]) {
				if len(token.Args) < 3 {
					continue
				}
				typeName := token.Args[1]
				value := strings.Join(token.Args[2:], " ")
				additionalCode = append(additionalCode, fmt.Sprintf("const %s %s = %s", constName, typeName, value))
			} else {
				value := strings.Join(token.Args[1:], " ")
				additionalCode = append(additionalCode, fmt.Sprintf("const %s = %s", constName, value))
			}
		case "import":
			if len(token.Args) >= 1 {
				pkg := strings.Join(token.Args, " ")
				imports[pkg] = true
			}
		}
	}

	importList := []string{}
	for pkg := range imports {
		importList = append(importList, fmt.Sprintf("\t\"%s\"", pkg))
	}
	output[2] = strings.Join(importList, "\n")
	fullOutput := []string{output[0], output[1], output[2], output[3]}
	fullOutput = append(fullOutput, additionalCode...)
	fullOutput = append(fullOutput, "func main() {")
	fullOutput = append(fullOutput, output[4:]...)
	fullOutput = append(fullOutput, "}")

	return strings.Join(fullOutput, "\n")
}

func BuildExecutable(filename string, goCode string) error {
	goFile := strings.TrimSuffix(filename, ".skid") + ".go"
	err := ioutil.WriteFile(goFile, []byte(goCode), 0644)
	if err != nil {
		return fmt.Errorf("error writing intermediate Go file: %v", err)
	}

	outputBinary := strings.TrimSuffix(filename, ".skid")
	if runtime.GOOS == "windows" {
		outputBinary += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", outputBinary, goFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error building binary: %v", err)
	}

	err = os.Remove(goFile)
	if err != nil {
		return fmt.Errorf("error removing intermediate Go file: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <compiler> <filename>.skid")
		return
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".skid") {
		fmt.Println("Error: Input file must have .skid extension")
		return
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokens := Lex(string(content))
	goCode := Compile(tokens)

	err = BuildExecutable(filename, goCode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	outputBinary := strings.TrimSuffix(filename, ".skid")
	if runtime.GOOS == "windows" {
		outputBinary += ".exe"
	}

	fmt.Printf("Built successfully: %s\n", outputBinary)
}