package global

var (
	globalValues = map[string]map[string]interface{}{}
)

func Load(name string) map[string]interface{} {
	return globalValues[name]
}

func Save(name string, data map[string]interface{}) {
	globalValues[name] = data
}
