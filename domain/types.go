package domain

type Project struct {
	Name         string
	Path         string
	Description  string
	Architecture string
	GenerateEnv  bool
	GenerateReadME bool
	GenerateGitIgnore bool
}

type Attribute struct{
	Field interface{}
	Label string
}
type Template struct {
	Name string `yaml:"name"`
	Description string `yaml:"description"`
	Folders interface {} `yaml:"folders"`
}

type BoilerPlate struct{
	Name string `yaml:"name"`
	Directory string `yaml:"directory"`
	Description string `yaml:"description"`
	Content string `yaml:"content"`


}