package config

const (
	AppName                  = "_canu"
	DefaultRelativePath      = ".aws/config"
	CanuSaveProfileFileName  = ".canu"
	DefaultShellConfigScript = ".bash_profile"
	DefaultAliasName         = "canu"
)

type Config struct {
	ConfigPath   string
	UserHomePath string
	Include      []string
	Exclude      []string

	// install flags
	IsSkipAlias           bool
	ShellConfigScriptPath string
	AliasName             string
}
