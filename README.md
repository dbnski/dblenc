# dblenc

**dblenc** is a Go library for correcting double-encoded (and multiply-encoded) text values that can be found in some MySQL databases. The code was inspired by [`golang.org/x/text/encoding/charmap`](https://pkg.go.dev/golang.org/x/text/encoding/charmap), but which cannot handle these conversions correctly due to differences between true Latin1 and its implementation in MySQL.

## Background: The Problem

In MySQL databases, especially those with legacy migration histories, it's not uncommon to encounter double-encoded text. It occurs when UTF-8 data is mistakenly interpreted and re-encoded as Latin1 (or similar) due to character set misconfiguration between the server and applications, resulting in unreadable or mangled strings like: Ã© instead of é.

## Features

- Detects if given byte string has been double- or multiply-encoded with reasonable accuracy.
- Corrects corrupted strings by reversing the encoding layers.
- Discards invalid trailing Unicode sequences.

## Caveat

Double-encoded strings use a small subset of UTF-8 characters. As a result, it can be difficult to distinguish between a valid UTF-8 string and a double-encoded one at a glance. Therefore, `Detect()` may occasionally return "double-encoded" for valid UTF-8 strings when such strings don't include any characters that are out of scope for double-encoded sequences. However, `Transform()` called on a false-positive will not alter the string. For a practical illustration, see the example in the [`example/`](example/) directory.
