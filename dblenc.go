package dblenc

import (
    "errors"
    "unicode/utf8"
)

var DEBUG = false

var (
    ErrInvalid = errors.New("invalid byte sequence")
    ErrNoop    = errors.New("nothing changed")
)

var charMap = [256]uint32{
    0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007,
    0x0008, 0x0009, 0x000A, 0x000B, 0x000C, 0x000D, 0x000E, 0x000F,
    0x0010, 0x0011, 0x0012, 0x0013, 0x0014, 0x0015, 0x0016, 0x0017,
    0x0018, 0x0019, 0x001A, 0x001B, 0x001C, 0x001D, 0x001E, 0x001F,
    0x0020, 0x0021, 0x0022, 0x0023, 0x0024, 0x0025, 0x0026, 0x0027,
    0x0028, 0x0029, 0x002A, 0x002B, 0x002C, 0x002D, 0x002E, 0x002F,
    0x0030, 0x0031, 0x0032, 0x0033, 0x0034, 0x0035, 0x0036, 0x0037,
    0x0038, 0x0039, 0x003A, 0x003B, 0x003C, 0x003D, 0x003E, 0x003F,
    0x0040, 0x0041, 0x0042, 0x0043, 0x0044, 0x0045, 0x0046, 0x0047,
    0x0048, 0x0049, 0x004A, 0x004B, 0x004C, 0x004D, 0x004E, 0x004F,
    0x0050, 0x0051, 0x0052, 0x0053, 0x0054, 0x0055, 0x0056, 0x0057,
    0x0058, 0x0059, 0x005A, 0x005B, 0x005C, 0x005D, 0x005E, 0x005F,
    0x0060, 0x0061, 0x0062, 0x0063, 0x0064, 0x0065, 0x0066, 0x0067,
    0x0068, 0x0069, 0x006A, 0x006B, 0x006C, 0x006D, 0x006E, 0x006F,
    0x0070, 0x0071, 0x0072, 0x0073, 0x0074, 0x0075, 0x0076, 0x0077,
    0x0078, 0x0079, 0x007A, 0x007B, 0x007C, 0x007D, 0x007E, 0x007F,
    0x20AC, 0x0081, 0x201A, 0x0192, 0x201E, 0x2026, 0x2020, 0x2021,
    0x02C6, 0x2030, 0x0160, 0x2039, 0x0152, 0x008D, 0x017D, 0x008F,
    0x0090, 0x2018, 0x2019, 0x201C, 0x201D, 0x2022, 0x2013, 0x2014,
    0x02DC, 0x2122, 0x0161, 0x203A, 0x0153, 0x009D, 0x017E, 0x0178,
    0x00A0, 0x00A1, 0x00A2, 0x00A3, 0x00A4, 0x00A5, 0x00A6, 0x00A7,
    0x00A8, 0x00A9, 0x00AA, 0x00AB, 0x00AC, 0x00AD, 0x00AE, 0x00AF,
    0x00B0, 0x00B1, 0x00B2, 0x00B3, 0x00B4, 0x00B5, 0x00B6, 0x00B7,
    0x00B8, 0x00B9, 0x00BA, 0x00BB, 0x00BC, 0x00BD, 0x00BE, 0x00BF,
    0x00C0, 0x00C1, 0x00C2, 0x00C3, 0x00C4, 0x00C5, 0x00C6, 0x00C7,
    0x00C8, 0x00C9, 0x00CA, 0x00CB, 0x00CC, 0x00CD, 0x00CE, 0x00CF,
    0x00D0, 0x00D1, 0x00D2, 0x00D3, 0x00D4, 0x00D5, 0x00D6, 0x00D7,
    0x00D8, 0x00D9, 0x00DA, 0x00DB, 0x00DC, 0x00DD, 0x00DE, 0x00DF,
    0x00E0, 0x00E1, 0x00E2, 0x00E3, 0x00E4, 0x00E5, 0x00E6, 0x00E7,
    0x00E8, 0x00E9, 0x00EA, 0x00EB, 0x00EC, 0x00ED, 0x00EE, 0x00EF,
    0x00F0, 0x00F1, 0x00F2, 0x00F3, 0x00F4, 0x00F5, 0x00F6, 0x00F7,
    0x00F8, 0x00F9, 0x00FA, 0x00FB, 0x00FC, 0x00FD, 0x00FE, 0x00FF,
}

type Encoding byte
const (
    UNKNOWN                  Encoding = iota
    ASCII
    MAYBE_OTHER
    OTHER_CHARSET
    MAYBE_DOUBLE_ENCODED
    DOUBLE_ENCODED_TRUNCATED
    DOUBLE_ENCODED
    ERROR
)

func (r Encoding) String() string {
    switch r {
    case ASCII:
        return "ascii"
    case OTHER_CHARSET:
        return "other-charset"
    case DOUBLE_ENCODED:
        return "double-encoded"
    case MAYBE_DOUBLE_ENCODED:
        return "maybe-double-encoded"
    case DOUBLE_ENCODED_TRUNCATED:
        return "double-encoded-truncated"
    case MAYBE_OTHER:
        return "maybe-other"
    case ERROR:
        return "error"
    default:
        return "unknown"
    }
}

type byteMap struct {
    byteMap [256]byte
    next    [256]*byteMap
}

func newByteMap() *byteMap {
    m := &byteMap{}
    buf := make([]byte, utf8.UTFMax)

    for i, u := range charMap {
        ptr := m

        r := rune(u)
        n := utf8.EncodeRune(buf, r)

        for j := 0; j < utf8.UTFMax; j++ {
            code := buf[j]
            if n == j + 1 {
                ptr.byteMap[code] = byte(i)
                break
            }
            if ptr.next[code] == nil {
                ptr.next[code] = &byteMap{}
            }
            ptr = ptr.next[code]
        }
    }

    return m
}

type Callback func(...interface{})

type Decoder struct {
    byteMap *byteMap

    onRune      Callback
    onTransform Callback
}

func NewDecoder() *Decoder {
    d := &Decoder{
        byteMap: newByteMap(),
    }
    return d
}

func (d *Decoder) OnRune(callback Callback) *Decoder {
    d.onRune = callback
    return d
}

func (d *Decoder) OnTransform(callback Callback) *Decoder {
    d.onTransform = callback
    return d
}

// The function quickly tests a byte slice for presence of double-encoded
// characters. It will correctly classify strings that are not double-encoded,
// but may occasionally flag legitimate UTF-8 strings as double-encoded - such
// cases require separate verification which is intentionally omitted here
// in favour of the function's performance.
func (d *Decoder) Detect(data []byte) (Encoding, int, int, int) {
    // fast path for strings with long ascii prefixes
    f := 0 // fast-forward index
    data = data[:len(data):len(data)]
    for len(data) >= 8 {
        c1 := uint32(data[0]) | uint32(data[1]) << 8 |
              uint32(data[2]) << 16 | uint32(data[3]) << 24
        c2 := uint32(data[4]) | uint32(data[5]) << 8 |
              uint32(data[6]) << 16 | uint32(data[7]) << 24
        // test the highest bit in each byte code.
        if (c1 | c2) & 0x80808080 != 0 {
            // not ascii
            break
        }
        f += 8
        data = data[8:]
    }

    r := ASCII      // analysis result
    m := d.byteMap  // character map pointer

    a := 0          // ascii count
    c := 0          // code point counter
    e := 0          // double-encoded code point counter
    i := 0          // buffer position index
    n := 0          // decoded code units counter

    o := len(data)  // position of the first double-encoded sequence
    x := byte(0)    // decoded code unit
    u := uint32(0)  // decoded code unit sequence
    s := 1          // decoded code unit sequence size

    var currentRune rune
    var runeSequence [5]rune
    var sequenceLength int
    var isMultiple bool
    var isLatin bool = true
    var isLanguage Language = ^Language(0)
    var isDecodedLanguage Language = ^Language(0)

    for i < len(data) {
        // ASCII
        // FIRST BYTE
        currentByte := data[i]
        i++

        if currentByte < 0x80 {                 // ascii?
            if r == UNKNOWN {                   // incomplete sequence followed by an ascii
                return OTHER_CHARSET, c, e, f + i
            }
            a++
            c++

            continue
        }
        if currentByte < 0xC0 {                 // 0x80 - 0xBF cannot appear stand-alone
            return OTHER_CHARSET, c, e, f + i
        }

        m := m.next[currentByte]
        if m == nil {                           // byte sequence does not appear
            return OTHER_CHARSET, c, e, f + i   // in the map
        }
        if i == len(data) {                     // buffer ends mid-sequence
            return ERROR, c, e, f + i
        }
        firstByte := currentByte

        // SECOND BYTE
        currentByte = data[i]
        i++

        if m.byteMap[currentByte] != 0 {        // matches complete double-encoded character
            x = m.byteMap[currentByte]
            c++
            n++

            currentRune = rune(firstByte & 0x1F) << 6 | rune(currentByte & 0x3F)
            if isLatin {
                latin := false
                if int(currentRune) < len(Diacritics) {
                    latin = Diacritics[currentRune] > 0
                    isLanguage = isLanguage & Diacritics[currentRune]
                }
                isLatin = latin
            }

            if e > 0 {
                isMultiple = runeSequence[sequenceLength] != currentRune
            }
            runeSequence[sequenceLength] = currentRune
            sequenceLength++

            if n == 1 {                         // first byte of decoded code point
                switch {
                case x & 0xE0 == 0xC0:          // 2-byte code point
                    if x < 0xC2 {
                        return OTHER_CHARSET, c, e, f + i
                    }
                    s = 2
                case x & 0xF0 == 0xE0:          // 3-byte code point
                    s = 3
                case x & 0xF8 == 0xF0:          // 4-byte code point
                    if x >= 0xF5 {
                        return OTHER_CHARSET, c, e, f + i
                    }
                    s = 4
                default:                        // not utf8
                    return OTHER_CHARSET, c, e, f + i
                }
                u = uint32(x)
                r = UNKNOWN
            } else {                            // continuation bytes of decoded code point
                if (x & 0xC0) != 0x80 {         // check if valid continuation byte
                    return OTHER_CHARSET, c, e, f + i
                }
                u = (u << 8) | uint32(x)

                if n == s {                     // decoded complete code unit sequence
                    if s == 2 {
                        decodedRune := rune((((u >> 8) & 0x1F) << 6) | ((u & 0xFF) & 0x3F))
                        if int(decodedRune) < len(DecodedDiacritics) {
                            isDecodedLanguage = isDecodedLanguage & DecodedDiacritics[decodedRune]
                        }
                    }
                    if s == 3 {
                        // UTF16 code points
                        if u >= 0xEDA080 && u <= 0xEDBFBF {
                            return OTHER_CHARSET, c, e, f + i
                        }
                    }
                    if s == 4 {
                        // UTF16 code points
                        if u >= 0xF4908080 {
                            return OTHER_CHARSET, c, e, f + i
                        }
                    }
                    // out-of-scope code points
                    if u > 0xF3A087BF {
                        return OTHER_CHARSET, c, e, f + i
                    }

                    if d.onRune != nil {
                        d.onRune(sequenceLength, runeSequence)
                    }

                    r = DOUBLE_ENCODED
                    e++                         // found double-encoded code point
                    n = 0                       // reset decoded code units counter
                    u = 0                       // reset decoded char

                    sequenceLength = 0
                }
            }

            o = min(o, i - 1)

            continue
        }

        m = m.next[currentByte]
        if m == nil {
            return OTHER_CHARSET, c, e, f + (i - 1)
        }
        if i == len(data) {
            return ERROR, c, e, f + i
        }
        secondByte := currentByte

        // THIRD BYTE
        currentByte = data[i]
        i++

        if m.byteMap[currentByte] != 0 {        // matches complete double-encoded character
            x = m.byteMap[currentByte]
            c++
            n++

            currentRune = rune(firstByte & 0x0F) << 12 | rune(secondByte & 0x3F) << 6 | rune(currentByte & 0x3F)
            if isLatin {
                latin := false
                if int(currentRune) < len(Diacritics) {
                    latin = Diacritics[currentRune] > 0
                    isLanguage = isLanguage & Diacritics[currentRune]
                }
                isLatin = latin
            }

            if e > 0 {
                isMultiple = runeSequence[sequenceLength] != currentRune
            }
            runeSequence[sequenceLength] = currentRune
            sequenceLength++

            if n == 1 {                         // analyse the first byte
                switch {
                case x & 0xE0 == 0xC0:          // 2-byte code point
                    if x < 0xC2 {
                        return OTHER_CHARSET, c, e, f + i
                    }
                    s = 2
                case x & 0xF0 == 0xE0:          // 3-byte code point
                    s = 3
                case x & 0xF8 == 0xF0:          // 4-byte code point
                    if x >= 0xF5 {
                        return OTHER_CHARSET, c, e, f + i
                    }
                    s = 4
                default:                        // not utf8
                    return OTHER_CHARSET, c, e, f + i
                }
                u = uint32(x)
                r = UNKNOWN
            } else {                            // analyse continuation bytes
                if (x & 0xC0) != 0x80 {         // check if valid continuation byte
                    return OTHER_CHARSET, c, e, f + i
                }
                u = (u << 8) | uint32(x)

                if n == s {                     // decoded complete code unit sequence
                    if s == 2 {
                        decodedRune := rune((((u >> 8) & 0x1F) << 6) | ((u & 0xFF) & 0x3F))
                        if int(decodedRune) < len(DecodedDiacritics) {
                            isDecodedLanguage = isDecodedLanguage & DecodedDiacritics[decodedRune]
                        }
                    }
                    if s == 3 {
                        // UTF16 code points
                        if u >= 0xEDA080 && u <= 0xEDBFBF {
                            return OTHER_CHARSET, c, e, f + i
                        }
                    }
                    if s == 4 {
                        // UTF16 code points
                        if u >= 0xF4908080 {
                            return OTHER_CHARSET, c, e, f + i
                        }
                    }
                    if s > 0xF3A087BF {
                        return ERROR, c, e, f + i
                    }

                    if d.onRune != nil {
                        d.onRune(sequenceLength, runeSequence)
                    }

                    r = DOUBLE_ENCODED
                    e++                         // found double-encoded code point
                    n = 0                       // reset decoded code units counter
                    u = 0                       // reset decoded char

                    sequenceLength = 0
                }
            }

            o = min(o, i - 2)

            continue
        }

        // FOURTH BYTE
        return OTHER_CHARSET, c, e, f + (i - 2) // no 4-byte code points exist
    }

    if r != ASCII && d.onRune != nil && sequenceLength > 0 {
        d.onRune(sequenceLength, runeSequence)
    }

    switch {
    case r == ASCII:
        // do not touch me

    case r == OTHER_CHARSET:
        panic("we should not be here")

    case r == ERROR:
        panic("we should not be here")

    case r == UNKNOWN:                          // if string ends halfway in what could be a double encoded character
        switch {
        case isLatin:                           // if all suspects are made exclusively of cp1252 letters,
            r = MAYBE_OTHER                     // assume the string is not double encoded (e.g. "Úžasná")

        case e > 0:                             // if there's at least one other suspect,
            r = DOUBLE_ENCODED_TRUNCATED        // assume it's a truncated double encoded string (e.g. "MATÄšJ [..] Tomáš")

        case isClosingPunctuation(currentRune): // if it's the only suspect and the final char is "closing" punctuation (e.g. "qué¡"),
            r = MAYBE_OTHER                     // assume it's not double encoded
        }

    case r == DOUBLE_ENCODED:                   // if the string was classified as double encoded
        if !isMultiple {                        // but, it has just one unique double encoded sequence,
            r = MAYBE_DOUBLE_ENCODED            // assume it's double-encoded (e.g. "ÄŽakujem [..] ÄŽakujem")

            if isLatin {                        // if the suspect is made exclusively of cp1252 letters
                if isLanguage > 0 {             // and all those letters are used by the same language,
                    r = MAYBE_OTHER             // assume it's not double encoded (e.g. "Úžasna")
                }
                if isDecodedLanguage > 0 &&            // except if the decoded letter(s) is a known exception
                   isDecodedLanguage < ^Language(0) {  // like ĊČĎĚĞğġŌŞşŚƟΟ (e.g. "DoÄŸan" -> "Doğan")
                    r = MAYBE_DOUBLE_ENCODED
                }
            }
        }
    }

    return r, c, e, f + min(o, i)
}

func (d *Decoder) Transform(b []byte) ([]byte, error) {
    o := b

    enc, _, _, _ := d.Detect(o)
    // try to recover from an incomplete trailing sequence
    if enc == ERROR {
        for p := len(b) - 1; p >= 0 && p >= len(b) - utf8.UTFMax; p-- {
            if o[p] == 0xC2 || o[p] == 0xC3 || o[p] == 0xC5 ||
                o[p] == 0xC6 || o[p] == 0xCB || o[p] == 0xE2 {
                // re-check the shorter string
                enc, _, _, _ = d.Detect(o[:p])
                if enc != ERROR {
                    // discard the broken sequence
                    o = o[:p]
                    if d.onTransform != nil {
                        d.onTransform(enc, o)
                    }
                }
                break
            }
        }
    }

    transformErr := ErrNoop

    for enc == DOUBLE_ENCODED || enc == DOUBLE_ENCODED_TRUNCATED || enc == MAYBE_DOUBLE_ENCODED {
        x, err := d.transform(o)
        if err != nil {
            break
        }
        if d.onTransform != nil {
            d.onTransform(enc, x)
        }

        // validate the suffix if it's not an ascii char
        if x[len(x) - 1] >= 0x80 {
            p := len(x) - 1
            // search for the start of a utf8 sequence
            for p >= 0 && p >= len(x) - utf8.UTFMax {
                if x[p] >= 0xC2 && x[p] <= 0xF4 {
                    valid := utf8.Valid(x[p:])
                    if !valid {
                        // trailing sequence is invalid, discard it
                        x = x[:p]
                        if d.onTransform != nil {
                            d.onTransform(enc, x)
                        }
                    }
                    break
                }
                p--
            }
            if p < len(x) - utf8.UTFMax {
                // no ascii or utf8 sequence found
                break
            }
        }

        valid := utf8.Valid(x)
        if !valid {
            // this iteration got us nowhere good
            // break
        }

        transformErr = nil
        enc, _, _, _ = d.Detect(x)

        o = x  // new candidate
    }

    if transformErr != nil {
        o = b
    }

    return o, transformErr
}

func (this *Decoder) TransformOnce(src []byte) (dst []byte, err error) {
    return this.transform(src)
}

func (d *Decoder) transform(src []byte) (dst []byte, err error) {
    // fast path for strings with long ascii prefixes
    pDst := 0
    src = src[:len(src):len(src)]
    for len(src) >= 8 {
        c1 := uint32(src[0]) | uint32(src[1]) << 8 |
              uint32(src[2]) << 16 | uint32(src[3]) << 24
        c2 := uint32(src[4]) | uint32(src[5]) << 8 |
              uint32(src[6]) << 16 | uint32(src[7]) << 24
        // test the highest bit in each byte code.
        if (c1 | c2) & 0x80808080 != 0 {
            // not ascii
            break
        }

        if pDst > len(dst) - 8 {
            dst = grow(dst, pDst)
        }
        copy(dst[pDst:], src[:8])
        pDst += 8

        src = src[8:]
    }

    m := d.byteMap  // character map pointer
    n := 0          // decoded code units counter
    x := byte(0)    // decoded code unit
    s := 1          // decoded code unit sequence size
    u := uint32(0)

    pSrc := 0
    for pSrc < len(src) {
        // FIRST BYTE
        currentByte := src[pSrc]
        pSrc++

        if currentByte < 0x80 {                 // ascii?
            if pDst == len(dst) {
                dst = grow(dst, pDst)
            }
            dst[pDst] = currentByte
            pDst++

            continue
        }

        if currentByte < 0xC0 {                 // 0x80 - 0xBF cannot appear stand-alone
            return nil, ErrInvalid
        }

        m := m.next[currentByte]
        if m == nil {                           // byte sequence does not appear
            return nil, ErrInvalid              // in the map
        }
        if pSrc == len(src) {                   // buffer ends mid-sequence
            return nil, ErrInvalid
        }

        // SECOND BYTE
        currentByte = src[pSrc]
        pSrc++

        if m.byteMap[currentByte] != 0 {        // matches complete double-encoded character
            x = m.byteMap[currentByte]
            n++

            if n == 1 {                         // first byte of decoded code point
                switch {
                case x & 0xE0 == 0xC0:          // 2-byte code point
                    if x < 0xC2 {
                        return nil, ErrInvalid
                    }
                    s = 2
                case x & 0xF0 == 0xE0:          // 3-byte code point
                    s = 3
                case x & 0xF8 == 0xF0:          // 4-byte code point
                    if x >= 0xF5 {
                        return nil, ErrInvalid
                    }
                    s = 4
                default:                        // not utf8
                    return nil, ErrInvalid
                }
            } else {                            // continuation bytes of decoded code point
                if (x & 0xC0) != 0x80 {         // check if valid continuation byte
                    return nil, ErrInvalid
                }
                u = (u << 8) | uint32(x)

                if n == s {                     // decoded complete code unit sequence
                    if s == 3 {
                        // UTF16 code points
                        if u >= 0xEDA080 && u <= 0xEDBFBF {
                            return nil, ErrInvalid
                        }
                    }
                    if s == 4 {
                        // UTF16 code points
                        if u >= 0xF4908080 {
                            return nil, ErrInvalid
                        }
                    }
                    // out-of-scope code points
                    if u > 0xF3A087BF {
                        return nil, ErrInvalid
                    }

                    n = 0                       // reset decoded code units counter
                    u = 0                       // reset decoded char
                }
            }

            if pDst == len(dst) {
                dst = grow(dst, pDst)
            }
            dst[pDst] = x
            pDst++

            continue
        }

        m = m.next[currentByte]
        if m == nil {
            return nil, ErrInvalid
        }
        if pSrc == len(src) {
            return nil, ErrInvalid
        }

        // THIRD BYTE
        currentByte = src[pSrc]
        pSrc++

        if m.byteMap[currentByte] != 0 {        // matches complete double-encoded character
            x = m.byteMap[currentByte]
            n++

            if n == 1 {                         // first byte of decoded code point
                switch {
                case x & 0xE0 == 0xC0:          // 2-byte code point
                    if x < 0xC2 {
                        return nil, ErrInvalid
                    }
                    s = 2
                case x & 0xF0 == 0xE0:          // 3-byte code point
                    s = 3
                case x & 0xF8 == 0xF0:          // 4-byte code point
                    if x >= 0xF5 {
                        return nil, ErrInvalid
                    }
                    s = 4
                default:                        // not utf8
                    return nil, ErrInvalid
                }
            } else {                            // continuation bytes of decoded code point
                if (x & 0xC0) != 0x80 {         // check if valid continuation byte
                    return nil, ErrInvalid
                }
                u = (u << 8) | uint32(x)

                if n == s {                     // decoded complete code unit sequence
                    if s == 3 {
                        // UTF16 code points
                        if u >= 0xEDA080 && u <= 0xEDBFBF {
                            return nil, ErrInvalid
                        }
                    }
                    if s == 4 {
                        // UTF16 code points
                        if u >= 0xF4908080 {
                            return nil, ErrInvalid
                        }
                    }
                    // out-of-scope code points
                    if u > 0xF3A087BF {
                        return nil, ErrInvalid
                    }

                    n = 0                       // reset decoded code units counter
                    u = 0                       // reset decoded char
                }
            }

            if pDst == len(dst) {
                dst = grow(dst, pDst)
            }
            dst[pDst] = x
            pDst++

            continue
        }

        // FOURTH BYTE
        return nil, ErrInvalid
    }

    dst = dst[:pDst:pDst]

    return dst, nil
}


func grow(b []byte, n int) []byte {
    m := len(b)
    if m <= 32 {
        m = 64
    } else if m <= 256 {
        m *= 2
    } else {
        m += m >> 1
    }
    buf := make([]byte, m)
    copy(buf, b[:n])
    return buf
}