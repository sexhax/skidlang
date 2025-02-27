# Skidlang Compiler

Skidlang is a simple programming language that compiles to Go (Golang). The language provides basic features such as variable declarations, conditionals, loops, functions, and more. This repository provides the compiler for Skidlang that reads `.skid` files, tokenizes the code, and compiles it into an executable Go binary.

## Features

- **Lexical analysis (tokenization)**: The Skidlang code is parsed into tokens.
- **Syntax processing**: The code is transformed into Go code with equivalent functionality.
- **Compile to Go**: Skidlang code is compiled into Go code and can be built into an executable.
- **Basic constructs**: Includes variable declarations, conditionals, loops, functions, error handling, and more.

## How to Use

### Installation

To use the Skidlang compiler, you need Go installed on your system. If you haven't installed Go yet, download and install it from [here](https://golang.org/dl/).

### Usage

1. **Prepare your Skidlang code**: Write your Skidlang code in a file with the `.skid` extension. Example:

   ```skid
   let x int 10
   let y int 20
   if x < y
       print "x is less than y"
   end
   ```
   
## Compile and Build: To compile your Skidlang code into an executable Go binary, run the following command:

```
bash

go run main.go <filename>.skid
```

Replace <filename> with the name of your Skidlang file. The compiler will generate a .go file which will be deleted after builing and then build it into an executable for your system.

Example:

```
bash

go run skidlang.go example.skid
```
Run the Executable: After successful compilation, you can run the generated executable file:

```
bash

./hello      # or hello.exe on Windows
```
# Syntax
Skidlang follows a basic syntax for most constructs. Here's a breakdown of the major elements.

## 1. Variables
let: Declare a variable with an optional type.
```
let <varName> <type> <value>   # Example: let x int 10
let <varName> <value>          # Example: let x 10
```

set: Assign a value to an already declared variable.

```
set <varName> <value>
```
## 2. Conditionals
if: Start a conditional block.

```
if <condition>
    <commands>
end
```
else: Define the else block.

```
else
    <commands>
end
```
else if: Not directly supported; use multiple if statements.

## 3. Loops
while: Start a while loop.

```
while <condition>
    <commands>
end
```
for: Start a for loop.

```
for <initialization> <condition> <increment>
    <commands>
end
```
## 4. Functions
func: Define a function.

```
func <functionName> (<params>)
    <commands>
endfunc
```
return: Return a value from a function.

```
return <value>
```
call: Call a function.

```
call <functionName> <args>
```
## 5. Error Handling
try: Start a try block.

```
try
    <commands>
catch <errorVariable>
    <commands>
endtry
```
## 6. Other Constructs
print: Print to the console.

```
print "<message>"
```

input: Read input from the user.

```
input <variableName> <type>
```

switch / case: Start a switch block.

```
switch <expression>
case <value>
    <commands>
default
    <commands>
endswitch
```
struct: Define a structure.

```
struct <structName>
    field <fieldName> <fieldType>
endstruct
```

# Example Skidlang Program

```
// A simple Skidlang program

let x int 10
let y int 20

if x < y
    print "x is less than y"
else
    print "x is greater than or equal to y"
end
```
This Skidlang program declares two variables x and y, compares them, and prints a message based on the result of the comparison.

## How it Works
Lexical Analysis: The Lex function tokenizes the input Skidlang code into commands and arguments.
Compilation: The Compile function transforms the tokens into Go code.
Build: The BuildExecutable function compiles the Go code into a binary executable.


## License
This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.
