# Hit

### Description

Hit is a CLI tool that provides HTTP server performance metrics

Developed by following [Effective Go by Inanc Gumus](https://www.manning.com/books/effective-go)

### Usage

To build binaries, use:
```bash
make
```

To run the tool, use:
```bash
bin/hit_darwin_arm64 -n 1000 -c 10 http://localhost:3000
```
To list all available flags, use:
```bash
bin/hit_darwin_arm64 -h
```
