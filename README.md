# go-todo-check
`go-todo-check` finds TODO comments that do not contain issue links, and if the TODO comment contains a github issue or asana task, it finds one that has already been completed.

## Install

```sh
go get github.com/mikito39/go-todo-check/cmd/go-todo-check
```

## Requirement

`go-todo-check` requires api access token

```sh
export ASANA_PERSONAL_ACCESS_TOKEN="YOUR_ASANA_PERSONAL_ACCESS_TOKEN"
export GITHUB_ACCESS_TOKEN="YOUR_GITHUB_ACCESS_TOKEN"
```

## Usage

```sh
go vet -vettool=`which go-todo-check` [flag] pkgname
Flags:
      -issue s must contain github issue link (ex. https://github.com/hogehoge/hogehoge)
      -workspaces s must contain asana workspaces link (ex. https://app.asana.com/0/hogehoge,https://app.asana.com/0/hoge1hoge1)

if you put "-ignore" on TODO comment, it ignores whether the given link has already been completed or not.
```
## Example
```sh
go vet -vettool=$(which go-todo-check) ./...
go vet -vettool=$(which go-todo-check) -issue https://github.com/hogehoge/hogehoge -workspaces https://app.asana.com/0/hogehoge,https://app.asana.com/0/hoge1hoge1 ./...
```

## Reference
- https://github.com/MakotoNaruse/todocomment
- https://github.com/google/go-github
- https://github.com/range-labs/go-asana

# License

`go-todo-check` is under [MIT license](https://en.wikipedia.org/wiki/MIT_License).