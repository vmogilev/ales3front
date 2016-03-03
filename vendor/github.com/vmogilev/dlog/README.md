# dlog

It's the D Logger - simple logger pkg to setup four handlers for your basic logging needs:

* Trace
* Info
* Warning
* Error

## Example

    package main

    import (
        "flag"
        "io/ioutil"
        "os"

        "github.com/vmogilev/dlog"
    )

    func main() {
        var myVar string
        var debug bool

        flag.StringVar(&myVar, "myVar", "", "Mandatory Var")
        flag.BoolVar(&debug, "debug", false, "Debug")
        flag.Parse()

        if debug {
            dlog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
        } else {
            dlog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
        }

        if myVar == "" {
            dlog.Error.Fatalf("myVar is needed! Exiting ...")
        }

        dlog.Info.Printf("thanks - got myVar=%s\n", myVar)
    }

## Example Output of Above

    # dlog-example
    ERROR: 2015/06/11 16:33:06 main.go:26: myVar is needed! Exiting ...

    # dlog-example --help
    Usage of dlog-example:
      -debug=false: Debug
      -myVar="": Mandatory Var

    # dlog-example --myVar="Here You Go"
    INFO: 2015/06/11 16:33:57 main.go:29: thanks - got myVar=Here You Go

## License

[The MIT License](http://opensource.org/licenses/MIT)


