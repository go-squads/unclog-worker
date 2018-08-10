package alertic

type ConfigStruct struct {
	errorLog map[string]interface{}
	debugLog map[string]interface{}
}

func SendAlert(appName, logLevel string, current int, handler string) {

	if handler == "email" {
		SendEmail(appName, logLevel, current)
	}
}
