# Creating a container from scratch

This repository contains the demo files for how to create a Linux container from scratch. For detailed information about the files in this repo and how to use them go to https://hectorcorrea.com/blog/tiny-container/83

## Quick rundown 

If you want to compile the code follow the following steps on a Linux machine with Go installed:

```
$ git clone https://github.com/hectorcorrea/tiny-container.git
$ cd tiny-container
$ GOOS=linux go build -o tc tinyContainer.go
$ GOOS=linux go build -o ts tinyShell.go

$ ./tc -root=/root/tiny-container -shell=./ts
Tiny shell started
ts: _
```

## Quick rundown (without the source code)

Download TinyContainer (`tc`) and TinyShell (`ts`) from https://github.com/hectorcorrea/tiny-container/releases and run

```
$ pwd
/root/tiny-container

$ ./tc -root=/root/tiny-container -shell=./ts
Tiny shell started
ts: _
```

