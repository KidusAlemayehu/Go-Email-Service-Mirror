package models

type EmailTask struct {
	To          string       `json:"to"`
	From        string       `json:"from"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	Cc          []string     `json:"cc,omitempty"`
	Bcc         []string     `json:"bcc,omitempty"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
}
