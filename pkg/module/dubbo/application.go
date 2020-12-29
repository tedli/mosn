package dubbo

var (
	application string
)

func GetApplicationName() string {
	return application
}

func SetApplicationName(name string) {
	application = name
}
