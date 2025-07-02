# lazylog

Uma biblioteca de logging para Go inspirada na Winston do NodeJS, com foco em flexibilidade, extensibilidade e facilidade de uso.

## Instalação

```sh
go get github.com/chmenegatti/lazylog
```

## Principais Features

- **Múltiplos transportes**: console, arquivo, rotação de arquivo (lumberjack), customizáveis
- **Níveis de log customizáveis**: registre seus próprios níveis além de DEBUG, INFO, WARN, ERROR
- **Formatadores customizáveis**: texto, JSON, ou implemente o seu
- **Metadata/contexto extra**: adicione campos extras (ex: user, request_id, etc)
- **Hooks**: execute funções antes/depois de cada log, ou em caso de erro de transporte
- **Filtros por transporte**: lógica customizada para decidir se um log será aceito
- **Configuração dinâmica via struct/map**: inicialize o logger a partir de configuração
- **Suporte a context.Context**: integração fácil com tracing/distribuição
- **Formatação customizada por mensagem**: sobrescreva o formatter só para um log
- **Campos aninhados/estruturados**: suporte a mapas dentro de mapas no JSON
- **Remoção/adicionamento dinâmico de transportes**

---

## Uso Básico

```go
import "github.com/chmenegatti/lazylog"

func main() {
 logger := lazylog.NewLogger(
  &lazylog.ConsoleTransport{
   Level:     lazylog.DEBUG,
   Formatter: &lazylog.TextFormatter{},
  },
 )
 logger.Info("Servidor iniciado.")
 logger.Warn("A conexão está lenta.")
 logger.Error("Falha ao processar requisição.")
}
```

**Saída:**

```bash
2025-07-01T21:00:00-03:00 [INFO] Servidor iniciado.
2025-07-01T21:00:00-03:00 [WARN] A conexão está lenta.
2025-07-01T21:00:00-03:00 [ERROR] Falha ao processar requisição.
```

---

## Múltiplos Transportes (console + arquivo)

```go
fileTransport, _ := lazylog.NewFileTransport("app.log", lazylog.INFO, &lazylog.JSONFormatter{})
logger := lazylog.NewLogger(
 &lazylog.ConsoleTransport{Level: lazylog.DEBUG, Formatter: &lazylog.TextFormatter{}},
 fileTransport,
)
logger.Info("Mensagem de info (console e arquivo)")
logger.Warn("Mensagem de aviso")
logger.Error("Mensagem de erro")
logger.ComFields(map[string]any{"user": "cesar", "request_id": 12345}).Info("Log com contexto extra")
```

**Saída no console:**

```bash
2025-07-01T21:00:00-03:00 [INFO] Mensagem de info (console e arquivo)
2025-07-01T21:00:00-03:00 [WARN] Mensagem de aviso
2025-07-01T21:00:00-03:00 [ERROR] Mensagem de erro
2025-07-01T21:00:00-03:00 [INFO] Log com contexto extra user=cesar request_id=12345
```

**Saída no arquivo (JSON):**

```json
{"level":"INFO","message":"Mensagem de info (console e arquivo)","timestamp":"2025-07-01T21:00:00-03:00"}
{"level":"WARN","message":"Mensagem de aviso","timestamp":"2025-07-01T21:00:00-03:00"}
{"level":"ERROR","message":"Mensagem de erro","timestamp":"2025-07-01T21:00:00-03:00"}
{"level":"INFO","message":"Log com contexto extra","request_id":12345,"timestamp":"2025-07-01T21:00:00-03:00","user":"cesar"}
```

---

## Rotação de Arquivo (Lumberjack)

```go
rotating := lazylog.NewLumberjackTransport(
 "app_rotating.log", lazylog.INFO, &lazylog.JSONFormatter{},
 maxSize=1, maxBackups=3, maxAge=7, compress=true,
)
logger := lazylog.NewLogger(rotating)
logger.Info("Log com rotação!")
```

---

## Hooks (antes, depois, erro)

```go
logger.AddHook(func(e *lazylog.Entry) {
 if e.Fields == nil {
  e.Fields = make(map[string]any)
 }
 e.Fields["runtime"] = time.Now().UnixNano()
}, true) // before

logger.AddHook(func(e *lazylog.Entry) {
 if e.Level == lazylog.ERROR {
  fmt.Println("[ALERTA] Um erro foi registrado:", e.Message)
 }
}, false) // after

logger.AddErrorHook(func(e *lazylog.Entry, t lazylog.Transport, err error) {
 fmt.Println("[ERRO] Falha ao gravar log:", err)
})
```

---

## Filtros por transporte

```go
// Só loga mensagens que contenham "importante"
filter := func(e *lazylog.Entry) bool {
 return strings.Contains(e.Message, "importante")
}
logger.AddTransport(&lazylog.TransportWithFilter{
 Transport: &lazylog.ConsoleTransport{Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{}},
 Filter:    filter,
})
```

---

## Níveis customizados

```go
lazylog.RegisterLevel("NOTICE", 10)
logger.Info("Log normal")
logger.log(lazylog.ParseLevel("NOTICE"), "Log de nível NOTICE")
```

---

## Logging com context.Context

```go
ctx := context.WithValue(context.Background(), "trace_id", "abc-123")
logger.InfoCtx(ctx, "Log com trace", nil)
```

---

## Configuração dinâmica via struct/map

```go
cfg := lazylog.LoggerConfig{
 Transports: []lazylog.TransportConfig{
  {
   Type:      "console",
   Level:     "DEBUG",
   Formatter: "text",
  },
  {
   Type:      "file",
   Level:     "INFO",
   Formatter: "json",
   Options: map[string]any{"path": "app.log"},
  },
 },
}
logger, _ := lazylog.NewLoggerFromConfig(cfg)
```

---

## Formatação customizada por mensagem

```go
logger.WithFormatter(&lazylog.TextFormatter{TimestampFormat: time.RFC822}).Info("Log com timestamp customizado!")
```

---

## Campos aninhados/estruturados

```go
logger.ComFields(map[string]any{
 "user": "cesar",
 "request": map[string]any{
  "id": 123,
  "ip": "1.2.3.4",
 },
}).Info("Log estruturado")
```

**Saída JSON:**

```json
{"level":"INFO","message":"Log estruturado","timestamp":"2025-07-01T21:00:00-03:00","user":"cesar","request":{"id":123,"ip":"1.2.3.4"}}
```

---

## Testes e Benchmark

Para rodar os testes:

```sh
go test ./...
```

---

## Integração com frameworks web

Basta usar hooks para capturar request_id, user, etc, e passar context.Context nos logs.

---

## Contribuindo

Pull requests são bem-vindos! Abra uma issue para bugs ou sugestões.

---

## Licença

MIT
