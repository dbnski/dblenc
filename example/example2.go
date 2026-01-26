package main

import (
    "fmt"
    "os"
    "unicode/utf8"

    "github.com/dbnski/dblenc"
    "github.com/dbnski/go-helpers/binary"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("usage: %s <string> ...\n", os.Args[0])
        os.Exit(1)
    }

    xformer := dblenc.NewDecoder().
        OnTransform(func (encoding dblenc.Encoding, data []byte) {
            fmt.Printf("# transform encoding=%s value=\"%s\" bytes=[%s]\n",
                encoding, data, binary.HexifyBytesToString(data))
        })

    detector := dblenc.NewDecoder().
        OnRune(func (encoded []byte) {
            decoded, _ := xformer.JustTransform(encoded)
            decodedRune, decodedLength := utf8.DecodeRune(decoded)
            fmt.Printf("# suspect encoded=\"%s\" length=%d bytes=[%s] rune=%x decoded=\"%s\"\n",
                encoded, decodedLength, binary.HexifyBytesToString(encoded), decodedRune, string(decodedRune))
        })


    for i := 1; i < len(os.Args); i++ {
        value := []byte(os.Args[i])
        encoding, chars, suspects, _ := detector.Detect(value)
        decoded, _ := xformer.Transform(value)
        fmt.Printf("detected=%s length=%d chars=%d suspects=%d decoded=\"%s\"\n",
            encoding, len(value), chars, suspects, decoded)
    }
}