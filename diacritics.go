package dblenc

type Language uint32

// Langauges that use CP1252
const (
    L_NONE Language = 0
    L_FR   Language = 1 << (iota - 1) // french
    L_PT                              // portuguese
    L_ES                              // spanish
    L_IT                              // italian
    L_DE                              // german
    L_DA                              // danish
    L_NO                              // norwegian
    L_FI                              // finnish
    L_IS                              // icelandic
    L_FO                              // faroese
    L_NL                              // dutch
    L_CY                              // welsh
    L_HU                              // hungarian
    L_CZ                              // czech
    L_SK                              // slovak
    L_RO                              // romanian
    L_ET                              // estonian
    L_SV                              // swedish
    L_GA                              // irish
    L_SQ                              // albanian
    L_TR                              // turkish
    L_AZ                              // azerbaijani
    L_MT                              // maltese
    L_PL                              // polish
    L_GR                              // greek
    L_FK                              // some other
)

// Letters with diacritics
var Diacritics = [0x0400]Language{
    // 0x00A1: L_ES,                                                                // ¡
    // 0x00BF: L_ES,                                                                // ¿
    0x00C0: L_FR | L_IT | L_PT | L_CY,                                           // À
    0x00C1: L_IS | L_FO | L_GA | L_CZ | L_SK | L_HU | L_ES | L_PT,               // Á
    0x00C2: L_FR | L_RO | L_PT | L_CY | L_TR,                                    // Â
    0x00C3: L_PT,                                                                // Ã
    0x00C4: L_DE | L_FI | L_SV | L_ET | L_SK,                                    // Ä
    0x00C5: L_SV | L_DA | L_NO | L_FI,                                           // Å
    0x00C6: L_IS | L_FO | L_DA | L_NO,                                           // Æ
    0x00C7: L_FR | L_PT | L_TR | L_AZ | L_SQ,                                    // Ç
    0x00C8: L_FR | L_IT | L_PT,                                                  // È
    0x00C9: L_FR | L_PT | L_ES | L_IS | L_HU | L_CZ | L_SK | L_DA | L_NO | L_SV, // É
    0x00CA: L_FR | L_PT | L_CY,                                                  // Ê
    0x00CB: L_SQ | L_FR | L_NL,                                                  // Ë
    0x00CC: L_IT,                                                                // Ì
    0x00CD: L_IS | L_FO | L_CZ | L_SK | L_HU | L_GA | L_PT | L_ES,               // Í
    0x00CE: L_FR | L_RO,                                                         // Î
    0x00CF: L_FR | L_NL,                                                         // Ï
    0x00D0: L_IS | L_FO,                                                         // Ð
    0x00D1: L_ES,                                                                // Ñ
    0x00D2: L_IT | L_PT,                                                         // Ò
    0x00D3: L_IS | L_FO | L_GA | L_CZ | L_SK | L_HU | L_ES | L_PT,               // Ó
    0x00D4: L_FR | L_PT | L_CY,                                                  // Ô
    0x00D5: L_PT,                                                                // Õ
    0x00D6: L_DE | L_SV | L_FI | L_ET | L_HU | L_TR | L_AZ,                      // Ö
    0x00D8: L_DA | L_NO | L_FO,                                                  // Ø
    0x00D9: L_FR | L_IT | L_PT,                                                  // Ù
    0x00DA: L_IS | L_FO | L_CZ | L_SK | L_HU | L_ES | L_PT,                      // Ú
    0x00DB: L_FR | L_CY | L_PT,                                                  // Û
    0x00DC: L_DE | L_HU | L_TR | L_AZ | L_ET,                                    // Ü
    0x00DD: L_IS | L_FO,                                                         // Ý
    0x00DE: L_IS,                                                                // Þ
    0x00DF: L_DE,                                                                // ß
    0x00E0: L_FR | L_IT | L_PT | L_CY,                                           // à
    0x00E1: L_IS | L_FO | L_GA | L_CZ | L_SK | L_HU | L_ES | L_PT,               // á
    0x00E2: L_FR | L_RO | L_PT | L_CY | L_TR,                                    // â
    0x00E3: L_PT | L_ET,                                                         // ã
    0x00E4: L_DE | L_FI | L_SV | L_ET | L_SK,                                    // ä
    0x00E5: L_SV | L_DA | L_NO | L_FI,                                           // å
    0x00E6: L_IS | L_FO | L_DA | L_NO,                                           // æ
    0x00E7: L_FR | L_PT | L_TR | L_AZ | L_SQ,                                    // ç
    0x00E8: L_FR | L_IT | L_PT,                                                  // è
    0x00E9: L_FR | L_PT | L_ES | L_IS | L_HU | L_CZ | L_SK | L_DA | L_NO | L_SV, // é
    0x00EA: L_FR | L_PT | L_CY,                                                  // ê
    0x00EB: L_SQ | L_FR | L_NL,                                                  // ë
    0x00EC: L_IT,                                                                // ì
    0x00ED: L_IS | L_FO | L_CZ | L_SK | L_HU | L_GA | L_PT | L_ES,               // í
    0x00EE: L_FR | L_RO,                                                         // î
    0x00EF: L_FR | L_NL,                                                         // ï
    0x00F0: L_IS | L_FO,                                                         // ð
    0x00F1: L_ES,                                                                // ñ
    0x00F2: L_IT | L_PT,                                                         // ò
    0x00F3: L_IS | L_FO | L_GA | L_CZ | L_SK | L_HU | L_ES | L_PT,               // ó
    0x00F4: L_FR | L_PT | L_CY,                                                  // ô
    0x00F5: L_ET | L_PT,                                                         // õ
    0x00F6: L_DE | L_SV | L_FI | L_ET | L_HU | L_TR | L_AZ,                      // ö
    0x00F8: L_DA | L_NO | L_FO,                                                  // ø
    0x00F9: L_FR | L_IT | L_PT,                                                  // ù
    0x00FA: L_IS | L_FO | L_CZ | L_SK | L_HU | L_ES | L_PT,                      // ú
    0x00FB: L_FR | L_CY | L_PT,                                                  // û
    0x00FC: L_DE | L_HU | L_TR | L_AZ | L_ET,                                    // ü
    0x00FD: L_IS | L_FO,                                                         // ý
    0x00FE: L_IS,                                                                // þ
    0x00FF: L_FR | L_NL,                                                         // ÿ

    0x010A: L_MT,                                                                // Ċ - exception
    0x010C: L_CZ | L_SK,                                                         // Č - exception
    0x010E: L_CZ | L_SK,                                                         // Ď - exception
    0x011A: L_CZ | L_SK,                                                         // Ě - exception
    0x011E: L_TR,                                                                // Ğ - exception
    0x011F: L_TR,                                                                // ğ - exception
    0x0121: L_MT,                                                                // ġ - exception
    0x014C: L_FK,                                                                // Ō - exception
    0x015E: L_AZ | L_TR,                                                         // Ş - exception
    0x015F: L_AZ | L_TR,                                                         // ş - exception
    0x015A: L_PL,                                                                // Ś - exception

    0x0160: L_CZ | L_SK | L_ET,                                                  // Š
    0x0161: L_CZ | L_SK | L_ET,                                                  // š
    0x0178: L_FR | L_NL,                                                         // Ÿ
    0x017D: L_CZ | L_SK | L_ET,                                                  // Ž
    0x017E: L_CZ | L_SK | L_ET,                                                  // ž

    0x019F: L_GR,                                                                // Ɵ - exception
    0x039F: L_GR,                                                                // Ο - exception
}