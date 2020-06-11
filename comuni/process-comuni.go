// The following directive is necessary to make the package coherent:

// +build ignore

// Based on:
//  - https://blog.carlmjohnson.net/post/2016-11-27-how-to-use-go-generate/

// This program generates comuni/comuni-generated-data.go. It can be invoked by running
// go generate
package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/panta/codicefiscale/comuni"
)

const (
	COMUNI_URL = "https://github.com/matteocontrini/comuni-json/blob/master/comuni.json?raw=true"
	NAZIONI_URL = "https://www.istat.it/it/files//2011/01/Elenco-codici-e-denominazioni-unita-territoriali-estere.zip"
)

var (
	outputFilename string = path.Join("comuni", "comuni-generated-data.go")
)


func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readComuniJSONFromUrl(url string) ([]comuni.Comune, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var comuniList []comuni.Comune
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &comuniList); err != nil {
		return nil, err
	}

	return comuniList, nil
}

func readNazioniCSVFromUrl(url string) ([]comuni.Nazione, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// zipReader, err := zip.OpenReader(buf)
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	var records []comuni.Nazione

	processZipCsv := func (zipFile *zip.File) error {
		ext := path.Ext(zipFile.Name)
		if strings.ToLower(ext) != ".csv" {
			return nil
		}
		rc, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		csvReader := csv.NewReader(charmap.ISO8859_1.NewDecoder().Reader(rc))
		csvReader.Comma = ';'
		// headerRec, err := csvReader.Read()
		// if err != nil {
		// 	log.Fatal(err)
		// 	return err
		// }
		// for _, v := range headerRec {
		// 	fmt.Println(v)
		// }
		for {
			fields, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading CSV '%v': %w", zipFile.Name, err)
			}

			if len(fields) < 15 {
				log.Printf("too few columns (%v) in fields, skipping...\n", len(fields))
				continue
			}
			record := comuni.Nazione{
				Stato_Territorio:              fields[0],
				CodiceContinente:              fields[1],
				DenominazioneContinente:       fields[2],
				CodiceArea:                    fields[3],
				DenominazioneArea:             fields[4],
				CodiceISTAT:                   fields[5],
				DenominazioneIT:               fields[6],
				DenominazioneEN:               fields[7],
				CodiceMIN:                     fields[8],
				CodiceAT:                      fields[9],
				CodiceUNSD_M49:                fields[10],
				Codice_ISO_3166_alpha2:        fields[11],
				Codice_ISO_3166_alpha3:        fields[12],
				Codice_ISTAT_Stato_Padre:      fields[13],
				Codice_ISO_alpha3_Stato_Padre: fields[14],
			}
			records = append(records, record)
		}

		return nil
	}

	for _, zipFile := range zipReader.File {
		err := processZipCsv(zipFile)
		if err != nil {
			return nil, err
		}
		if len(records) > 0 {
			break
		}
	}

	return records, nil
}

func main() {
	nazioniList, err := readNazioniCSVFromUrl(NAZIONI_URL)
	die(err)

	comuniList, err := readComuniJSONFromUrl(COMUNI_URL)
	die(err)

	italyIndex := -1
	codice2nazione := map[string]int{}
	for nazioneIndex, nazione := range nazioniList {
		if (strings.ToLower(nazione.DenominazioneIT) == "italia") || (strings.ToLower(nazione.DenominazioneEN) == "italy") {
			italyIndex = nazioneIndex
		}
		if nazione.CodiceAT == "" || nazione.CodiceAT == "n.d." {
			continue
		}
		codice2nazione[nazione.CodiceAT] = nazioneIndex
	}

	catastale2comune := map[string]int{}
	for comuneIndex, comune := range comuniList {
		catastale2comune[comune.CodiceCatastale] = comuneIndex
	}

	now := time.Now()

	outputFile, err := os.Create(outputFilename)
	die(err)
	defer outputFile.Close()

	packageTemplate := template.Must(template.New("comuni-go").Parse(packageTemplateText))
	context := map[string]interface{}{
		"Nazioni": nazioniList,
		"Codice2Nazione": codice2nazione,
		"ItalyIndex": italyIndex,
		"Comuni": comuniList,
		"Catastale2Comune": catastale2comune,
		"URL": COMUNI_URL,
		"Timestamp": now,
		"TimestampFormatted": now.Format(time.RFC3339),
	}
	packageTemplate.Execute(outputFile, context)
}

var packageTemplateText = `// Code generated by go generate; DO NOT EDIT.
//
// This file was automagically generated at
//   {{.Timestamp}}
// using data from
//   {{ .URL }}

package comuni

var Nazioni = []Nazione{
{{range .Nazioni}}
	Nazione{
		Stato_Territorio:              "{{.Stato_Territorio}}",
		CodiceContinente:              "{{.CodiceContinente}}",
		DenominazioneContinente:       "{{.DenominazioneContinente}}",
		CodiceArea:                    "{{.CodiceArea}}",
		DenominazioneArea:             "{{.DenominazioneArea}}",
		CodiceISTAT:                   "{{.CodiceISTAT}}",
		DenominazioneIT:               "{{.DenominazioneIT}}",
		DenominazioneEN:               "{{.DenominazioneEN}}",
		CodiceMIN:                     "{{.CodiceMIN}}",
		CodiceAT:                      "{{.CodiceAT}}",
		CodiceUNSD_M49:                "{{.CodiceUNSD_M49}}",
		Codice_ISO_3166_alpha2:        "{{.Codice_ISO_3166_alpha2}}",
		Codice_ISO_3166_alpha3:        "{{.Codice_ISO_3166_alpha3}}",
		Codice_ISTAT_Stato_Padre:      "{{.Codice_ISTAT_Stato_Padre}}",
		Codice_ISO_alpha3_Stato_Padre: "{{.Codice_ISO_alpha3_Stato_Padre}}",
	},
{{end}}
}

var ItalyIndex int = {{.ItalyIndex}}

var Codice2Nazione = map[string]int{
{{ range $key, $value := .Codice2Nazione }}
   "{{ $key }}": {{ $value }},
{{ end }}
}

var Comuni = []Comune{
{{range .Comuni}}
	Comune{
		Cap:             []string{ {{range .Cap}}"{{.}}",{{end}} },
		Codice:          "{{.Codice}}",
		CodiceCatastale: "{{.CodiceCatastale}}",
		Nome:            "{{.Nome}}",
		Popolazione:     {{.Popolazione}},
		Provincia: Provincia{
			Codice: "{{.Provincia.Codice}}",
			Nome:   "{{.Provincia.Nome}}",
		},
		Regione: Regione{
			Codice: "{{.Regione.Codice}}",
			Nome:   "{{.Regione.Nome}}",
		},
		Sigla: "{{.Sigla}}",
		Zona: Zona{
			Codice: "{{.Zona.Codice}}",
			Nome:   "{{.Zona.Nome}}",
		},
	},
{{end}}
}

var Catastale2Comune = map[string]int{
{{ range $key, $value := .Catastale2Comune }}
   "{{ $key }}": {{ $value }},
{{ end }}
}

// Code generated - DO NOT EDIT
`
