package main

import (
    "io"
    "os"
    "strings"
)

func main() {
    fs := os.Args[1:]
    out, _ := os.Create("textfile_constants.go")
    out.Write([]byte("package main \n\nconst (\n"))
    for _, f := range fs {
        out.Write([]byte(strings.TrimSuffix(f, ".txt") + " = `"))
        f, _ := os.Open(f)
        io.Copy(out, f)
        out.Write([]byte("`\n"))
    }
    out.Write([]byte(")\n"))
}
