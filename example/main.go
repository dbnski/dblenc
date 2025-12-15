package main

import (
    "fmt"
    "github.com/dbnski/dblenc"
)

func main() {
    samples := [][]byte {
        // Simple ASCII.
        []byte("elan vital"),
        // Valid UTF-8.
        []byte("生命冲力"),
        // Valid UTF-8, but lacks any character sequences that would
        // clearly differentiate it from a double-encoded string.
        []byte("élan vital"),
        // Double-encoded version of "élan vital".
        []byte("Ã©lan vital"),
    }

    decoder := dblenc.NewDecoder()

    for i, sample := range(samples) {
        result, _ := decoder.Detect(sample)
        fixed, err := decoder.Transform(sample)
        if err != nil {
            panic(err)
        }

        fmt.Printf(
            "[%d] original: %s, suspected type: %s, transformed: %s\n",
            i + 1, string(sample), result, string(fixed),
        )
    }
}
