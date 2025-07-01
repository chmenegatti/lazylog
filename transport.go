package lazylog

// Transport define a interface para destinos de log (ex: arquivo, console, etc).
type Transport interface {
	// WriteLog recebe uma Entry e a escreve no destino configurado.
	WriteLog(entry *Entry) error
	// MinLevel retorna o nível mínimo deste transporte.
	MinLevel() Level
}

// FilterFunc permite lógica customizada para decidir se um log deve ser aceito pelo transporte.
type FilterFunc func(entry *Entry) bool

// TransportWithFilter é um wrapper para adicionar filtro a qualquer transporte.
type TransportWithFilter struct {
	Transport Transport
	Filter    FilterFunc
}

func (t *TransportWithFilter) WriteLog(entry *Entry) error {
	if t.Filter != nil && !t.Filter(entry) {
		return nil // Ignora o log
	}
	return t.Transport.WriteLog(entry)
}

func (t *TransportWithFilter) MinLevel() Level {
	return t.Transport.MinLevel()
}
