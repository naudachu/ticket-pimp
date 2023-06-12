package domain

type Git struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Url      string `json:"url"`
	CloneUrl string `json:"clone_url"`
	HtmlUrl  string `json:"Html_url"`
	SshUrl   string `json:"ssh_url"`
}
