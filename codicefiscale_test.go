package codicefiscale

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	fixtures := []struct {
		CF 			string
		Valid		bool
		Surname		string
		Name		string
		Sex			Sex
		BirthYear	int
		BirthMonth	int
		BirthDay	int
		BirthPlace	string
	}{
		{ "", false, "", "", UnknownSex, 0, 0, 0, "" },
		{ "AAA", false, "", "", UnknownSex, 0, 0, 0, "" },
		{ "RSSMRA77L18H501W", true, "RSS", "MRA", Male, 1977, 7, 18, "Roma" },
		{ "RSSMRA77L18H501Z", false, "RSS", "MRA", Male, 1977, 7, 18, "Roma" },
		{ "RSSMRA77L18H501WA", false, "RSS", "MRA", Male, 1977, 7, 18, "Roma" },
		{ "BNCGRC69A41G048P", true, "BNC", "GRC", Female, 1969, 1, 1, "Olmo Gentile" },
		{ "BNCGRC69A41G048A", false, "BNC", "GRC", Female, 1969, 1, 1, "Olmo Gentile" },
		{ "RSSMRA77L18Z103I", true, "RSS", "MRA", Male, 1977, 7, 18, "Belgio" },
		{ "RSSLVR64M44Z602P", true, "RSS", "LVR", Female, 1964, 8, 4, "Brasile" },
	}

	for _, fixture := range fixtures {
		cf, err := Decode(fixture.CF)
		if !fixture.Valid {
			if err == nil {
				t.Errorf("no error with invalid code: '%v'", fixture.CF)
				continue
			}
			continue
		}
		// fmt.Printf("%#v\n", cf)

		if err != nil {
			t.Errorf("error decoding code '%v': %w", fixture.CF, err)
		}
		if cf.Surname != fixture.Surname {
			t.Errorf("Surname mismatch, expected:'%v' got:'%v'", fixture.Surname, cf.Surname)
		}
		if cf.Name != fixture.Name {
			t.Errorf("Name mismatch, expected:'%v' got:'%v'", fixture.Name, cf.Name)
		}
		if cf.Sex != fixture.Sex {
			t.Errorf("Sex mismatch, expected:'%v' got:'%v'", fixture.Sex, cf.Sex)
		}
		if cf.BirthDateYear != fixture.BirthYear {
			t.Errorf("Birth year mismatch, expected:'%v' got:'%v'", fixture.BirthYear, cf.BirthDateYear)
		}
		if cf.BirthDateMonth != fixture.BirthMonth {
			t.Errorf("Birth month mismatch, expected:'%v' got:'%v'", fixture.BirthMonth, cf.BirthDateMonth)
		}
		if cf.BirthDateDay != fixture.BirthDay {
			t.Errorf("Birth day mismatch, expected:'%v' got:'%v'", fixture.BirthDay, cf.BirthDateDay)
		}
		if strings.ToUpper(cf.BirthPlaceName) != strings.ToUpper(fixture.BirthPlace) {
			t.Errorf("Birth place mismatch, expected:'%v' got:'%v'", fixture.BirthPlace, cf.BirthPlaceName)
		}
		if cf.BirthPlaceNazione.DenominazioneIT == "Italia" && cf.BirthPlaceComune.CodiceCatastale == "" {
			t.Errorf("Missing 'codice catastale' from results")
		}
	}
}
