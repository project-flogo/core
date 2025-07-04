package debugger

var appName, appVersion string

func SetAppInfo(name, version string) {
	appName = name
	appVersion = version
}

func GetAppName() string {
	return appName
}

func GetAppVersion() string {
	return appVersion
}
