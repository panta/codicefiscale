package comuni

type Provincia struct {
	Codice string `json:"codice"`
	Nome   string `json:"nome"`
}

type Regione struct {
	Codice string `json:"codice"`
	Nome   string `json:"nome"`
}

type Zona struct {
	Codice string `json:"codice"`
	Nome   string `json:"nome"`
}

type Comune struct {
	Cap             []string `json:"cap"`
	Codice          string   `json:"codice"`
	CodiceCatastale string   `json:"codiceCatastale"`
	Nome            string   `json:"nome"`
	Popolazione     int64    `json:"popolazione"`
	Provincia       Provincia `json:"provincia"`
	Regione 		Regione	  `json:"regione"`
	Sigla			string		`json:"sigla"`
	Zona  			Zona	  `json:"zona"`
}
