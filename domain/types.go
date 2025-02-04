package domain

type Project struct {
	Name         string
	Path         string
	Description  string
	Architecture string
}

type Attribute struct{
	Field *string
	Label string
}
type Template struct {
	Architecture string `yaml:"architecture"`
	Description string `yaml:"description"`
	Folders interface {} `yaml:"folders"`
}