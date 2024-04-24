package models

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type DocumentListResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	MaxFiles  string `json:"maxFiles"`
	Documents []struct {
		Name            string `json:"name"`
		MimeType        string `json:"mimeType"`
		Size            string `json:"size"`
		FileType        string `json:"fileType"`
		NumberOfRecords int    `json:"numberOfRecords"`
		Priority        string `json:"priority"`
	} `json:"documents"`
}

type DocumentListRequest struct {
	DepartmentId string `json:"departmentId"`
	StoreId      string `json:"storeId"`
	Source       string `json:"source"`
}
