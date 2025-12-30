package dblenc

import (
    "encoding/hex"
    "testing"

    "github.com/stretchr/testify/assert"
)

type TestCase struct {
    Name              string
    Original          string
    OriginalHex       []byte
    TestString        string
    TestStringHex     []byte
    TransformHex      []byte
    TransformError    error
    TransformedHex    []byte
    TransformedError  error
    DetectResult      Encoding
    DetectOffset      int
}

var testCases = []TestCase{
    {
        Name:             "ASCII",
        Original:         "Hello!",
        OriginalHex:      []byte("Hello!"),
        TestString:       "Hello!",
        TestStringHex:    []byte("Hello!"),
        TransformHex:     []byte("Hello!"),
        TransformedHex:   []byte("Hello!"),
        TransformedError: ErrNoop,
        DetectResult:     ASCII,
        DetectOffset:     6,
    },
    {
        Name:             "UTF8",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wąż",
        TestStringHex:    decode("77c485c5bc"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   decode("77c485c5bc"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_False_Positive",
        Original:         "Tomáš",
        OriginalHex:      decode("546f6dc3a1c5a1"),
        TestString:       "Tomáš",
        TestStringHex:    decode("546f6dc3a1c5a1"),
        TransformHex:     decode("546f6de19a"),
        TransformedHex:   decode("546f6d"),
        DetectResult:     UNKNOWN,
        DetectOffset:     4,
    },
    {
        Name:             "UTF8_Japanese",
        Original:         "全然分からない",
        OriginalHex:      decode("e585a8e784b6e58886e3818be38289e381aae38184"),
        TestString:       "å…¨ç„¶åˆ†ã‹ã‚‰ãªã„",
        TestStringHex:    decode("c3a5e280a6c2a8c3a7e2809ec2b6c3a5cb86e280a0c3a3c281e280b9c3a3e2809ae280b0c3a3c281c2aac3a3c281e2809e"),
        TransformHex:     decode("e585a8e784b6e58886e3818be38289e381aae38184"),
        TransformedHex:   decode("e585a8e784b6e58886e3818be38289e381aae38184"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Bengali",
        Original:         "আজ বিশ্ব ভালবাসা দিবস",
        OriginalHex:      decode("e0a686e0a69c20e0a6ace0a6bfe0a6b6e0a78de0a6ac20e0a6ade0a6bee0a6b2e0a6ace0a6bee0a6b8e0a6be20e0a6a6e0a6bfe0a6ace0a6b8"),
        TestString:       "à¦†à¦œ à¦¬à¦¿à¦¶à§à¦¬ à¦­à¦¾à¦²à¦¬à¦¾à¦¸à¦¾ à¦¦à¦¿à¦¬à¦¸",
        TestStringHex:    decode("c3a0c2a6e280a0c3a0c2a6c59320c3a0c2a6c2acc3a0c2a6c2bfc3a0c2a6c2b6c3a0c2a7c28dc3a0c2a6c2ac20c3a0c2a6c2adc3a0c2a6c2bec3a0c2a6c2b2c3a0c2a6c2acc3a0c2a6c2bec3a0c2a6c2b8c3a0c2a6c2be20c3a0c2a6c2a6c3a0c2a6c2bfc3a0c2a6c2acc3a0c2a6c2b8"),
        TransformHex:     decode("e0a686e0a69c20e0a6ace0a6bfe0a6b6e0a78de0a6ac20e0a6ade0a6bee0a6b2e0a6ace0a6bee0a6b8e0a6be20e0a6a6e0a6bfe0a6ace0a6b8"),
        TransformedHex:   decode("e0a686e0a69c20e0a6ace0a6bfe0a6b6e0a78de0a6ac20e0a6ade0a6bee0a6b2e0a6ace0a6bee0a6b8e0a6be20e0a6a6e0a6bfe0a6ace0a6b8"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Short_Emoji",
        Original:         "🙂",
        OriginalHex:      decode("f09f9982"),
        TestString:       "ðŸ™‚",
        TestStringHex:    decode("c3b0c5b8e284a2e2809a"),
        TransformHex:     decode("f09f9982"),
        TransformedHex:   decode("f09f9982"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Long_Emoji",
        Original:         "🚵🏻‍♀️",
        OriginalHex:      decode("f09f9ab5f09f8fbbe2808de29980efb88f"),
        TestString:       "ðŸšµðŸ»â€â™€ï¸",
        TestStringHex:    decode("c3b0c5b8c5a1c2b5c3b0c5b8c28fc2bbc3a2e282acc28dc3a2e284a2e282acc3afc2b8c28f"),
        TransformHex:     decode("f09f9ab5f09f8fbbe2808de29980efb88f"),
        TransformedHex:   decode("f09f9ab5f09f8fbbe2808de29980efb88f"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Fancy_Text",
        Original:         "§.•´¨'°÷•..×   🎀  𝒞𝒽𝒶𝓇𝒶𝒸𝓉𝑒𝓇𝓈  🎀   ×..•÷°'¨´•.§",
        OriginalHex:      decode("c2a72ee280a2c2b4c2a827c2b0c3b7e280a22e2ec397202020f09f8e802020f09d929ef09d92bdf09d92b6f09d9387f09d92b6f09d92b8f09d9389f09d9192f09d9387f09d93882020f09f8e80202020c3972e2ee280a2c3b7c2b027c2a8c2b4e280a22ec2a7"),
        TestString:       "Â§.â€¢Â´Â¨'Â°Ã·â€¢..Ã—   ðŸŽ€  ð’žð’½ð’¶ð“‡ð’¶ð’¸ð“‰ð‘’ð“‡ð“ˆ  ðŸŽ€   Ã—..â€¢Ã·Â°'Â¨Â´â€¢.Â§",
        TestStringHex:    decode("c382c2a72ec3a2e282acc2a2c382c2b4c382c2a827c382c2b0c383c2b7c3a2e282acc2a22e2ec383e28094202020c3b0c5b8c5bde282ac2020c3b0c29de28099c5bec3b0c29de28099c2bdc3b0c29de28099c2b6c3b0c29de2809ce280a1c3b0c29de28099c2b6c3b0c29de28099c2b8c3b0c29de2809ce280b0c3b0c29de28098e28099c3b0c29de2809ce280a1c3b0c29de2809ccb862020c3b0c5b8c5bde282ac202020c383e280942e2ec3a2e282acc2a2c383c2b7c382c2b027c382c2a8c382c2b4c3a2e282acc2a22ec382c2a7"),
        TransformHex:     decode("c2a72ee280a2c2b4c2a827c2b0c3b7e280a22e2ec397202020f09f8e802020f09d929ef09d92bdf09d92b6f09d9387f09d92b6f09d92b8f09d9389f09d9192f09d9387f09d93882020f09f8e80202020c3972e2ee280a2c3b7c2b027c2a8c2b4e280a22ec2a7"),
        TransformedHex:   decode("c2a72ee280a2c2b4c2a827c2b0c3b7e280a22e2ec397202020f09f8e802020f09d929ef09d92bdf09d92b6f09d9387f09d92b6f09d92b8f09d9389f09d9192f09d9387f09d93882020f09f8e80202020c3972e2ee280a2c3b7c2b027c2a8c2b4e280a22ec2a7"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Double_Encoded",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÄ…Å¼",
        TestStringHex:    decode("77c384e280a6c385c2bc"),
        TransformHex:     decode("77c485c5bc"),
        TransformedHex:   decode("77c485c5bc"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Double_Encoded_Truncated",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5"),
        TestString:       "wÄ…Å",
        TestStringHex:    decode("77c384e280a6c385"),
        TransformHex:     decode("77c485c5"),
        TransformedHex:   decode("77c485"),
        DetectResult:     DOUBLE_ENCODED_TRUNCATED,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Double_Encoded_Truncated_2",
        Original:         "wąwż",
        OriginalHex:      decode("77c48577c5"),
        TestString:       "wÄ…wÅ",
        TestStringHex:    decode("77c384e280a677c385"),
        TransformHex:     decode("77c48577c5"),
        TransformedHex:   decode("77c48577"),
        DetectResult:     DOUBLE_ENCODED_TRUNCATED,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Double_Encoded_Truncated_3",
        Original:         "waż",
        OriginalHex:      decode("7761c5bc"),
        TestString:       "waÅ",
        TestStringHex:    decode("7761c385"),
        TransformHex:     decode("7761c5"),
        TransformedHex:   decode("7761"),
        DetectResult:     UNKNOWN,
        DetectOffset:     3,
    },
    {
        Name:             "UTF8_Malformed_Suffix_Outer",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÄ…Å�",
        TestStringHex:    decode("77c384e280a6c385c2"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   decode("77c485"),
        DetectResult:     ERROR,
        DetectOffset:     9,
    },
    {
        Name:             "UTF8_Malformed_Suffix_Inner",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÄ…Å",
        TestStringHex:    decode("77c384e280a6c385"),
        TransformHex:     decode("77c485c5"),
        TransformedHex:   decode("77c485"),
        DetectResult:     DOUBLE_ENCODED_TRUNCATED,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Malformed_Infix_Outer",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÄ…�¼",
        TestStringHex:    decode("77c384e280a6c300c2bc"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   decode("77c384e280a6c300c2bc"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     7,
    },
    {
        Name:             "UTF8_Malformed_Infix_Inner",
        Original:         "w�ż",
        OriginalHex:      decode("77c400c5bc"),
        TestString:       "wÄ Å¼",
        TestStringHex:    decode("77c38400c385c2bc"),
        TransformHex:     decode("77c400c5bc"),
        TransformedHex:   decode("77c38400c385c2bc"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     4,
    },
    {
        Name:             "UTF8_Triple_Encoded",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÃ„â€¦Ã…Â¼",
        TestStringHex:    decode("77c383e2809ec3a2e282acc2a6c383e280a6c382c2bc"),
        TransformHex:     decode("77c384e280a6c385c2bc"),
        TransformedHex:   decode("77c485c5bc"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Triple_Encoded_Truncated_Last_Byte",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÃ„â€¦Ã…Â",
        TestStringHex:    decode("77c383e2809ec3a2e282acc2a6c383e280a6c382c2"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   decode("77c485"),
        DetectResult:     ERROR,
        DetectOffset:     21,
    },
    {
        Name:             "UTF8_Triple_Encoded_Truncated_Last_Sequence",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÃ„â€¦Ã…",
        TestStringHex:    decode("77c383e2809ec3a2e282acc2a6c383e280a6c382"),
        TransformHex:     decode("77c384e280a6c385c2"),
        TransformedHex:   decode("77c485"),
        DetectResult:     DOUBLE_ENCODED_TRUNCATED,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Quarduple_Encoded",
        Original:         "wąż",
        OriginalHex:      decode("77c485c5bc"),
        TestString:       "wÃƒâ€žÃ¢â‚¬Â¦Ãƒâ€¦Ã‚Â¼",
        TestStringHex:    decode("77c383c692c3a2e282acc5bec383c2a2c3a2e2809ac2acc382c2a6c383c692c3a2e282acc2a6c383e2809ac382c2bc"),
        TransformHex:     decode("77c383e2809ec3a2e282acc2a6c383e280a6c382c2bc"),
        TransformedHex:   decode("77c485c5bc"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     2,
    },
}

var (
    asciiShort    = decode("20202020202020") // whitespace
    asciiLong     = decode("2020202020202020") // whitespace
    wellEncoded   = decode("20e8a5bfe38282e69db1e38282e58886e3818be38289e381aae38184") // "西も東も分からない"
    doubleEncoded = decode("20c3a8c2a5c2bfc3a3e2809ae2809ac3a6c29dc2b1c3a3e2809ae2809ac3a5cb86e280a0c3a3c281e280b9c3a3e2809ae280b0c3a3c281c2aac3a3c281e2809e") // "西も東も分からない"
)

func decode(s string) []byte {
    if b, err := hex.DecodeString(s); err != nil {
        panic(err)
    } else {
        return b
    }
}

func TestDetect(t *testing.T) {
    d := NewDecoder()

    for _, tc := range testCases {
        t.Run(tc.Name, func(t *testing.T) {
            result, offset := d.Detect(tc.TestStringHex)
            assert.Equal(t, tc.DetectResult, result)
            assert.Equal(t, tc.DetectOffset, offset)
        })
    }
}

func TestTransform(t *testing.T) {
    d := NewDecoder()

    for _, tc := range testCases {
        t.Run(tc.Name, func(t *testing.T) {
            r, err := d.transform(tc.TestStringHex)
            assert.ErrorIs(t, tc.TransformError, err)
            assert.Equal(t, tc.TransformHex, r)
        })
    }
}

func TestTransformed(t *testing.T) {
    d := NewDecoder()

    for _, tc := range testCases {
        t.Run(tc.Name, func(t *testing.T) {
            r, err := d.Transform(tc.TestStringHex)
            assert.ErrorIs(t, tc.TransformedError, err)
            assert.Equal(t, tc.TransformedHex, r)
        })
    }
}

func BenchmarkTransformAsciiShort(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.transform(asciiShort)
    }
}

func BenchmarkTransformAsciiLong(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.transform(asciiLong)
    }
}

func BenchmarkTransformWellEncoded(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.transform(wellEncoded)
    }
}

func BenchmarkTransformDoubleEncoded(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.transform(doubleEncoded)
    }
}

func BenchmarkDetectAsciiShort(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.Detect(asciiShort)
    }
}

func BenchmarkDetectAsciiLong(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.Detect(asciiLong)
    }
}

func BenchmarkDetectWellEncoded(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.Detect(wellEncoded)
    }
}

func BenchmarkDetectDoubleEncoded(b *testing.B) {
    d := NewDecoder()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        d.Detect(doubleEncoded)
    }
}
