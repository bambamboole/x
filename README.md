# X
`x` is a Bash task runner written in Go. It has the goal to make your
utility scripts as easy as possible to execute. It is heavily inspired by
`Makefile` and [`b5`](https://github.com/team23/b5)

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
