package param

type NotifyType string

const (
	Activity NotifyType = "activity"
	Browser  NotifyType = "browser"
	WebView  NotifyType = "webView"
	Share    NotifyType = "share"
)

type NotifyParam struct {
	Type        NotifyType `json:"type" mapstructure:"type"`
	URL         string     `json:"url" mapstructure:"url"`
	Nav         string     `json:"nav_show" mapstructure:"nav_show"`
	Token       string     `json:"need_token" mapstructure:"need_token"`
	ActIOS      string     `json:"act_ios" mapstructrue:"act_ios"`
	ActAndroid  string     `json:"act_android" mapstructrue:"act_android"`
	Fragment    string     `json:"fragment" mapstructrue:"fragment"`
	DisplayMode string     `json:"displayMode" mapstructrue:"displayMode"`
}
