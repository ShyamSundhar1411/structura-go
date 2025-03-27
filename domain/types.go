package domain

type Project struct {
	Name         string
	Path         string
	PackageName  string
	Description  string
	Architecture string
	Dependencies []Dependency
}
type Attribute struct {
	Field        *string
	Label        string
	Type         string
	Options      interface{}
	Condition    func() bool
	DefaultValue string
}
type Template struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Readme 		string 		`yaml:"readme"`
	Folders     interface{} `yaml:"folders"`
}

type BoilerPlate struct {
	Name        string `yaml:"name"`
	Directory   string `yaml:"directory"`
	Description string `yaml:"description"`
	Content     string `yaml:"content"`
}
type FileContent struct {
	Files map[string]string `yaml:",inline"`
}

type Dependency struct {
	Name        string                 `yaml:"name"`
	Source      string                 `yaml:"source"`
	Description string                 `yaml:"description"`
	Version     string                 `yaml:"version"`
	Content     map[string]FileContent `yaml:"content,omitempty"`
}
