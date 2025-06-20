package dblenc

import (
    "encoding/hex"
    "testing"

    "github.com/stretchr/testify/assert"
)

var (
    asciiShort           = "616C61206D61"
    asciiLong            = "616C61206D61206B6F746120616C6120"
    utf8Only             = "C2A72EE280A2C2B4C2A827C2B0C3B7E280A22E2EC397202020F09F8E802020F09D929EF09D92BDF09D92B6F09D9387F09D92B6F09D92B8F09D9389F09D9192F09D9387F09D93887E2020F09F8E80202020C3972E2EE280A2C3B7C2B027C2A8C2B4E280A22EC2A7"
    utf8DblEnc           = "C3A0C2A6E280A0C3A0C2A6C59320C3A0C2A6C2ACC3A0C2A6C2BFC3A0C2A6C2B6C3A0C2A7C28DC3A0C2A6C2AC20C3A0C2A6C2ADC3A0C2A6C2BEC3A0C2A6C2B2C3A0C2A6C2ACC3A0C2A6C2BEC3A0C2A6C2B8C3A0C2A6C2BE20C3A0C2A6C2A6C3A0C2A6C2BFC3A0C2A6C2ACC3A0C2A6C2B8"
    utf8TrplEnc          = "C383C692C382C2A5C383E2809AC382C2A4C383E2809AC382C2A7C383C692C382C2A7C383E280B9C3A2E282ACC2A0C383C2A2C3A2E2809AC2ACC382C2A0C383C692C382C2A7C383C2A2C3A2E282ACC5BEC382C2A2C383E2809AC382C2BA"
    utf8DblEncTrunc      = "C3A5C2A4C2A7C3A7CB86E280A0C3A7E284A2"

    asciiShortCheck      = "616C61206D61"
    asciiLongCheck       = "616C61206D61206B6F746120616C6120"
    utf8OnlyCheck        = "C2A72EE280A2C2B4C2A827C2B0C3B7E280A22E2EC397202020F09F8E802020F09D929EF09D92BDF09D92B6F09D9387F09D92B6F09D92B8F09D9389F09D9192F09D9387F09D93887E2020F09F8E80202020C3972E2EE280A2C3B7C2B027C2A8C2B4E280A22EC2A7"
    utf8DblEncCheck      = "E0A686E0A69C20E0A6ACE0A6BFE0A6B6E0A78DE0A6AC20E0A6ADE0A6BEE0A6B2E0A6ACE0A6BEE0A6B8E0A6BE20E0A6A6E0A6BFE0A6ACE0A6B8"
    utf8TrplEncCheck     = "E5A4A7E78886E799BA"
    utf8DblEncTruncCheck = "E5A4A7E78886"
)

func decode(s string) []byte {
    if b, err := hex.DecodeString(s); err != nil {
        panic(err)
    } else {
        return b
    }
}

func TestTransform(t *testing.T) {
    tests := []struct {
        Name        string
        Input       []byte
        ReturnValue []byte
        ReturnError error
    }{
        {
            Name:        "short-ascii",
            Input:       decode(asciiShort),
            ReturnValue: decode(asciiShortCheck),
            ReturnError: nil,
        },
        {
            Name:        "utf8-only",
            Input:       decode(utf8Only),
            ReturnValue: []byte(nil),
            ReturnError: ErrInvalid,
        },
        {
            Name:        "double-encoded",
            Input:       decode(utf8DblEnc),
            ReturnValue: decode(utf8DblEncCheck),
            ReturnError: nil,
        },
    }

    und := NewUnDoubleEncoder()

    for _, tc := range tests {
        t.Run(tc.Name, func(t *testing.T) {
            b, err := und.transform(tc.Input)
            assert.ErrorIs(t, err, tc.ReturnError)
            assert.Equal(t, b, tc.ReturnValue)
        })
    }
}

func BenchmarkTransformAsciiShort(b *testing.B) {
    buf, err := hex.DecodeString(asciiShort)
    assert.NoError(b, err)

    und := NewUnDoubleEncoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        und.transform(buf)
    }
}

func BenchmarkTransformAsciiLong(b *testing.B) {
    buf, err := hex.DecodeString(asciiLong)
    assert.NoError(b, err)

    und := NewUnDoubleEncoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        und.transform(buf)
    }
}

func BenchmarkTransformUtf8Only(b *testing.B) {
    buf, err := hex.DecodeString(utf8Only)
    assert.NoError(b, err)

    und := NewUnDoubleEncoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        und.transform(buf)
    }
}

func BenchmarkTransformDoubleEncoded(b *testing.B) {
    buf, err := hex.DecodeString(utf8DblEnc)
    assert.NoError(b, err)

    und := NewUnDoubleEncoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        und.transform(buf)
    }
}

func TestDetect(t *testing.T) {
    type ReturnValue struct {
        Encoding Encoding
        At       int
    }
    tests := []struct {
        Name        string
        Input       []byte
        ReturnValue ReturnValue
    }{
        {
            Name:        "short-ascii",
            Input:       decode(asciiShort),
            ReturnValue: ReturnValue{ ASCII, 6 },
        },
        {
            Name:        "long-ascii",
            Input:       decode(asciiLong),
            ReturnValue: ReturnValue{ ASCII, 16 },
        },
        {
            Name:        "utf8-only",
            Input:       decode(utf8Only),
            ReturnValue: ReturnValue{ UNKNOWN, 26 },
        },
        {
            Name:        "double-encoded",
            Input:       decode(utf8DblEnc),
            ReturnValue: ReturnValue{ DOUBLE_ENCODED, 112 },
        },
    }

    bm := newByteMap()

    for _, tc := range tests {
        t.Run(tc.Name, func(t *testing.T) {
            enc, at := bm.detect(tc.Input)
            assert.Equal(t, tc.ReturnValue.Encoding, enc)
            assert.Equal(t, tc.ReturnValue.At, at)
        })
    }
}

func TestTransformEndToEnd(t *testing.T) {
    tests := []struct {
        Name        string
        Input       []byte
        ReturnValue []byte
    }{
        {
            Name:        "short-ascii",
            Input:       decode(asciiShort),
            ReturnValue: decode(asciiShortCheck),
        },
        {
            Name:        "long-ascii",
            Input:       decode(asciiLong),
            ReturnValue: decode(asciiLongCheck),
        },
        {
            Name:        "utf8-only",
            Input:       decode(utf8Only),
            ReturnValue: decode(utf8OnlyCheck),
        },
        {
            Name:        "double-encoded",
            Input:       decode(utf8DblEnc),
            ReturnValue: decode(utf8DblEncCheck),
        },
        {
            Name:        "triple-encoded",
            Input:       decode(utf8TrplEnc),
            ReturnValue: decode(utf8TrplEncCheck),
        },
        {
            Name:        "double-encoded-truncated",
            Input:       decode(utf8DblEncTrunc),
            ReturnValue: decode(utf8DblEncTruncCheck),
        },
    }

    und := NewUnDoubleEncoder()

    for _, tc := range tests {
        t.Run(tc.Name, func(t *testing.T) {
            b, err := und.Transform(tc.Input)
            assert.NoError(t, err)
            assert.Equal(t, tc.ReturnValue, b)
        })
    }
}

func BenchmarkDetectAsciiShort(b *testing.B) {
    buf, err := hex.DecodeString(asciiShort)
    assert.NoError(b, err)

    bm := newByteMap()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bm.detect(buf)
    }
}

func BenchmarkDetectAsciiLong(b *testing.B) {
    buf, err := hex.DecodeString(asciiLong)
    assert.NoError(b, err)

    bm := newByteMap()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bm.detect(buf)
    }
}

func BenchmarkDetectUtf8Only(b *testing.B) {
    buf, err := hex.DecodeString(utf8Only)
    assert.NoError(b, err)

    bm := newByteMap()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bm.detect(buf)
    }
}

func BenchmarkDetectDoubleEncoded(b *testing.B) {
    buf, err := hex.DecodeString(utf8DblEnc)
    assert.NoError(b, err)

    bm := newByteMap()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bm.detect(buf)
    }
}