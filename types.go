package main

import "time"

type Snippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	Files       []File `json:"files"`
}
type File struct {
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
}

type SnippetCreateResponse struct {
	ID            int                          `json:"id"`
	Title         string                       `json:"title"`
	Description   string                       `json:"description"`
	Visibility    string                       `json:"visibility"`
	Author        SnippetCreateResponseAuthor  `json:"author"`
	ExpiresAt     interface{}                  `json:"expires_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	CreatedAt     time.Time                    `json:"created_at"`
	ProjectID     interface{}                  `json:"project_id"`
	WebURL        string                       `json:"web_url"`
	RawURL        string                       `json:"raw_url"`
	SSHURLToRepo  string                       `json:"ssh_url_to_repo"`
	HTTPURLToRepo string                       `json:"http_url_to_repo"`
	FileName      string                       `json:"file_name"`
	Files         []SnippetCreateResponseFiles `json:"files"`
}

type SnippetCreateResponseAuthor struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type SnippetCreateResponseFiles struct {
	Path   string `json:"path"`
	RawURL string `json:"raw_url"`
}
