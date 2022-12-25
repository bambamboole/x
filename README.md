# X
Early alpha stage. Until now it is playground of ideas and everything can change.

`x` is a Bash task runner written in Go. It has the goal to make your
utility scripts as easy as possible to execute. It is heavily inspired by
`Makefile` and [`b5`](https://github.com/team23/b5)

## Installation
X is build in Go and has therefor a single dependency free binary. You can find 
the right one for your platform attached to the [releases](https://github.com/bambamboole/x/releases).
### Homebrew
x is available via Homebrew by adding a custom tap:
```shell
brew tap bambamboole/x
brew install x
```

## Features
* Taskfile
* config
* It recursively iterates up, so you can use it in a sub folder
* It is fast

### Taskfile
A `Taskfile` is a file which contains task definition in bash syntax. 
Here a simple example:

```shell
#!/bin/bash

task:hello (){
  echo "Hello $@!"
}
```
This is a very simple example, but it demonstrates a big feature. 
It os possible to pass through arguments. A feature which is 
missing in a `Makefile`.

To execute this task now, you have to execute `x hello`, or better with an
additional argument `x hello world`.
