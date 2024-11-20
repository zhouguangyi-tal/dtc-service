package request

type Result struct {
	Strategy Strategy           `json:"strategy"`
	Status   int                `json:"status"`
	Time     string             `json:"time"`
	Config   map[string]Version `json:"config"`
}

type Strategy struct {
	Force bool `json:"force"`
}

type Version struct {
	Package Package `json:"package"`
	Name    string  `json:"name"`
	Title   string  `json:"title"`
	Version string  `json:"version"`
	Content Content `json:"content"`
}

type Package struct {
	FullPackage FullPackage `json:"full_package"`
}

type FullPackage struct {
	DownloadURL string `json:"download_url"`
	MD5         string `json:"md5"`
}

type Content struct {
	En string `json:"en"`
	Cn string `json:"cn"`
	Zh string `json:"zh"`
}
