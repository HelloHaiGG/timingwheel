package ilogger

type LogMsg struct {
	Who   string `json:"who"`
	When  string `json:"when"`
	Where string `json:"where"`
	What  string `json:"what"`
	Level string `json:"level"`
}

