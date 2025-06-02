module github.com/chmenegatti/lazylog/examples/rotate_file

go 1.23.2 // Ou a versão do Go que você está usando

require (
	github.com/chmenegatti/lazylog v0.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

replace github.com/chmenegatti/lazylog => ../../
