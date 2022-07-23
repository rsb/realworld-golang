package main

import "github.com/rsb/realworld-golang/app/cli/conduit/cmd"

var build = "develop"

func main() {
	cmd.Execute(build)
}
