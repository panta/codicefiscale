package codicefiscale

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/panta/codicefiscale/comuni"
)

type Sex int

const (
	UnknownSex Sex = iota
	Male
	Female
)

var monthsRunes = []rune{'A', 'B', 'C', 'D', 'E', 'H', 'L', 'M', 'P', 'R', 'S', 'T'}

var cinOdd = map[rune]int {
	'0': 1, '1': 0, '2': 5, '3': 7, '4': 9, '5': 13,
	'6': 15, '7': 17, '8': 19, '9': 21, 'A': 1, 'B': 0,
	'C': 5, 'D': 7, 'E': 9, 'F': 13, 'G': 15, 'H': 17,
	'I': 19, 'J': 21, 'K': 2, 'L': 4, 'M': 18, 'N': 20,
	'O': 11, 'P': 3, 'Q': 6, 'R': 8, 'S': 12, 'T': 14,
	'U': 16, 'V': 10, 'W': 22, 'X': 25, 'Y': 24, 'Z': 23,
}

var cinEven = map[rune]int {
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5,
	'6': 6, '7': 7, '8': 8, '9': 9, 'A': 0, 'B': 1,
	'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7,
	'I': 8, 'J': 9, 'K': 10, 'L': 11, 'M': 12, 'N': 13,
	'O': 14, 'P': 15, 'Q': 16, 'R': 17, 'S': 18, 'T': 19,
	'U': 20, 'V': 21, 'W': 22, 'X': 23, 'Y': 24, 'Z': 25,
}

// normalizeText removes diacritics and accents from a string.
// Based on:
//   - https://medium.com/@swdream/golang-remove-all-accents-in-string-319abf6a7f5b
//   - https://stackoverflow.com/questions/26722450/remove-diacritics-using-go
//   - https://rosettacode.org/wiki/Strip_control_codes_and_extended_characters_from_a_string#Go
func normalizeText(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ := transform.String(t, text)
	return s
}

type CodiceFiscale struct {
	Code				string
	Surname				string
	Name				string
	Sex					Sex
	BirthDate			time.Time
	BirthDateYear		int
	BirthDateMonth		int
	BirthDateDay		int
	BirthPlaceName		string
	BirthPlaceComune	comuni.Comune
	BirthPlaceNazione	comuni.Nazione
	Raw					CodiceFiscaleRaw
}

type CodiceFiscaleRaw struct {
	Code           string
	Surname        string
	Name           string
	BirthDate      string
	BirthDateYear  string
	BirthDateMonth string
	BirthDateDay   string
	BirthPlace     string
	CIN            string
}

// decodeRaw decodes the "codice fiscale" into its raw constituents.
func decodeRaw(codice string) (*CodiceFiscaleRaw, error) {
	cf_re := regexp.MustCompile(`(?i)^` +
		`([a-z]{3})` +
		`([a-z]{3})` +
		`(([a-z\d]{2})([abcdehlmprst]{1})([a-z\d]{2}))` +
		`([a-z]{1}[a-z\d]{3})` +
		`([a-z]{1})$`)

	normalized := normalizeText(codice)
	normalized = strings.ToUpper(normalized)

	matches := cf_re.FindStringSubmatch(normalized)
	if len(matches) <= 0 {
		return nil, fmt.Errorf("invalid 'codice fiscale': %#v", normalized)
	}

	return &CodiceFiscaleRaw{
		Code:           normalized,
		Surname:        matches[1],
		Name:           matches[2],
		BirthDate:      matches[3],
		BirthDateYear:  matches[4],
		BirthDateMonth: matches[5],
		BirthDateDay:   matches[6],
		BirthPlace:     matches[7],
		CIN:            matches[8],
	}, nil
}

// Decode decodes the "codice fiscale".
func Decode(codice string) (*CodiceFiscale, error) {
	raw, err := decodeRaw(codice)
	if err != nil {
		return nil, err
	}

	birthdateYear, err := strconv.Atoi(omocodiaDecodeTrans(raw.BirthDateYear))
	if err != nil {
		return nil, fmt.Errorf("can't convert 'year' component: %w", err)
	}
	birthdateMonth := bytes.IndexByte([]byte(string(monthsRunes)), byte(rune(raw.BirthDateMonth[0]))) + 1
	birthdateDay, err := strconv.Atoi(omocodiaDecodeTrans(raw.BirthDateDay))
	if err != nil {
		return nil, fmt.Errorf("can't convert 'day' component: %w", err)
	}

	var sex Sex
	if birthdateDay > 40 {
		birthdateDay -= 40
		sex = Female
	} else {
		sex = Male
	}

	now := time.Now()
	currentYear := now.Year()

	currentCentury := currentYear / 100
	birthdateYear += (currentCentury * 100)
	if birthdateYear > currentYear {
		birthdateYear -= 100
	}

	birthdate := time.Date(birthdateYear, time.Month(birthdateMonth), birthdateDay,
		0, 0, 0, 0, time.Local)

	var birthplaceNazione comuni.Nazione
	var birthplaceComune comuni.Comune
	birthplaceName := ""
	catastale := string(raw.BirthPlace[0]) + omocodiaDecodeTrans(raw.BirthPlace[1:])
	comuneIndex, ok := comuni.Catastale2Comune[catastale]
	if !ok {
		nazioneIndex, ok := comuni.Codice2Nazione[catastale]
		if !ok {
			return nil, fmt.Errorf("birth place code not found")
		} else {
			birthplaceNazione = comuni.Nazioni[nazioneIndex]
			birthplaceName = birthplaceNazione.DenominazioneIT
		}
	} else {
		birthplaceComune = comuni.Comuni[comuneIndex]
		birthplaceNazione = comuni.Nazioni[comuni.ItalyIndex]
		birthplaceName = birthplaceComune.Nome
	}

	cinComputed, err := computeCIN(codice)
	if err != nil {
		return nil, err
	}
	if raw.CIN != cinComputed {
		return nil, fmt.Errorf("wrong CIN (computed:'%v' found:'%v')", cinComputed, raw.CIN)
	}

	return &CodiceFiscale{
		Code:           	raw.Code,
		Surname:        	raw.Surname,
		Name:           	raw.Name,
		Sex:				sex,
		BirthDate:			birthdate,
		BirthDateYear:		birthdateYear,
		BirthDateMonth:		birthdateMonth,
		BirthDateDay:		birthdateDay,
		BirthPlaceName:		birthplaceName,
		BirthPlaceNazione:	birthplaceNazione,
		BirthPlaceComune:	birthplaceComune,
		Raw:            	*raw,
	}, err
}

// computeCIN computes the CIN character from 'codice fiscale'.
func computeCIN(codice string) (string, error) {
	if (len(codice) < 15) || (len(codice) > 16) {
		return "", fmt.Errorf("the code length must be 15 or 16")
	}

	cinValue := 0
	for index, char := range codice[:15] {
		if (index + 1) % 2 != 0 {
			cinValue += cinOdd[char]
		} else {
			cinValue += cinEven[char]
		}
	}
	cinCode := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"[cinValue % 26]
	return string(cinCode), nil
}

var omocodiaDigitToLetter = map[rune]rune{
	'0': 'L',
	'1': 'M',
	'2': 'N',
	'3': 'P',
	'4': 'Q',
	'5': 'R',
	'6': 'S',
	'7': 'T',
	'8': 'U',
	'9': 'V',
}

var omocodiaLetterToDigit = map[rune]rune{
	'L': '0',
	'M': '1',
	'N': '2',
	'P': '3',
	'Q': '4',
	'R': '5',
	'S': '6',
	'T': '7',
	'U': '8',
	'V': '9',
}

// omocodiaDecodeTrans maps from letters to omocode digits.
func omocodiaDecodeTrans(str string) string {
	omocodiaDecoder := func (src rune) rune {
		dst, ok := omocodiaLetterToDigit[src]
		if ok {
			return dst
		}
		return src
	}
	return strings.Map(omocodiaDecoder, str)
}

// omocodiaEncodeTrans maps from digits to omocode letters.
func omocodiaEncodeTrans(str string) string {
	omocodiaEncoder := func (src rune) rune {
		dst, ok := omocodiaDigitToLetter[src];
		if ok {
			return dst
		}
		return src
	}
	return strings.Map(omocodiaEncoder, str)
}
