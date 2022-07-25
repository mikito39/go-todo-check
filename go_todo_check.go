package go_todo_check

import (
	"flag"
	"github.com/mikito39/go-todo-check/asana"
	"github.com/mikito39/go-todo-check/github"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"regexp"
	"strings"
)

const doc = "go-todo-check is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "go_todo_check",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var issue string
var workspaces string

func init() {
	flag.StringVar(&issue, "issue", "github", "Master Issue Link ex. https://github.com/test/test/issues")
	flag.StringVar(&workspaces, "workspaces", "asana", "Master Workspace Links ex. https://app.asana.com/0/test")
}

func parseIssueLink(s string) (org, repo, id string) {
	r := regexp.MustCompile(`https://github.com/[^\s]+`)
	val := strings.Split(r.FindString(s), "/")
	return val[3], val[4], val[6]
}

func parseTaskLink(s string) (taskListID, taskID string) {
	r := regexp.MustCompile(`https://app.asana.com/[^\s]+`)
	val := strings.Split(r.FindString(s), "/")
	return val[4], val[5]
}

func checkWorkspaceContains(s, workspaces string) bool {
	val := strings.Split(workspaces, ",")
	for _, i := range val {
		if strings.Contains(s, i) {
			return true
		}
	}
	return false
}

// TODO: リファクタリングする。
func run(pass *analysis.Pass) (interface{}, error) {

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			for _, comment := range n.Comments {
				s := comment.Text()
				if strings.Contains(s, "nolint: go-todo-check") {
					continue
				}
				if strings.Contains(s, "TODO ") || strings.Contains(s, "TODO:") {
					if strings.Contains(s, issue) {
						org, repo, issueID := parseIssueLink(s)
						if !github.IssueStatus(org, repo, issueID) {
							pass.Reportf(comment.Pos(), "TODO comment must contains open github issue's link")
						}
					} else if checkWorkspaceContains(s, workspaces) {
						taskListID, taskID := parseTaskLink(s)
						if !asana.TaskStatus(taskID, taskListID) {
							pass.Reportf(comment.Pos(), "TODO comment must contains incomplete asana task's link")
						}
					} else {
						pass.Reportf(comment.Pos(), "TODO comment must contains open github issue's link or asana task's link")
					}
				}
			}
		}
	})

	return nil, nil
}
