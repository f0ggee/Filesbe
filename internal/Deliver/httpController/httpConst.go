package httpController

import "errors"

var (
	Success            = "Success"
	DomainName         = "https://filesbes.com/"
	EncryptURLDownload = "https://filesbes.com/d2/"
	UrlDownload        = "https://filesbes.com/d/"
	MainPageUrl        = "/main"
	LoginPage          = "/login"
	Bots               = "Bot"
)

var (
	FileUrlName = "name"
	TypeFile    = "bool"
)

type AnswerRegister struct {
	StatusOfOperation string `json:"status_of_operation"`
	UrlToRedirect     string `json:"url_to_redirect"`
	Error             string `json:"error"`
}
type AnswerUserCheck struct {
	UrlToRedirect string `json:"url_to_redirect"`
	Error         string `json:"error"`
}
type AnswerLogin struct {
	StatusOfOperation string `json:"status_of_operation"`
	UrlToRedirect     string `json:"url_to_redirect"`
	ErrorMessage      string `json:"error_message"`
}
type AnswerUrlBuilder struct {
	StatusOperation string `json:"StatusOperation"`
	Url             string `json:"Url"`
	ErrorMessage    string `json:"ErrorMessage"`
}
type AnswerUploaderFileNoEncrypt struct {
	StatusOperation string `json:"StatusOperation"`
	UrlToRedirect   string `json:"UrlRedict"`
	Error           string `json:"Error"`
}
type AnswerDownloadNoEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
	Url             string `json:"url"`
}
type AnswerFileDownloadEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
}

type AnswerUploadEncrypt struct {
	StatusOperation string `json:"status_operation"`
	Error           string `json:"error"`
	UrlToRedirect   string `json:"url_to_redirect"`
}
type OutComingUrlData struct {
	NameFile string
	FileType string
}

// /the error list
var (
	ErrorFileNameEmpty   = errors.New("the file's filed is empty")
	ErrorCantGetFileName = errors.New("the file's name isn't set")
	ErrorFile            = errors.New("can't get a file")
	Break                = "break"
	NotStart             = "not_start"
	ErrorGetCookie       = errors.New("error get an user's cookie")
	ErrorAuthExpired     = errors.New("user's auth expired")
	ErrorSaveCookie      = errors.New("error save a cookie")
)
