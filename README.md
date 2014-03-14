# gobuildhistory

A simple tool which goes through the git commit history of the repo cloned in the working directory and runs `go build` on each commit.
The motivation for this tool was to find the code which caused `go build` to fail with an internal compiler error on a project when switching from go1.1.2 to go1.2.
