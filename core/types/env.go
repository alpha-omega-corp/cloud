package types

type ConfigUserService struct {
	JwtSecret string `mapstruct:"secret"`
}

type ConfigConfigService struct {
	Mode string `mapstruct:"mode"`
}

type EnvVars struct {
	Host    string `mapstruct:"host"`
	Config  string `mapstruct:"config"`
	Service string `mapstruct:"name"`
}

type Environment struct {
	Vars   EnvVars
	Config Config
	Host   ConfigHost
}
