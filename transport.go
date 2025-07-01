package lazylog

// Transport define a interface para destinos de log (ex: arquivo, console, etc).
type Transport interface {
	// WriteLog recebe uma Entry e a escreve no destino configurado.
	WriteLog(entry *Entry) error
	// MinLevel retorna o nível mínimo deste transporte.
	MinLevel() Level
}
