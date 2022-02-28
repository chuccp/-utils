package coke_log

import (
	"io"
	"log"
	"os"
)

type Config struct {
	level     Level //设置显示等级
	Out       io.Writer
	formatter *LogFormatter
	levelMap  map[Level]*cut
}

func NewConfig(Out io.Writer, Formatter *LogFormatter) *Config {
	return &Config{Out: Out, formatter: Formatter, levelMap: make(map[Level]*cut), level: TraceLevel}
}
func (config *Config) SetLevel(level Level) {
	config.level = level
}
func (config *Config) SetFormatter(level Level) {
	config.level = level
}
func (config *Config) GetCut(level Level)*cut {
	return config.levelMap[level]
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
func (config *Config) AddFileConfig(filePattern string, level ...Level) {
	for _, v := range level {
		cut, err := parse(filePattern)
		if err != nil {
			log.Panicln("解析错误：",filePattern)
		}
		config.levelMap[v] = cut
	}
}

var defaultConfig = NewConfig(os.Stdout, defaultFormatter)

func GetDefaultConfig() *Config {
	return defaultConfig
}
