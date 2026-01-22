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
        TransformHex:     decode("48e96c6c6f20e7e0207661207472e873206269656e"),
        TransformedHex:   []byte("Héllo çà va très bien"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     4,
    },
    {
        Name:             "UTF8_German",
        TestString:       "Guten Tag, wie geht's? Größe, Ärger, Übung",
        TestStringHex:    []byte("Guten Tag, wie geht's? Größe, Ärger, Übung"),
        TransformHex:     decode("477574656e205461672c20776965206765687427733f204772f6df652c20c4726765722c20dc62756e67"),
        TransformedHex:   []byte("Guten Tag, wie geht's? Größe, Ärger, Übung"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     27,
    },
    {
        Name:             "UTF8_Italian",
        TestString:       "Città, perché, caffè",
        TestStringHex:    []byte("Città, perché, caffè"),
        TransformHex:     decode("43697474e02c207065726368e92c2063616666e8"),
        TransformedHex:   []byte("Città, perché, caffè"),
        TransformedError: ErrNoop,
        DetectResult:     OTHER_CHARSET,
        DetectOffset:     7,
    },
    {
        Name:             "UTF8_Portuguese",
        TestString:       "São Paulo, ação, coração",
        TestStringHex:    []byte("São Paulo, ação, coração"),
        TransformHex:     decode("53e36f205061756c6f2c2061e7e36f2c20636f7261e7e36f"),
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
        TransformedHex:   []byte("abc"),
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
