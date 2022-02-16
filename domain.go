package lingvanex

type Language struct {
	FullCode    string `json:"full_code"`
	EnglishName string `json:"englishName"`
}

type LanguagesResponse struct {
	Err    string     `json:"err,omitempty"`
	Result []Language `json:"result,omitempty"`
}

type TranslateRequest struct {
	Transliteration bool   `json:"enableTransliteration,omitempty"`
	From            string `json:"from,omitempty"`
	To              string `json:"to"`
	Data            string `json:"data"`
	TranslateMode   string `json:"translateMode,omitempty"`
	Platform        string `json:"platform,omitempty"`
}

type TranslateResponse struct {
	Err                   string `json:"err,omitempty"`
	Result                string `json:"from,omitempty"`
	SourceTransliteration string `json:"sourceTransliteration,omitempty"`
	TargetTransliteration string `json:"targetTransliteration,omitempty"`
}
