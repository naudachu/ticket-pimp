package domain

type Git struct {
	Name     string `json:"name"`      // "poop"
	FullName string `json:"full_name"` // "developer/poop"
	Private  bool   `json:"private"`
	Url      string `json:"url"`       // "http://localhost:8081/api/v3/repos/developer/poop"
	CloneUrl string `json:"clone_url"` // "http://localhost:8081/git/developer/poop.git"
	HtmlUrl  string `json:"Html_url"`  // "http://localhost:8081/developer/poop"
	SshUrl   string `json:"ssh_url"`   // ?!
}
