package prompts

type Info struct {
	Model      string  `json:"model"`
	Created_at string  `json:"created_at"`
	Message    message `json:"message"`
	Done       bool    `json:"done"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
