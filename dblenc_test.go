package dblenc

import (
    "encoding/hex"
    "testing"

    "github.com/stretchr/testify/assert"
)

type TestCase struct {
    Name              string
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
        Name:             "Simple_ASCII",
        TestString:       "Hello world!",
        TestStringHex:    []byte("Hello world!"),
        TransformHex:     []byte("Hello world!"),
        TransformedHex:   []byte("Hello world!"),
        TransformedError: ErrNoop,
        DetectResult:     ASCII,
        DetectOffset:     12,
    },
    {
        Name:             "UTF8_Polish",
        TestString:       "Zażółć gęślą jaźń",
        TestStringHex:    []byte("Zażółć gęślą jaźń"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Zażółć gęślą jaźń"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     3,
    },
    {
        Name:             "UTF8_French",
        TestString:       "Héllo çà va très bien",
        TestStringHex:    []byte("Héllo çà va très bien"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Héllo çà va très bien"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     4,
    },
    {
        Name:             "UTF8_German",
        TestString:       "Guten Tag, wie geht's? Größe, Ärger, Übung",
        TestStringHex:    []byte("Guten Tag, wie geht's? Größe, Ärger, Übung"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Guten Tag, wie geht's? Größe, Ärger, Übung"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     27,
    },
    {
        Name:             "UTF8_Italian",
        TestString:       "Città, perché, caffè",
        TestStringHex:    []byte("Città, perché, caffè"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Città, perché, caffè"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     7,
    },
    {
        Name:             "UTF8_Portuguese",
        TestString:       "São Paulo, ação, coração",
        TestStringHex:    []byte("São Paulo, ação, coração"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("São Paulo, ação, coração"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     4,
    },
    {
        Name:             "UTF8_Czech",
        TestString:       "Příliš žluťoučký kůň",
        TestStringHex:    []byte("Příliš žluťoučký kůň"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Příliš žluťoučký kůň"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Ukrainian",
        TestString:       "Доброго дня, Україна",
        TestStringHex:    []byte("Доброго дня, Україна"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Доброго дня, Україна"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Chinese",
        TestString:       "你好世界",
        TestStringHex:    []byte("你好世界"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("你好世界"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Japanese",
        TestString:       "こんにちは世界、カタカナもあります",
        TestStringHex:    []byte("こんにちは世界、カタカナもあります"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("こんにちは世界、カタカナもあります"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Korean",
        TestString:       "안녕하세요 세계",
        TestStringHex:    []byte("안녕하세요 세계"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("안녕하세요 세계"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Thai",
        TestString:       "สวัสดีครับ",
        TestStringHex:    []byte("สวัสดีครับ"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("สวัสดีครับ"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Vietnamese",
        TestString:       "Xin chào thế giới",
        TestStringHex:    []byte("Xin chào thế giới"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Xin chào thế giới"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     9,
    },
    {
        Name:             "UTF8_Arabic",
        TestString:       "مرحبا بالعالم",
        TestStringHex:    []byte("مرحبا بالعالم"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("مرحبا بالعالم"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Hebrew",
        TestString:       "שלום עולם",
        TestStringHex:    []byte("שלום עולם"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("שלום עולם"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Hindi",
        TestString:       "नमस्ते दुनिया",
        TestStringHex:    []byte("नमस्ते दुनिया"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("नमस्ते दुनिया"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Greek",
        TestString:       "Γεια σου κόσμε",
        TestStringHex:    []byte("Γεια σου κόσμε"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Γεια σου κόσμε"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Armenian",
        TestString:       "Բարեւ աշխարհ",
        TestStringHex:    []byte("Բարեւ աշխարհ"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Բարեւ աշխարհ"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Single_Emoji",
        TestString:       "😀",
        TestStringHex:    []byte("😀"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("😀"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Multiple_Emojis",
        TestString:       "😀😃😄😁",
        TestStringHex:    []byte("😀😃😄😁"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("😀😃😄😁"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Skin_Tone_Emoji",
        TestString:       "👍🏻👍🏼👍🏽👍🏾👍🏿",
        TestStringHex:    []byte("👍🏻👍🏼👍🏽👍🏾👍🏿"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("👍🏻👍🏼👍🏽👍🏾👍🏿"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Complex_Emoji",
        TestString:       "👨‍👩‍👧‍👦👨‍💻",
        TestStringHex:    []byte("👨‍👩‍👧‍👦👨‍💻"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("👨‍👩‍👧‍👦👨‍💻"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Emoji_with_Text",
        TestString:       "Hello 👋 World 🌍",
        TestStringHex:    []byte("Hello 👋 World 🌍"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Hello 👋 World 🌍"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     7,
    },
    {
        Name:             "Double_Encoded_Letter",
        TestString:       "é",
        TestStringHex:    []byte("\xC3\x83\xC2\xA9"),
        TransformHex:     []byte("é"),
        TransformedHex:   []byte("é"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Letters",
        TestString:       "éé",
        TestStringHex:    []byte("\xC3\x83\xC2\xA9\xC3\x83\xC2\xA9"),
        TransformHex:     []byte("éé"),
        TransformedHex:   []byte("éé"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Kanji",
        TestString:       "中",
        TestStringHex:    []byte("\xC3\xA4\xC2\xB8\xC2\xAD"),
        TransformHex:     []byte("中"),
        TransformedHex:   []byte("中"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Two_Kanji",
        TestString:       "日本",
        TestStringHex:    []byte("\xC3\xA6\xE2\x80\x94\xC2\xA5\xC3\xA6\xC5\x93\xC2\xAC"),
        TransformHex:     []byte("日本"),
        TransformedHex:   []byte("日本"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Emoji",
        TestString:       "😀",
        TestStringHex:    []byte("\xC3\xB0\xC5\xB8\xCB\x9C\xE2\x82\xAC"),
        TransformHex:     []byte("😀"),
        TransformedHex:   []byte("😀"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Complex_Emoji",
        TestString:       "👨‍👩‍👧‍👦",
        TestStringHex:    []byte("\xC3\xB0\xC5\xB8\xE2\x80\x98\xC2\xA8\xC3\xA2\xE2\x82\xAC\xC2\x8D\xC3\xB0\xC5\xB8\xE2\x80\x98\xC2\xA9\xC3\xA2\xE2\x82\xAC\xC2\x8D\xC3\xB0\xC5\xB8\xE2\x80\x98\xC2\xA7\xC3\xA2\xE2\x82\xAC\xC2\x8D\xC3\xB0\xC5\xB8\xE2\x80\x98\xC2\xA6"),
        TransformHex:     []byte("👨‍👩‍👧‍👦"),
        TransformedHex:   []byte("👨‍👩‍👧‍👦"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Polish",
        TestString:       "Zażółć gęślą jaźń",
        TestStringHex:    []byte("Za\xC3\x85\xC2\xBC\xC3\x83\xC2\xB3\xC3\x85\xE2\x80\x9A\xC3\x84\xE2\x80\xA1 g\xC3\x84\xE2\x84\xA2\xC3\x85\xE2\x80\xBAl\xC3\x84\xE2\x80\xA6 ja\xC3\x85\xC2\xBA\xC3\x85\xE2\x80\x9E"),
        TransformHex:     []byte("Zażółć gęślą jaźń"),
        TransformedHex:   []byte("Zażółć gęślą jaźń"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     3,
    },
    {
        Name:             "Double_Encoded_French",
        TestString:       "Héllo çà va très bien",
        TestStringHex:    []byte("H\xC3\x83\xC2\xA9llo \xC3\x83\xC2\xA7\xC3\x83\xC2\xA0 va tr\xC3\x83\xC2\xA8s bien"),
        TransformHex:     []byte("Héllo çà va très bien"),
        TransformedHex:   []byte("Héllo çà va très bien"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     2,
    },
    {
        Name:             "Double_Encoded_German",
        TestString:       "Guten Tag, wie geht's? Größe, Ärger, Übung",
        TestStringHex:    []byte("Guten Tag, wie geht's? Gr\xC3\x83\xC2\xB6\xC3\x83\xC5\xB8e, \xC3\x83\xE2\x80\x9Erger, \xC3\x83\xC5\x93bung"),
        TransformHex:     []byte("Guten Tag, wie geht's? Größe, Ärger, Übung"),
        TransformedHex:   []byte("Guten Tag, wie geht's? Größe, Ärger, Übung"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     26,
    },
    {
        Name:             "Double_Encoded_Italian",
        TestString:       "Città, perché, caffè",
        TestStringHex:    []byte("Citt\xC3\x83\xC2\xA0, perch\xC3\x83\xC2\xA9, caff\xC3\x83\xC2\xA8"),
        TransformHex:     []byte("Città, perché, caffè"),
        TransformedHex:   []byte("Città, perché, caffè"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     5,
    },
    {
        Name:             "Double_Encoded_Portuguese",
        TestString:       "São Paulo, ação, coração",
        TestStringHex:    []byte("S\xC3\x83\xC2\xA3o Paulo, a\xC3\x83\xC2\xA7\xC3\x83\xC2\xA3o, cora\xC3\x83\xC2\xA7\xC3\x83\xC2\xA3o"),
        TransformHex:     []byte("São Paulo, ação, coração"),
        TransformedHex:   []byte("São Paulo, ação, coração"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     2,
    },
    {
        Name:             "Double_Encoded_Czech",
        TestString:       "Příliš žluťoučký kůň",
        TestStringHex:    []byte("P\xC3\x85\xE2\x84\xA2\xC3\x83\xC2\xADli\xC3\x85\xC2\xA1 \xC3\x85\xC2\xBElu\xC3\x85\xC2\xA5ou\xC3\x84\xC2\x8Dk\xC3\x83\xC2\xBD k\xC3\x85\xC2\xAF\xC3\x85\xCB\x86"),
        TransformHex:     []byte("Příliš žluťoučký kůň"),
        TransformedHex:   []byte("Příliš žluťoučký kůň"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     2,
    },
    {
        Name:             "Double_Encoded_Ukrainian",
        TestString:       "Доброго дня, Україна",
        TestStringHex:    decode("C390E2809DC390C2BEC390C2B1C391E282ACC390C2BEC390C2B3C390C2BE20C390C2B4C390C2BDC391C28F2C20C390C2A3C390C2BAC391E282ACC390C2B0C391E28094C390C2BDC390C2B0"),
        TransformHex:     []byte("Доброго дня, Україна"),
        TransformedHex:   []byte("Доброго дня, Україна"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Chinese",
        TestString:       "你好世界",
        TestStringHex:    decode("C3A4C2BDC2A0C3A5C2A5C2BDC3A4C2B8E28093C3A7E280A2C592"),
        TransformHex:     []byte("你好世界"),
        TransformedHex:   []byte("你好世界"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Japanese",
        TestString:       "こんにちは世界、カタカナもあります",
        TestStringHex:    decode("C3A3C281E2809CC3A3E2809AE2809CC3A3C281C2ABC3A3C281C2A1C3A3C281C2AFC3A4C2B8E28093C3A7E280A2C592C3A3E282ACC281C3A3E2809AC2ABC3A3E2809AC2BFC3A3E2809AC2ABC3A3C692C5A0C3A3E2809AE2809AC3A3C281E2809AC3A3E2809AC5A0C3A3C281C2BEC3A3C281E284A2"),
        TransformHex:     []byte("こんにちは世界、カタカナもあります"),
        TransformedHex:   []byte("こんにちは世界、カタカナもあります"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Japanese_Long",
        TestString:       "...",
        TestStringHex:    decode("C3A4C2BBC5A0C3A3C281C2B3C3A8C2A1C592C3A9C281C2B8C3A3C692C2B1C3A3C692C5BEC3A6E280B0E280A2C3A5C2BDC2B0C3A3E2809AC2BDC3A8C2A1C5923831C3A8C2BCE280B0C3A3E2809AE280A0C3A3C281C2AEC3A3E2809AE280A1C3A3C281C593C3A5C5BDC5B8C3A8C2A6E280B9C3A3E2809AE2809EC3A5C2A4C5A1C3A6C5A0E280A2C3A3C692CB86C3A3C281C2A4C3A7C2B6C5A1C3A8C2ACE280BAC3A3C692C2ABC3A3C692C281C3A3E2809AC2B9C3A3C692C5BDC3A5C281C2B4C3A9E28093E280B9C3A3C281C2A1C3A3C281C2A6C3A3E2809ACB86C3A5E280A1C2BAC3A5C2B0E280A0C3A3C692C281C3A3C692E28099C3A3C692C2ADC3A5E280B9C29D37C3A6E280B0C2B9C3A3C692E28099C3A3E2809AC2AFC3A4C2BAE2809DC3A8C2AAC28DC3A6E284A2E2809AC3A6C2AEE280B9C3A3E2809AE280B9C3A3C281C2B1C3A3C281E280BAC3A5C2BCCB9CC3A8C2A6E280B9C3A3C281C2A0C3A3C281C5BDC3A3C281E28099C3A3C692E280B0C3A6C2BCE2809DC3A5C2A0C2B1C3A3C692C2A1C3A7C2ADE280933439C3A5E280B9E284A2C3A8E28094C2A4C3A5E2809CC281C3A5E280B0C2A4C3A7C5A1C2AEC3A3C281C2A5C3A3E282ACE2809AC3A8C2AAC2B2C3A3C692C28DC3A3E2809AC2B5C3A3C692E280A2C3A7C2A4C2BAC3A4C2BCC2BCC3A3E2809AC2B1C3A3C692C2B1C3A7E28099C2B0C3A8C2B2C2A9C3A3C281C2B1C3A3C281E28098C3A3C281C2A5C3A5E280A6E280A0C3A8E280B0C2AFC3A5C2B8C2B0C3A3C281C2A7C3A3C281C2AEC3A5C290E280BAC3A5C2A0C2B1C3A3E2809AC2AAC3A5E280B0C2AFC3A8C2A8CB86C3A3C692C2ACC3A3E2809AC2ABC3A3E2809AC2B7C3A6CB86C2A6C3A7C2B4C2A2C3A3E2809AE280B0C3A7E2809DC2B7C3A5C2A9C5A1C3A6C2B0E28094C3A3E2809AC2A4C3A3C281C2B9C3A7E280BAC2B4C3A5C2B0C28FC3A3E2809AE2809AC3A3C692C2B3C3A3E2809AC592C3A3C281C28DC3A6E28093C2ADC3A7C2A0E2809DC3A3C281E2809EC3A3C692E280A2C3A3E2809AE280A6C3A3C281C2BCC3A5C2B1E280B9C3A4C2B8C5A0C3A3C281E280BAC3A3C281C2B4C3A3C692C2ABC3A3C692E280A2C3A5E280A6C2B5C3A6C2B5C281C3A3C692C2A2C3A6C2A8C2A9C3A4C2BCC5A1C3A6C5BDC2B2C3A3C692C2B2C3A9C2A0E28098C3A5E2809AC2ACC3A6E280B9E2809CC3A3C281C2A3C3A3E2809AE280B0C3A3C281E284A2C3A6E280A2E28094C3A8C2AAE280A1C3A8C2AAC2A4C3A3C281C28DC3A3C281E28098C3A3E282ACE2809AC3A6C593C2ACC3A3C281C2A0C3A3C281CB9CC3A3E2809AE2809EC3A3C281C2A1C3A7C2B7C5A1C3A6E280BAC2B4C3A3E2809AC2B1C3A3C692C2A2C3A3C692C5BEC3A6E280BAC2B8C3A5C2B9C2B3C3A3E2809AC2ADC3A3E2809AC2BDC3A3C692E280B9C3A3C692C5A0C3A5C2ADC2A6C3A5E280A1C2BAC3A3C692C2B1C3A3C692C2B2C3A3C692E280A0C3A3C692C5A0C3A7C2B6C2B8C3A6E284A2E2809AC3A6C290C2BAC3A3C692C2B2C3A3C692C592C3A3E2809AC2BBC3A3C692C2A6C3A7C2A9C28DC3A5E280A6C2A8C3A3C281C5BDC3A3C281C2A0C3A8C2BFE28098C3A6CB86C290C3A9C5A0C2ADC3A7C2A4C2BAC3A3C281C5BDC3A3E2809AC2A4C3A9C2A1C592C3A8C2ABE28093C3A3C281C28FC3A3C281C2BEC3A3C281C2A9C3A6C5BDC2B2C3A7E280BAC2AEC3A3C281C2BDC3A5C290CB863237C3A7C2B4C5A1C3A4C2B8E282ACC3A7C593C592C3A4C2B9C5BEC3A4C2BCC28DC3A4C2BDE280BAC3A3C281C2B3C3A3E282ACE2809AC3A6C28FC290C3A3C281C2A3C3A3C281C290C3A3E2809AC281C3A3E2809AC592C3A8C2BFC2ABC3A5C2B1C2B1C3A3C692C5A0C3A3C692C2A0C3A9E280A2C2B73531C3A6CB86E28099C3A3C692C5BDC3A3C692C2A4C3A3C692E280A0C3A5C28FC5BDC3A6E28093C2ADC3A3C692C2ACC3A3E2809AC2B5C3A3E2809AC2BBC3A5C2B1C2B13935C3A8C2B0C2B7C3A3C281C28FC3A3E2809AE280A1C3A3C692E280A2C3A3C281E284A2C3A8C2A1C2A8C3A6E280B0C2BFC3A3C692C2ACC3A3E2809AC2B3C3A3C692E2809EC3A6E28093C2B0C3A5C2A0E282ACC3A6C2B3C2A8C3A3C281C2AFC3A6E28093C2ADC3A5C2BFE280A6C3A3C281E280A2C3A3C281C2A7C3A3C281C29DC3A6C28FC28FC3A5C2AEC5A1C3A3C692CB86C3A3E2809AC2AAC3A3C692C5A0C3A3E2809AC2BFC3A6C5BDC2AAC3A6E282ACC2A7C3A3C281C2AFC3A3C281C290C3A4C2BCC29DC3A6E280B9C2A0C3A3C692C28DC3A3C692C592C3A3C692C5BEC3A3C692C2A4C3A5C2BFC593C3A5C5BDC2BBC3A3C281E28098C3A3E282ACE2809AC3A7C5A0C2ACC3A3C692CB86C3A3E2809AC2ABC3A3C692E280B9C3A7E280B0C2A9C3A8C2B3E282ACC3A3E2809AC2BFC3A3E2809AC2BBC3A5C2BCC2B7C3A8C2A8C2B3C3A6C592C2AFC3A6C2AFE280BAC3A3C692C5BDC3A3C692C2A0C3A3E2809AC2B5C3A8C2B7C2AFC3A5C28FC2B3C3A3C281C2BEC3A3C281C2AEC3A3E2809AC281C3A3C281C2B3C3A6E28099C2AEC3A6C2A1CB86C3A3C281C29DC3A5C2ADC2A6C3A5C2B1E282ACC3A3C281C2BCC3A3C692E280B0C3A6C2B5C593C3A6E284A2E2809AC3A3E2809AE280A6C3A5C28DE28094C3A7E2809DC2BAC3A3E2809AC5A0C3A3C281C2B2C3A3C281C2B5C3A5C2AEC2B6C3A9E28093E2809CC3A3E2809AC2ADC3A8C2A1E280943933C3A6C2B3C281C3A3E2809AE2809EC3A3E2809AC2A4C3A3C281C2BEC3A8C2A8E280943837C3A8C2A1C2A8C3A4C2B9C2B3C3A3C281E28099C3A3C281C2AEC3A9E282ACC5A1C3A5E28098CB86C3A3C692E280B9C3A3C692C2ABC3A6C2B5C281C3A5C2A4C2B1C3A3C281E2809DC3A3C281E280A2C3A3C281C2A0C3A3C692E280A2C3A3E282ACE2809A"),
        TransformHex:     []byte("今び行選ヱマ払彰ソ行81載ゆのょぜ原見や多投トつ続講ルチスノ側開ちてよ出将チヒロ勝7批ヒク五認時残るぱせ弘見だぎげド演報メ策49務藤品剤皮づ。課ネサフ示似ケヱ環販ぱけづ兆良帰での君報オ副計レカシ戦索ら男婚気イべ直小もンれき断研いフゅぼ屋上せぴルフ兵流モ権会掲ヲ頑催拓っらす敗誇誤きけ。本だじやち線更ケモマ書平キソニナ学出ヱヲテナ綸時携ヲヌセユ積全ぎだ近成銭示ぎイ題論くまど掲目ぽ合27級一県乞伍佛び。提っぐめれ迫山ナム長51戒ノヤテ収断レサセ山95谷くょフす表承レコツ新堀注は断必さでそ描定トオナタ措性はぐ伝拠ネヌマヤ応去け。犬トカニ物賀タセ強訳振毛ノムサ路右まのめび撮案そ学局ぼド浜時ゅ南町りひふ家間キ街93況やイま託87表乳げの通呈ニル流失ごさだフ。"),
        TransformedHex:   []byte("今び行選ヱマ払彰ソ行81載ゆのょぜ原見や多投トつ続講ルチスノ側開ちてよ出将チヒロ勝7批ヒク五認時残るぱせ弘見だぎげド演報メ策49務藤品剤皮づ。課ネサフ示似ケヱ環販ぱけづ兆良帰での君報オ副計レカシ戦索ら男婚気イべ直小もンれき断研いフゅぼ屋上せぴルフ兵流モ権会掲ヲ頑催拓っらす敗誇誤きけ。本だじやち線更ケモマ書平キソニナ学出ヱヲテナ綸時携ヲヌセユ積全ぎだ近成銭示ぎイ題論くまど掲目ぽ合27級一県乞伍佛び。提っぐめれ迫山ナム長51戒ノヤテ収断レサセ山95谷くょフす表承レコツ新堀注は断必さでそ描定トオナタ措性はぐ伝拠ネヌマヤ応去け。犬トカニ物賀タセ強訳振毛ノムサ路右まのめび撮案そ学局ぼド浜時ゅ南町りひふ家間キ街93況やイま託87表乳げの通呈ニル流失ごさだフ。"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Korean",
        TestString:       "안녕하세요 세계",
        TestStringHex:    decode("C3ACE280A2CB86C3ABE280A6E280A2C3ADE280A2CB9CC3ACE2809EC2B8C3ACC5A1E2809D20C3ACE2809EC2B8C3AAC2B3E2809E"),
        TransformHex:     []byte("안녕하세요 세계"),
        TransformedHex:   []byte("안녕하세요 세계"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Thai",
        TestString:       "สวัสดีครับ",
        TestStringHex:    decode("C3A0C2B8C2AAC3A0C2B8C2A7C3A0C2B8C2B1C3A0C2B8C2AAC3A0C2B8E2809DC3A0C2B8C2B5C3A0C2B8E2809EC3A0C2B8C2A3C3A0C2B8C2B1C3A0C2B8C5A1"),
        TransformHex:     []byte("สวัสดีครับ"),
        TransformedHex:   []byte("สวัสดีครับ"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Vietnamese",
        TestString:       "Xin chào thế giới",
        TestStringHex:    []byte("Xin ch\xC3\x83\xC2\xA0o th\xC3\xA1\xC2\xBA\xC2\xBF gi\xC3\xA1\xC2\xBB\xE2\x80\xBAi"),
        TransformHex:     []byte("Xin chào thế giới"),
        TransformedHex:   []byte("Xin chào thế giới"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     7,
    },
    {
        Name:             "Double_Encoded_Arabic",
        TestString:       "مرحبا بالعالم",
        TestStringHex:    decode("C399E280A6C398C2B1C398C2ADC398C2A8C398C2A720C398C2A8C398C2A7C399E2809EC398C2B9C398C2A7C399E2809EC399E280A6"),
        TransformHex:     []byte("مرحبا بالعالم"),
        TransformedHex:   []byte("مرحبا بالعالم"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Hebrew",
        TestString:       "שלום עולם",
        TestStringHex:    decode("C397C2A9C397C593C397E280A2C397C29D20C397C2A2C397E280A2C397C593C397C29D"),
        TransformHex:     []byte("שלום עולם"),
        TransformedHex:   []byte("שלום עולם"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Hindi",
        TestString:       "नमस्ते दुनिया",
        TestStringHex:    decode("C3A0C2A4C2A8C3A0C2A4C2AEC3A0C2A4C2B8C3A0C2A5C28DC3A0C2A4C2A4C3A0C2A5E280A120C3A0C2A4C2A6C3A0C2A5C281C3A0C2A4C2A8C3A0C2A4C2BFC3A0C2A4C2AFC3A0C2A4C2BE"),
        TransformHex:     []byte("नमस्ते दुनिया"),
        TransformedHex:   []byte("नमस्ते दुनिया"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Greek",
        TestString:       "Γεια σου κόσμε",
        TestStringHex:    decode("C38EE2809CC38EC2B5C38EC2B9C38EC2B120C38FC692C38EC2BFC38FE280A620C38EC2BAC38FC592C38FC692C38EC2BCC38EC2B5"),
        TransformHex:     []byte("Γεια σου κόσμε"),
        TransformedHex:   []byte("Γεια σου κόσμε"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Armenian",
        TestString:       "Բարեւ աշխարհ",
        TestStringHex:    decode("C394C2B2C395C2A1C396E282ACC395C2A5C396E2809A20C395C2A1C395C2B7C395C2ADC395C2A1C396E282ACC395C2B0"),
        TransformHex:     []byte("Բարեւ աշխարհ"),
        TransformedHex:   []byte("Բարեւ աշխարհ"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Mixed_Language",
        TestString:       "Caféでコーヒーを飲む",
        TestStringHex:    []byte("Caf\xC3\x83\xC2\xA9\xC3\xA3\xC2\x81\xC2\xA7\xC3\xA3\xE2\x80\x9A\xC2\xB3\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xC6\x92\xE2\x80\x99\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xE2\x80\x9A\xE2\x80\x99\xC3\xA9\xC2\xA3\xC2\xB2\xC3\xA3\xE2\x80\x9A\xE2\x82\xAC"),
        TransformHex:     []byte("Caféでコーヒーを飲む"),
        TransformedHex:   []byte("Caféでコーヒーを飲む"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     4,
    },
    {
        Name:             "Triple_Encoded_Letter",
        TestString:       "é",
        TestStringHex:    decode("C383C692C382C2A9"),
        TransformHex:     []byte("Ã©"),
        TransformedHex:   []byte("é"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Triple_Encoded_Letters",
        TestString:       "éé",
        TestStringHex:    decode("C383C692C382C2A9C383C692C382C2A9"),
        TransformHex:     []byte("Ã©Ã©"),
        TransformedHex:   []byte("éé"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Triple_Encoded_Mixed_Language",
        TestString:       "Caféでコーヒーを飲む",
        TestStringHex:    []byte("Caf\xC3\x83\xC6\x92\xC3\x82\xC2\xA9\xC3\x83\xC2\xA3\xC3\x82\xC2\x81\xC3\x82\xC2\xA7\xC3\x83\xC2\xA3\xC3\xA2\xE2\x82\xAC\xC5\xA1\xC3\x82\xC2\xB3\xC3\x83\xC2\xA3\xC3\x86\xE2\x80\x99\xC3\x82\xC2\xBC\xC3\x83\xC2\xA3\xC3\x86\xE2\x80\x99\xC3\xA2\xE2\x82\xAC\xE2\x84\xA2\xC3\x83\xC2\xA3\xC3\x86\xE2\x80\x99\xC3\x82\xC2\xBC\xC3\x83\xC2\xA3\xC3\xA2\xE2\x82\xAC\xC5\xA1\xC3\xA2\xE2\x82\xAC\xE2\x84\xA2\xC3\x83\xC2\xA9\xC3\x82\xC2\xA3\xC3\x82\xC2\xB2\xC3\x83\xC2\xA3\xC3\xA2\xE2\x82\xAC\xC5\xA1\xC3\xA2\xE2\x80\x9A\xC2\xAC"),
        TransformHex:     []byte("Caf\xC3\x83\xC2\xA9\xC3\xA3\xC2\x81\xC2\xA7\xC3\xA3\xE2\x80\x9A\xC2\xB3\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xC6\x92\xE2\x80\x99\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xE2\x80\x9A\xE2\x80\x99\xC3\xA9\xC2\xA3\xC2\xB2\xC3\xA3\xE2\x80\x9A\xE2\x82\xAC"),
        TransformedHex:   []byte("Caféでコーヒーを飲む"),
        DetectResult:     DOUBLE_ENCODED,
        DetectOffset:     4,
    },
    {
        Name:             "Double_Encoded_Letter_Truncated_Byte",
        TestString:       "�",
        TestStringHex:    []byte("\xC3\x83\xC2"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("\xC3\x83\xC2"),
        TransformedError: ErrNoop,
        DetectResult:     ERROR,
        DetectOffset:     3,
    },
    {
        Name:             "Double_Encoded_Two_Letters_Truncated_Byte",
        TestString:       "é�",
        TestStringHex:    []byte("\xC3\x83\xC2\xA9\xC3\x83\xC2"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("é"),
        DetectResult:     ERROR,
        DetectOffset:     7,
    },
    {
        Name:             "Double_Encoded_Mixed_Language_Truncated_Byte",
        TestString:       "Caféでコーヒーを飲�",
        TestStringHex:    []byte("Caf\xC3\x83\xC2\xA9\xC3\xA3\xC2\x81\xC2\xA7\xC3\xA3\xE2\x80\x9A\xC2\xB3\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xC6\x92\xE2\x80\x99\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xE2\x80\x9A\xE2\x80\x99\xC3\xA9\xC2\xA3\xC2\xB2\xC3\xA3\xE2\x80\x9A\xE2\x82"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Caféでコーヒーを飲"),
        DetectResult:     ERROR,
        DetectOffset:     60,
    },
    {
        Name:             "Double_Encoded_Two_Letters_Truncated_Rune",
        TestString:       "é�",
        TestStringHex:    []byte("\xC3\x83\xC2\xA9\xC3\x83"),
        TransformHex:     []byte("\xC3\xA9\xC3"),
        TransformedHex:   []byte("é"),
        DetectResult:     DOUBLE_ENCODED_TRUNCATED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Mixed_Language_Truncated_Rune",
        TestString:       "Caféでコーヒーを飲�",
        TestStringHex:    []byte("Caf\xC3\x83\xC2\xA9\xC3\xA3\xC2\x81\xC2\xA7\xC3\xA3\xE2\x80\x9A\xC2\xB3\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xC6\x92\xE2\x80\x99\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xE2\x80\x9A\xE2\x80\x99\xC3\xA9\xC2\xA3\xC2\xB2\xC3\xA3\xE2\x80\x9A"),
        TransformHex:     []byte("Caf\xC3\xA9\xE3\x81\xA7\xE3\x82\xB3\xE3\x83\xBC\xE3\x83\x92\xE3\x83\xBC\xE3\x82\x92\xE9\xA3\xB2\xE3\x82"),
        TransformedHex:   []byte("Caféでコーヒーを飲"),
        DetectResult:     DOUBLE_ENCODED_TRUNCATED,
        DetectOffset:     4,
    },
    {
        Name:             "Triple_Encoded_Mixed_Language_Truncated_Rune",
        TestString:       "Caféでコーヒーを飲�",
        TestStringHex:    []byte("Caf\xC3\x83\xC6\x92\xC3\x82\xC2\xA9\xC3\x83\xC2\xA3\xC3\x82\xC2\x81\xC3\x82\xC2\xA7\xC3\x83\xC2\xA3\xC3\xA2\xE2\x82\xAC\xC5\xA1\xC3\x82\xC2\xB3\xC3\x83\xC2\xA3\xC3\x86\xE2\x80\x99\xC3\x82\xC2\xBC\xC3\x83\xC2\xA3\xC3\x86\xE2\x80\x99\xC3\xA2\xE2\x82\xAC\xE2\x84\xA2\xC3\x83\xC2\xA3\xC3\x86\xE2\x80\x99\xC3\x82\xC2\xBC\xC3\x83\xC2\xA3\xC3\xA2\xE2\x82\xAC\xC5\xA1\xC3\xA2\xE2\x82\xAC\xE2\x84\xA2\xC3\x83\xC2\xA9\xC3\x82\xC2\xA3\xC3\x82\xC2\xB2\xC3\x83\xC2\xA3\xC3\xA2\xE2\x82\xAC\xC5\xA1\xC3\xA2\xE2\x80\x9A"),
        TransformHex:     []byte("Caf\xC3\x83\xC2\xA9\xC3\xA3\xC2\x81\xC2\xA7\xC3\xA3\xE2\x80\x9A\xC2\xB3\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xC6\x92\xE2\x80\x99\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xE2\x80\x9A\xE2\x80\x99\xC3\xA9\xC2\xA3\xC2\xB2\xC3\xA3\xE2\x80\x9A\xE2\x82"),
        TransformedHex:   []byte("Caféでコーヒーを飲"),
        DetectResult:     DOUBLE_ENCODED_TRUNCATED,
        DetectOffset:     4,
    },
    {
        Name:             "Double_Encoded_Mixed_Language_Irrecoverable",
        TestString:       "Caféでコーヒーを�む",
        TestStringHex:    []byte("Caf\xC3\x83\xC2\xA9\xC3\xA3\xC2\x81\xC2\xA7\xC3\xA3\xE2\x80\x9A\xC2\xB3\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xC6\x92\xE2\x80\x99\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xE2\x80\x9A\xE2\x80\x99\xC3\xA9\xC2\x20\xC2\xB2\xC3\xA3\xE2\x80\x9A"),
        TransformHex:     nil,
        TransformError:   ErrInvalid,
        TransformedHex:   []byte("Caf\xC3\x83\xC2\xA9\xC3\xA3\xC2\x81\xC2\xA7\xC3\xA3\xE2\x80\x9A\xC2\xB3\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xC6\x92\xE2\x80\x99\xC3\xA3\xC6\x92\xC2\xBC\xC3\xA3\xE2\x80\x9A\xE2\x80\x99\xC3\xA9\xC2\x20\xC2\xB2\xC3\xA3\xE2\x80\x9A"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     50,
    },
    {
        Name:             "Double_Encoded_Edge_Case_1",
        TestStringHex:    []byte("MATÄšJ"),
        TransformHex:     []byte("MATĚJ"),
        TransformedHex:   []byte("MATĚJ"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     4,
    },
    {
        Name:             "Double_Encoded_Edge_Case_2",
        TestStringHex:    []byte("KONECNÄš DOBRA"),
        TransformHex:     []byte("KONECNĚ DOBRA"),
        TransformedHex:   []byte("KONECNĚ DOBRA"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     7,
    },
    {
        Name:             "Double_Encoded_Edge_Case_3",
        TestStringHex:    []byte("ÄŠhess"),
        TransformHex:     []byte("Ċhess"),
        TransformedHex:   []byte("Ċhess"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Edge_Case_4",
        TestStringHex:    []byte("ÄŽakujem"),
        TransformHex:     []byte("Ďakujem"),
        TransformedHex:   []byte("Ďakujem"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Edge_Case_5",
        TestStringHex:    []byte("DoÄŸan"),
        TransformHex:     []byte("Doğan"),
        TransformedHex:   []byte("Doğan"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     3,
    },
    {
        Name:             "Double_Encoded_Edge_Case_6",
        TestStringHex:    []byte("Knock-ÎŸut"),
        TransformHex:     []byte("Knock-Οut"),
        TransformedHex:   []byte("Knock-Οut"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     7,
    },
    {
        Name:             "Double_Encoded_Edge_Case_7",
        TestStringHex:    []byte("Åšwinia"),
        TransformHex:     []byte("Świnia"),
        TransformedHex:   []byte("Świnia"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     1,
    },
    {
        Name:             "Double_Encoded_Edge_Case_8",
        TestStringHex:    []byte("koÍžok"),
        TransformHex:     []byte("ko͞ok"),
        TransformedHex:   []byte("ko͞ok"),
        DetectResult:     MAYBE_DOUBLE_ENCODED,
        DetectOffset:     3,
    },
    {
        Name:             "UTF8_Edge_Case_1",
        TestStringHex:    []byte("Úžasná"),
        TransformHex:     decode("da9e61736ee1"),
        TransformedHex:   []byte("Úžasná"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Edge_Case_2",
        TestStringHex:    []byte("Úžasn"),
        TransformHex:     decode("da9e61736e"),
        TransformedHex:   []byte("Úžasn"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Edge_Case_3",
        TestStringHex:    []byte("MÍšhrå"),
        TransformHex:     decode("4dcd9a6872e5"),
        TransformedHex:   []byte("MÍšhrå"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Edge_Case_4",
        TestStringHex:    []byte("2×"),
        TransformHex:     decode("32d7"),
        TransformedHex:   []byte("2×"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Edge_Case_5",
        TestStringHex:    decode("c3a0c2a0"), // à\u00a0
        TransformHex:     decode("e0a0"),
        TransformedHex:   decode("c3a0c2a0"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
        DetectOffset:     1,
    },
    {
        Name:             "UTF8_Edge_Case_6",
        TestStringHex:    []byte("abcðŸ˜"),
        TransformHex:     decode("616263f09f98"),
        TransformedHex:   []byte("abcðŸ˜"),
        TransformedError: ErrNoop,
        DetectResult:     UNKNOWN,
        DetectOffset:     4,
    },
    {
        Name:             "UTF8_Edge_Case_7",
        TestStringHex:    []byte("nè…"),
        TransformHex:     decode("6ee885"),
        TransformedHex:   []byte("nè…"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
        DetectOffset:     2,
    },
    {
        Name:             "UTF8_Edge_Case_8",
        TestStringHex:    []byte("qué¡"),
        TransformHex:     decode("7175e9a1"),
        TransformedHex:   []byte("qué¡"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
        DetectOffset:     3,
    },
    {
        Name:             "UTF8_Edge_Case_9",
        TestStringHex:    []byte("JÜŠt GØ"),
        TransformHex:     decode("4adc8a742047d8"),
        TransformedHex:   []byte("JÜŠt GØ"),
        TransformedError: ErrNoop,
        DetectResult:     MAYBE_OTHER,
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
            result, _, _, offset := d.Detect(tc.TestStringHex)
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
            assert.ErrorIs(t, err, tc.TransformError)
            assert.Equal(t, tc.TransformHex, r)
        })
    }
}

func TestTransformed(t *testing.T) {
    d := NewDecoder()

    for _, tc := range testCases {
        t.Run(tc.Name, func(t *testing.T) {
            r, err := d.Transform(tc.TestStringHex)
            assert.ErrorIs(t, err, tc.TransformedError)
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
