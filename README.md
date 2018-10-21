# app - the Library to Create the Most Convinient Command Line Applications in Go

It started when I tried to set up bash completion at [urfave/cli.v2](https://godoc.org/gopkg.in/urfave/cli.v2) and failed.

In the end I decided to recreate all the library from my point of view. Mostly it has similar workflow.

## Key Points of the Project

### Fully customizable

You can set your handlers to these actions:
* Command action
* Command Before/After actions
* Command Completion action
* Command Help action
* Command Error handler
* Flag Parse method (implementing interface)
* Flag Before/After actions
* Flag Completion action
* Default command help action
* Default command completion action
* Default flag completion action

### Convinient default behaviour / fast setup for most applications

It's needed about 10 lines of code per command to get started.

Flags could be placed after command they defined at and up to the end of arguments if not shadowed.

Everything could be completed by default functions.

### Simple and efficient code

Performance is not the key feature of the project because all the operations are simple and does only O(1) times per application run. Altough it does only what you ask and what is needed, so no problems here.

```
BenchmarkUrfaveAppCreateRun-8   	   30000	     42632 ns/op	   12633 B/op	     413 allocs/op
BenchmarkNikAppCreateRun-8      	 1000000	      1618 ns/op	    1432 B/op	      15 allocs/op
```

## Bash Completion

It supports reach bash completion out of the box. Subcommands and flags are completed from app configuration. Leaf command arguments are completed by default bash `ls`like completion. Even flag arguments could be completed. You can define custom completion behavoiur: just show help message, show some alternatives to choose from, alternatives with spaces works fine, execute default bash `ls`like completion, combine these methods, code your own completion having raw information about command and cursor position and haveing all the bash power to postprocess results.
