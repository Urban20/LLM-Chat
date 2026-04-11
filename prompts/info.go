package prompts

type Info struct {
	Model                string `json:"model"`
	Created_at           string `json:"created_at"`
	Response             string `json:"response"`
	Done                 bool   `json:"done"`
	Context              []int  `json:"context"`
	Total_duration       int    `json:"total_duration"`
	Prompt_eval_duration int    `json:"prompt_eval_duration"`
	Eval_count           int    `json:"eval_count"`
	Eval_duration        int    `json:"eval_duration"`
}
