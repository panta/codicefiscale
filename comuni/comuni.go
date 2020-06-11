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

type Nazione struct {
	Stato_Territorio				string
	CodiceContinente				string
	DenominazioneContinente			string
	CodiceArea						string
	DenominazioneArea				string
	CodiceISTAT						string
	DenominazioneIT					string
	DenominazioneEN					string
	CodiceMIN						string
	CodiceAT						string
	CodiceUNSD_M49					string
	Codice_ISO_3166_alpha2			string
	Codice_ISO_3166_alpha3			string
	Codice_ISTAT_Stato_Padre		string
	Codice_ISO_alpha3_Stato_Padre	string
}
