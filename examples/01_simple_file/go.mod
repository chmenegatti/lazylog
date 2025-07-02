module github.com/chmenegatti/lazylog/examples/simple_file

go 1.23.2 // Ou a versão do Go que você está usando

require github.com/chmenegatti/lazylog v0.0.0

require (
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/chmenegatti/lazylog => ../../
