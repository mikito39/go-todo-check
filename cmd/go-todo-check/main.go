package main

import (
	"github.com/mikito39/go-todo-check"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(go_todo_check.Analyzer) }
