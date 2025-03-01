package model

type Credential struct {
	Uci               string
	Password          string
	ApplicationNumber string
}

type ProfileResponse struct {
	Uci         string `json:"uci"`
	Apps        []App  `json:"apps"`
	LastUpdated int64  `json:"lastUpdated"`
	Type        string `json:"type"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}

type App struct {
	Status       string `json:"status"`
	AppNumber    string `json:"appNumber"`
	Tasks        uint32 `json:"tasks"`
	LastActivity int64  `json:"lastActivity"`
	Uci          string `json:"uci"`
}

type AuthResponse struct {
	AuthenticationResult struct {
		AccessToken  string `json:"AccessToken"`
		ExpiresIn    uint32 `json:"ExpiresIn"`
		IdToken      string `json:"IdToken"`
		RefreshToken string `json:"RefreshToken"`
		TokenType    string `json:"TokenType"`
	} `json:"AuthenticationResult"`
	ChallengeParameters struct{} `json:"ChallengeParameters"`
}

type StatusResponse struct {
	ApplicationNumber string        `json:"applicationNumber"`
	Uci               string        `json:"uci"`
	LastUpdatedTime   int64         `json:"lastUpdatedTime"`
	Status            string        `json:"status"`
	Actions           []Action      `json:"actions"`
	Activities        []Activity    `json:"activities"`
	History           []HistoryItem `json:"history"`
}

type Action struct {
	IsNew    bool     `json:"isNew"`
	Activity string   `json:"activity"`
	LoadTime int64    `json:"loadTime"`
	Title    Language `json:"title"`
	Content  Language `json:"content"`
}

type Activity struct {
	Activity string `json:"activity"`
	Order    uint8  `json:"order"`
	Status   string `json:"status"`
}

type HistoryItem struct {
	Time     int64    `json:"time"`
	Type     string   `json:"type"`
	Activity string   `json:"activity"`
	Title    Language `json:"title"`
	Text     Language `json:"text"`
}

type Language struct {
	En string `json:"en"`
	Fr string `json:"fr"`
}
