package internal

type Track struct {
	Id       int
	Path     string
	Filename string
}

type SoundSwitchVersion struct {
	Build  int
	Hotfix int
	Major  int
	Minor  int
}

type SoundSwitchProject struct {
	Id       string
	ReadOnly bool
	Version  SoundSwitchVersion
}
