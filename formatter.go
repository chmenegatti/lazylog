package lazylog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Formatter define a interface para formatadores de log.
// Qualquer tipo que implementar o método Format pode ser usado como um formatador.
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// --- Implementação do TextFormatter ---

// TextFormatter formata logs como texto simples.
type TextFormatter struct {
	// TimestampFormat especifica o formato do timestamp. Usa time.RFC3339 se vazio.
	TimestampFormat string
}

// Format implementa a interface Formatter para TextFormatter.
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	var b bytes.Buffer

	// Escreve o timestamp formatado
	b.WriteString(entry.Timestamp.Format(timestampFormat))

	// Adiciona um espaço
	b.WriteString(" ")

	// Escreve o nível
	b.WriteString(fmt.Sprintf("[%s]", entry.Level.String()))

	// Adiciona outro espaço
	b.WriteString(" ")

	// Escreve a mensagem
	b.WriteString(entry.Message)
	if entry.Fields != nil && len(entry.Fields) > 0 {
		b.WriteString(" ")
		for k, v := range entry.Fields {
			b.WriteString(fmt.Sprintf("%s=%v ", k, v))
		}
	}
	// Adiciona uma nova linha no final
	b.WriteString("\n")

	return b.Bytes(), nil
}

// --- Implementação do JSONFormatter ---

// JSONFormatter formata logs como JSON.
type JSONFormatter struct{}

// Format implementa a interface Formatter para JSONFormatter.
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	// Para serializar o nível como string, criamos um tipo anônimo.
	data := map[string]interface{}{
		"timestamp": entry.Timestamp.Format(time.RFC3339Nano), // JSON geralmente usa alta precisão
		"level":     entry.Level.String(),
		"message":   entry.Message,
	}
	for k, v := range entry.Fields {
		data[k] = v
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Adiciona uma nova linha para que cada log JSON fique em sua própria linha
	return append(b, '\n'), nil
}
