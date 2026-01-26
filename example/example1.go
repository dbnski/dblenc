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
        // Double-encoded version of "生命冲力".
        []byte("\xC3\xA7\xE2\x80\x9D\xC5\xB8\xC3\xA5\xE2\x80\x98\xC2\xBD\xC3\xA5\xE2\x80\xA0\xC2\xB2\xC3\xA5\xC5\xA0\xE2\x80\xBA"),
        // Valid UTF-8.
        []byte("élan vital"),
        // Double-encoded version of "élan vital".
        []byte("Ã©lan vital"),
    }

    decoder := dblenc.NewDecoder()

    for i, sample := range(samples) {
        result, _, _, _ := decoder.Detect(sample)
        fixed, err := decoder.Transform(sample)
        if err != nil && err != dblenc.ErrNoop {
            panic(err)
        }

        fmt.Printf(
            "[%d] before: %s, suspected type: %s, after: %s\n",
            i + 1, string(sample), result, string(fixed),
        )
    }
}
