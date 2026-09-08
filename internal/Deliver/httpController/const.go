package httpController

const ContentType = "Content-Type"

const Json = "application/json"
const Success = "Success"
const DomainName = "https://filesbes.com/"
const EncryptURLDownload = "https://filesbes.com/d2/"
const UrlDownload = "https://filesbes.com/d/"
const LocalHostName = "http://localhost:8080/"
const Bots = "Bot"
const RequestId = "RequestId"
const InfoPageUrl = "/informationPage"

const MainPageUrl = "/main"
const LoginPage = "/login"
const (
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

const (
	ErrorFileNameEmpty   = "the file's filed is empty"
	ErrorCantGetFileName = "the file's name isn't set"
	ErrorFile            = "can't get a file"
)
const (
	MethodNotAllowed = "the method isn't allowed"
)
const (
	Break    = "break"
	NotStart = "not_start"
)
