# go-source
`go-source` is a tool for "chunk"ing source code. Check it out in action! 


###  Installation
```bash 
git clone https://github.com/cheese-head/go-source
go build -o go-source
```


### Usage 
```bash
./go-source --help
NAME:
   go-source - A new cli application

USAGE:
   go-source [global options] command [command options] 

DESCRIPTION:
   A tool for parsing and chunking source code

COMMANDS:
   chunk    chunk files
   read     read content from chunks
   server   start a http server
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

```bash
./go-source chunk --files "sample.py" --include ".py"
./go-source chunk --files "pkg/*" --include ".go"
./go-source chunk --files "pkg/*" --include ".go" | ./go-source read
./go-source server
```