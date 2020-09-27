# vmt
![](https://github.com/HagaSpa/vmt/workflows/go%20test%20&%20go%20build/badge.svg)

vmt is Virtual Machine Translater for nand2tetris.
Virtual Machine Translater translates VM programs into assembly code.


## Requirements
* Docker 19.03 or later

You need a docker environment to build with Docker.

Also, since it uses Docker BuildKit, please using enabled the [Docker BuildKit flag](https://docs.docker.com/develop/develop-images/build_enhancements/) by using Docker 19.03 or later.


## Build
This Virtual Machine Translater can cross compile the binary for the host operating system by adding arguments to make command. 
(using cross compile in golang)

```
// compile the binary for you host operating system.
$ make

// macOS
$ make PLATFORM=darwin/amd64 

// linux
$ make PLATFORM=linux/amd64

// Windows
& make PLATFORM=windows/amd64
```

## Args
```
$./bin/main {arg1}
```
`{arg1}` is vm file name or the directory name with multiple vm files.


## Run
```
$ make

$ ./bin/main StaticsTest
2020/09/27 10:33:22 translated multiple vm: StaticsTest.asm

$ cat StaticsTest.asm 

// initialize asm
@256
D=A
@SP
M=D

// call Sys.init args nums 0
@Sys.init$0
D=A
@SP
A=M
M=D
@SP
M=M+1

// push register LCL
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1

// push register ARG
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
.
.
.
```

`Note: Confirmed to work only on macOS now`
