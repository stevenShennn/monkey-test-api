package types

// ParamConfig 参数配置
type ParamConfig struct {
	String StringParamConfig `toml:"string"`
	Number NumberParamConfig `toml:"number"`
	Bool   BoolParamConfig   `toml:"bool"`
}

// StringParamConfig 字符串参数配置
type StringParamConfig struct {
	Tests []ParamTest `toml:"tests"`
}

// NumberParamConfig 数字参数配置
type NumberParamConfig struct {
	Tests []ParamTest `toml:"tests"`
}

// BoolParamConfig 布尔参数配置
type BoolParamConfig struct {
	Tests []ParamTest `toml:"tests"`
}

// ParamTest 参数测试用例
type ParamTest struct {
	Value       interface{} `toml:"value"`
	Description string      `toml:"description"`
	Reason      string      `toml:"reason"`
	Risk        string      `toml:"risk"`
} 