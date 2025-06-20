package dblenc

import (
    "errors"
    // "fmt"
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

var encodeMap = [256]uint32{
    0x00000000, 0x01000001, 0x02000002, 0x03000003,
    0x04000004, 0x05000005, 0x06000006, 0x07000007,
    0x08000008, 0x09000009, 0x0a00000a, 0x0b00000b,
    0x0c00000c, 0x0d00000d, 0x0e00000e, 0x0f00000f,
    0x10000010, 0x11000011, 0x12000012, 0x13000013,
    0x14000014, 0x15000015, 0x16000016, 0x17000017,
    0x18000018, 0x19000019, 0x1a00001a, 0x1b00001b,
    0x1c00001c, 0x1d00001d, 0x1e00001e, 0x1f00001f,
    0x20000020, 0x21000021, 0x22000022, 0x23000023,
    0x24000024, 0x25000025, 0x26000026, 0x27000027,
    0x28000028, 0x29000029, 0x2a00002a, 0x2b00002b,
    0x2c00002c, 0x2d00002d, 0x2e00002e, 0x2f00002f,
    0x30000030, 0x31000031, 0x32000032, 0x33000033,
    0x34000034, 0x35000035, 0x36000036, 0x37000037,
    0x38000038, 0x39000039, 0x3a00003a, 0x3b00003b,
    0x3c00003c, 0x3d00003d, 0x3e00003e, 0x3f00003f,
    0x40000040, 0x41000041, 0x42000042, 0x43000043,
    0x44000044, 0x45000045, 0x46000046, 0x47000047,
    0x48000048, 0x49000049, 0x4a00004a, 0x4b00004b,
    0x4c00004c, 0x4d00004d, 0x4e00004e, 0x4f00004f,
    0x50000050, 0x51000051, 0x52000052, 0x53000053,
    0x54000054, 0x55000055, 0x56000056, 0x57000057,
    0x58000058, 0x59000059, 0x5a00005a, 0x5b00005b,
    0x5c00005c, 0x5d00005d, 0x5e00005e, 0x5f00005f,
    0x60000060, 0x61000061, 0x62000062, 0x63000063,
    0x64000064, 0x65000065, 0x66000066, 0x67000067,
    0x68000068, 0x69000069, 0x6a00006a, 0x6b00006b,
    0x6c00006c, 0x6d00006d, 0x6e00006e, 0x6f00006f,
    0x70000070, 0x71000071, 0x72000072, 0x73000073,
    0x74000074, 0x75000075, 0x76000076, 0x77000077,
    0x78000078, 0x79000079, 0x7a00007a, 0x7b00007b,
    0x7c00007c, 0x7d00007d, 0x7e00007e, 0x7f00007f,
    0x81000081, 0x8d00008d, 0x8f00008f, 0x90000090,
    0x9d00009d, 0xa00000a0, 0xa10000a1, 0xa20000a2,
    0xa30000a3, 0xa40000a4, 0xa50000a5, 0xa60000a6,
    0xa70000a7, 0xa80000a8, 0xa90000a9, 0xaa0000aa,
    0xab0000ab, 0xac0000ac, 0xad0000ad, 0xae0000ae,
    0xaf0000af, 0xb00000b0, 0xb10000b1, 0xb20000b2,
    0xb30000b3, 0xb40000b4, 0xb50000b5, 0xb60000b6,
    0xb70000b7, 0xb80000b8, 0xb90000b9, 0xba0000ba,
    0xbb0000bb, 0xbc0000bc, 0xbd0000bd, 0xbe0000be,
    0xbf0000bf, 0xc00000c0, 0xc10000c1, 0xc20000c2,
    0xc30000c3, 0xc40000c4, 0xc50000c5, 0xc60000c6,
    0xc70000c7, 0xc80000c8, 0xc90000c9, 0xca0000ca,
    0xcb0000cb, 0xcc0000cc, 0xcd0000cd, 0xce0000ce,
    0xcf0000cf, 0xd00000d0, 0xd10000d1, 0xd20000d2,
    0xd30000d3, 0xd40000d4, 0xd50000d5, 0xd60000d6,
    0xd70000d7, 0xd80000d8, 0xd90000d9, 0xda0000da,
    0xdb0000db, 0xdc0000dc, 0xdd0000dd, 0xde0000de,
    0xdf0000df, 0xe00000e0, 0xe10000e1, 0xe20000e2,
    0xe30000e3, 0xe40000e4, 0xe50000e5, 0xe60000e6,
    0xe70000e7, 0xe80000e8, 0xe90000e9, 0xea0000ea,
    0xeb0000eb, 0xec0000ec, 0xed0000ed, 0xee0000ee,
    0xef0000ef, 0xf00000f0, 0xf10000f1, 0xf20000f2,
    0xf30000f3, 0xf40000f4, 0xf50000f5, 0xf60000f6,
    0xf70000f7, 0xf80000f8, 0xf90000f9, 0xfa0000fa,
    0xfb0000fb, 0xfc0000fc, 0xfd0000fd, 0xfe0000fe,
    0xff0000ff, 0x8c000152, 0x9c000153, 0x8a000160,
    0x9a000161, 0x9f000178, 0x8e00017d, 0x9e00017e,
    0x83000192, 0x880002c6, 0x980002dc, 0x96002013,
    0x97002014, 0x91002018, 0x92002019, 0x8200201a,
    0x9300201c, 0x9400201d, 0x8400201e, 0x86002020,
    0x87002021, 0x95002022, 0x85002026, 0x89002030,
    0x8b002039, 0x9b00203a, 0x800020ac, 0x99002122,
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
        // buf = buf[:0]
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
    bm *byteMap
}

func NewUnDoubleEncoder() *UnDoubleEncoder {
    return &UnDoubleEncoder{
        bm: newByteMap(),
    }
}

func (this *UnDoubleEncoder) Transform(b []byte) ([]byte, error) {
    // {
    //     enc, at := this.Detect(b)
    //     fmt.Printf(
    //         "INPUT:\nstr   = %s\nhex   = %x\ntest  = %v [at %d/%d]\n\n",
    //         b, b[:], enc, at, len(b),
    //     )
    // }

    enc, _ := this.bm.detect(b)
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
        enc, _ = this.bm.detect(x)
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
            got := encodeMap[mid]
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
