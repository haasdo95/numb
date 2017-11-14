package bootstrap

type NmbConfig struct {
	Name  string   `json:"name"`
	Train []string `json:"train"`
	Test  []string `json:"test"`
}
