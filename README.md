# dblenc

**dblenc** is a Go library for correcting double-encoded (and multiply-encoded) text values that can be found in some MySQL databases. The code was inspired by [`golang.org/x/text/encoding/charmap`](https://pkg.go.dev/golang.org/x/text/encoding/charmap), but which cannot handle these conversions correctly due to slight differences with the character map MySQL uses.

## Background: The Problem

In MySQL databases, especially those with legacy migration histories, it's not uncommon to encounter double-encoded text. It occurs when UTF-8 data is mistakenly interpreted and re-encoded as Latin1 (or similar) due to character set misconfiguration between the server and applications, resulting in unreadable or mangled strings like: Ã© instead of é.

## Features

- Detects if given byte string has been double- or multiply-encoded with reasonable accuracy.
- Corrects corrupted strings by reversing the encoding layers.
- Discards invalid trailing Unicode sequences.

## Caveat

Double‑encoded strings consist of regular UTF‑8 characters and are themselves well‑formed UTF‑8. There are many edge cases that make it difficult to distinguish a valid UTF‑8 string from a double‑encoded one. As a result, `Detect()` may occasionally produce false positives when a string contains only characters that fall within the range of possible double‑encoded sequences. Similarly, it may produce false negatives when a string includes only a single double‑encoded character. Both cases should be very rare and limited to text in a language whose character set uses - or partially overlaps with - the Windows‑1252 code page.

## Usage

For a practical illustration, see the example in the [`example/`](example/) directory.
