package debugger

var appName, appVersion string

func SetAppInfo(name, version string) {
	appName = name
	appVersion = version
}

func getAppName() string {
	return appName
}

func getAppVersion() string {
	return appVersion
}
