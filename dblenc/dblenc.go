package dblenc

import (
    "errors"
    // "fmt"
    "sort"
    "unicode/utf8"
)

var (
    ErrInvalid = errors.New("invalid characters")
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
    ASCII          Encoding = iota
    DOUBLE_ENCODED
    ERROR
    UNKNOWN
)

func (r Encoding) String() string {
    switch r {
    case ASCII:
        return "ascii"
    case DOUBLE_ENCODED:
        return "double-encoded"
    case ERROR:
        return "error"
    default:
        return "unknown"
    }
}

type byteMap struct {
    Map  [256]bool
    Next [256]*byteMap
}

func newByteMap() *byteMap {
    m := &byteMap{}
    buf := make([]byte, utf8.UTFMax)
    ptr := m

    for _, u := range charMap {
        // convert runes to bytes
        r := rune(u)
        n := utf8.EncodeRune(buf, r)

        if n == 1 {
            code := buf[0]
            ptr.Map[code] = true
        } else {
            code := buf[0]
            if ptr.Next[code] == nil {
                ptr.Next[code] = &byteMap{}
            }
            ptr := ptr.Next[code]
            if n == 2 {
                code = buf[1]
                ptr.Map[code] = true
            } else {
                code = buf[1]
                if ptr.Next[code] == nil {
                    ptr.Next[code] = &byteMap{}
                }
                ptr := ptr.Next[code]
                code = buf[2]
                ptr.Map[code] = true
            }
        }
    }

    return m
}

func (this *byteMap) detect(data []byte) (Encoding, int) {
    // Fast path for strings with long ascii prefixes.
    // Based on utf8.Valid().
    f := 0
    data = data[:len(data):len(data)]
    for len(data) >= 8 {
        c1 := uint32(data[0]) | uint32(data[1]) << 8 |
              uint32(data[2]) << 16 | uint32(data[3]) << 24
        c2 := uint32(data[4]) | uint32(data[5]) << 8 |
              uint32(data[6]) << 16 | uint32(data[7]) << 24
        // Test the highest bit in each byte code.
        if (c1 | c2) & 0x80808080 != 0 {
            // not ascii
            break
        }
        f += 8
        data = data[8:]
    }

    p := 0
    m := this
    s := len(data)
    r := ASCII
    for p < s {
        c := data[p]
        p++
        if m.Map[c] {
            // ascii
            continue
        }
        if p == s {
            return ERROR, f + p
        }
        m := m.Next[c]
        if m == nil {
            return UNKNOWN, f + p
        }

        c = data[p]
        p++
        if m.Map[c] {
            // complete double-encoded sequence (2 bytes)
            r = DOUBLE_ENCODED
            continue
        }
        if p == s {
            return ERROR, f + p
        }
        m = m.Next[c]
        if m == nil {
            return UNKNOWN, f + p
        }

        c = data[p]
        p++
        if m.Map[c] {
            // complete double-encoded sequence (3 bytes)
            r = DOUBLE_ENCODED
            continue
        }
        return UNKNOWN, f + p
    }

    return r, f + p
}

type UnDoubleEncoder struct {
    byteMap *byteMap
    runeMap []uint32
}

func NewUnDoubleEncoder() *UnDoubleEncoder {
    und := &UnDoubleEncoder{
        byteMap: newByteMap(),
        runeMap: make([]uint32, len(charMap)),
    }
    for i := 0; i < len(charMap); i++ {
        // Store the byte code in the upper byte and the rune value
        // in the lower bytes.
        und.runeMap[i] = (uint32(i) << 24) | (charMap[i] & 0x00ffffff)
    }
    // Sort rune map on rune value.
    sort.Slice(und.runeMap, func(i, j int) bool {
        return (und.runeMap[i] & 0x00ffffff) < (und.runeMap[j] & 0x00ffffff)
    })
    return und
}

func (this *UnDoubleEncoder) Transform(b []byte) ([]byte, error) {
    // {
    //     enc, at := this.Detect(b)
    //     fmt.Printf(
    //         "INPUT:\nstr   = %s\nhex   = %x\ntest  = %v [at %d/%d]\n\n",
    //         b, b[:], enc, at, len(b),
    //     )
    // }

    enc, _ := this.byteMap.detect(b)
    o := b
    for enc == DOUBLE_ENCODED {
        x, err := this.transform(o)
        if err != nil {
            return nil, err
        }

        // validate the suffix if it's not an ascii char
        if x[len(x) - 1] >= 0x80 {
            p := len(x) - 1
            // search for the start of a utf8 sequence in the last four bytes
            for p >= 0 && p >= len(x) - 4 {
                if x[p] >= 0xC2 && x[p] <= 0xF4 {
                    // {
                    //     valid := utf8.Valid(x[p:])
                    //     fmt.Printf(
                    //         "\tSUFFIX CHECK:\n\tstr   = %s\n\thex   = " +
                    //         "%x\n\tat    = %d/%d\n\tvalid = %v\n\n",
                    //         string(x[p:]), x[p:], p + 1, len(x), valid,
                    //     )
                    // }

                    valid := utf8.Valid(x[p:])
                    if !valid {
                        // discard the incomplete sequence
                        x = x[:p]
                    }
                    break
                }
                p--
            }
            if p < len(x) - 4 {
                // no ascii or utf8 sequence found
                return nil, ErrInvalid
            }
        }

        // {
        //     enc, at = this.Detect(x)
        //     valid := utf8.Valid(x)
        //     fmt.Printf(
        //         "\tPASS:\n\tstr   = %s\n\thex   = %x\n\ttest  = %v " +
        //         "[at %d/%d]\n\tvalid = %v\n\n",
        //         x, x[:], enc, at, len(x), utf8.Valid(x),
        //     )
        // }
        enc, _ = this.byteMap.detect(x)
        valid := utf8.Valid(x)
        if !valid {
            break
        }

        o = x   // new candidate
    }

    // fmt.Printf("OUTPUT:\nstr   = %s\nhex   = %x\n\n\n", o, o[:])

    return o, nil
}

func (this *UnDoubleEncoder) transform(src []byte) (dst []byte, err error) {
    pDst := 0
    src = src[:len(src):len(src)]
    // Fast path for strings with long ascii prefixes.
    // Based on utf8.Valid().
    for len(src) >= 8 {
        c1 := uint32(src[0]) | uint32(src[1]) << 8 |
              uint32(src[2]) << 16 | uint32(src[3]) << 24
        c2 := uint32(src[4]) | uint32(src[5]) << 8 |
              uint32(src[6]) << 16 | uint32(src[7]) << 24
        // Test the highest bit in each byte code.
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

    // Based on golang.org/x/text/encoding/charmap
    pSrc := 0
    r, size := rune(0), 0
    for pSrc < len(src) {
        r = rune(src[pSrc])

        if r < 0x80 {
            pSrc++
            if pDst == len(dst) {
                dst = grow(dst, pDst)
            }
            dst[pDst] = byte(r)
            pDst++
            continue
        } else {
            r, size = utf8.DecodeRune(src[pSrc:])
            if size == 1 {
                if !utf8.FullRune(src[pSrc:]) {
                    return nil, ErrInvalid
                }
            }
            pSrc += size
        }

        for low, high := 0x80, 0x100; ; {
            if low >= high {
                return nil, ErrInvalid
            }
            mid := (low + high) / 2
            got := this.runeMap[mid]
            gotRune := rune(got & (1<<24 - 1))
            if gotRune < r {
                low = mid + 1
            } else if gotRune > r {
                high = mid
            } else {
                if pDst == len(dst) {
                    dst = grow(dst, pDst)
                }
                dst[pDst] = byte(got >> 24)
                pDst++
                break
            }
        }
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
