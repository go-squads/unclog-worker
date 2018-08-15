package alertic

type ConfigStruct struct {
	errorLog map[string]interface{}
	debugLog map[string]interface{}
}

type Alertic struct {
	alerticType string
	appName     string
	logLevel    string
	limit       int
	handler     func()
}

func SendAlert(appName, logLevel string, current int, handler string, receivers []string) {

	if handler == "email" {
		SendEmail(appName, logLevel, current, receivers)
	}
}
