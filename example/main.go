package main

import (
    "encoding/hex"
    "fmt"

    "github.com/dbnski/dblenc/dblenc"
)

func decode(s string) []byte {
    if b, err := hex.DecodeString(s); err != nil {
        panic(err)
    } else {
        return b
    }
}

func main() {
    items := []struct {
        String  string
        Encoded []byte
    }{
        {
            String:  "Plain ascii string",
            Encoded: []byte("Plain ascii string"),
        },
        {
            String:  "全然分からない",
            Encoded: decode("c3a5e280a6c2a8c3a7e2809ec2b6c3a5cb86e280a0c3a3c281e280b9c3a3e2809ae280b0c3a3c281c2aac3a3c281e2809e"),
        },
    }

    m := dblenc.NewUnDoubleEncoder()
    for i, item := range items {
        r, err := m.Transform(item.Encoded)
        if err != nil {
            panic(err)
        }
        if i > 0 {
            fmt.Println()
        }
        fmt.Printf("Original:       %s\n", item.String)
        fmt.Printf("Double encoded: %s\n", item.Encoded)
        fmt.Printf("Recovered:      %s\n", r)
    }
}
