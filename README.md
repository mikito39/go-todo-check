# go-todo-check
`go-todo-check` finds TODO comments that do not contain issue links, and if the TODO comment contains a github issue or asana task, it finds one that has already been completed.

## Install

```sh
go get github.com/mikito39/go-todo-check/cmd/go-todo-check
```

## Usage

```sh
go vet -vettool=`which go-todo-check` [flag] pkgname
Flags:
      -issue s must contain github issue link (ex. https://github.com/hogehoge/hogehoge)
      -workspaces s must contain asana workspaces link (ex. https://app.asana.com/0/hogehoge,https://app.asana.com/0/hoge1hoge1)
```
## Example
```sh
go vet -vettool=$(which go-todo-check) ./...
go vet -vettool=$(which go-todo-check) -issue https://github.com/hogehoge/hogehoge -workspaces https://app.asana.com/0/hogehoge,https://app.asana.com/0/hoge1hoge1 ./...
```
