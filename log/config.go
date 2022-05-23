package log

type Config struct {
	level     Level //设置显示等级
	formatter *LogFormatter
	fileLevel Level
	filePattern string

}

func NewConfig(Formatter *LogFormatter) *Config {
	return &Config{ formatter: Formatter, level: ErrorLevel,fileLevel: ErrorLevel}
}
func (config *Config) SetLevel(level Level) {
	config.level = level
}
func (config *Config) SetFormatter(level Level) {
	config.level = level
}


/*AddFileConfig
按行数，按日志，按日期 切割
 FilePattern 规则 ${time:2006-01-02|15-04}-${line:2000}-${size:200mb}-${level}.log

${time:2006-01-02|15-04} 按日期切割

${line:2000}按行数切割

${size:200mb}按尺寸切割

${level}日志类型

哪一个条件先达到就以那一条件为准切割
*/
func (config *Config) AddFileConfig(filePattern string, level Level)error {
	config.filePattern = filePattern
	config.fileLevel = level
	return nil
}

var defaultConfig = NewConfig(defaultFormatter)

func GetDefaultConfig() *Config {
	return defaultConfig
}
