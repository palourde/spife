# Spife

*Spife* is a simple command line tool which helps bumping a cookbook by automatically updating the **metadata.rb** file and pushing the changes.

### Installation
#### Mac OS X
You can use the binary inside the *bin* directory

### From Source
- [Install Go](https://golang.org/doc/install).
- Run `go build spife.go`, which will generate the *spife* binary.

### Usage

```
# spife -h
Usage of spife:
  -bump-level="patch": Version level to bump the cookbook
  -git-push=true: Whether or not changes should be committed.
  -git-remotes="upstream, origin": Comma separated list of Git remotes
  -path="": Full or relative path to the cookbook directory. REQUIRED.
```

### Example

#### Quick example
`spife -path .`

#### Full example
`spife -bump-level major -git-push=true -git-remotes "upstream, origin" -path ~/cookbook/mycookbook`
