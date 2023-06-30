package domain

type Project struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}

type ProjectID struct {
	ID string `json:"id"`
}

type IssueCreateRequest struct {
	ProjectID   ProjectID `json:"project"`
	Key         string    `json:"idReadable"`
	ID          string    `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
}

// [ ] try `,omitempty` to remove extra struct;

type IssueUpdateRequest struct {
	IssueCreateRequest
	CustomFields []CustomField `json:"customFields"`
}

type CustomField struct {
	Name  string `json:"name"`
	Type  string `json:"$type"`
	Value string `json:"value"`
}
