# lazylog

Uma biblioteca de logging para Go inspirada na Winston do NodeJS, com foco em flexibilidade, extensibilidade e facilidade de uso.

## Instalação

```sh
go get github.com/chmenegatti/lazylog
```

## Principais Features

- **Múltiplos transportes**: console, arquivo, rotação de arquivo (lumberjack), syslog, customizáveis
- **Níveis de log customizáveis**: registre seus próprios níveis além de DEBUG, INFO, WARN, ERROR
- **Formatadores customizáveis**: texto, JSON, ou implemente o seu
- **Metadata/contexto extra**: adicione campos extras (ex: user, request_id, etc)
- **Hooks**: execute funções antes/depois de cada log, ou em caso de erro de transporte
- **Filtros por transporte**: lógica customizada para decidir se um log será aceito
- **Configuração dinâmica via struct/map/arquivo (JSON/YAML)**: inicialize o logger a partir de configuração
- **Suporte a context.Context**: integração fácil com tracing/distribuição
- **Formatação customizada por mensagem**: sobrescreva o formatter só para um log
- **Campos aninhados/estruturados**: suporte a mapas dentro de mapas no JSON
- **Remoção/adicionamento dinâmico de transportes**
- **Stacktrace automático**: inclui stacktrace em logs de erro/fatal
- **API de child loggers**: loggers derivados com contexto fixo
- **Métodos Fatal/Panic**: loga e encerra a aplicação ou faz panic
- **Benchmarks e testes automatizados**
- **Exemplos de integração com frameworks web (Gin, Echo, Fiber)**
- **Transporte para syslog**

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

## Configuração via Arquivo (JSON/YAML)

### logger_config.json

```json
{
  "Transports": [
    {
      "Type": "console",
      "Level": "DEBUG",
      "Formatter": "text"
    },
    {
      "Type": "file",
      "Level": "INFO",
      "Formatter": "json",
      "Options": {"path": "app.log"}
    }
  ]
}
```

### logger_config.yaml

```yaml
Transports:
  - Type: console
    Level: DEBUG
    Formatter: text
  - Type: file
    Level: INFO
    Formatter: json
    Options:
      path: app.log
```

### Uso

```go
cfg, _ := lazylog.LoadLoggerConfigJSON("logger_config.json")
logger, _ := lazylog.NewLoggerFromConfig(cfg)
logger.Info("Logger configurado via JSON!")
```

---

## Uso de Child Logger

```go
child := logger.WithFields(map[string]any{"service": "auth", "env": os.Getenv("ENV")})
child.Info("Log do serviço de autenticação")
child.Error("Erro no serviço de autenticação", map[string]any{"code": 401})
```

---

## Envio para Syslog

```go
syslogTransport, _ := lazylog.NewSyslogTransport(syslog.LOG_INFO|syslog.LOG_LOCAL0, "myapp", lazylog.INFO, &lazylog.TextFormatter{})
logger := lazylog.NewLogger(syslogTransport)
logger.Info("Log enviado para o syslog!")
```

---

## Benchmarks

Execute:

```sh
go test -bench=. -benchmem
```

---

## Exemplos de Integração com Frameworks Web

### Gin

```go
r := gin.New()
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{Level: lazylog.DEBUG, Formatter: &lazylog.TextFormatter{}})
r.Use(func(c *gin.Context) {
    start := time.Now()
    c.Next()
    latency := time.Since(start)
    logger.WithFields(map[string]any{
        "method": c.Request.Method,
        "path":   c.Request.URL.Path,
        "status": c.Writer.Status(),
        "latency": latency.String(),
    }).Info("request completed")
})
```

### Echo

```go
e := echo.New()
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{}})
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        err := next(c)
        latency := time.Since(start)
        logger.WithFields(map[string]any{
            "method": c.Request().Method,
            "path":   c.Request().URL.Path,
            "status": c.Response().Status,
            "latency": latency.String(),
        }).Info("request completed")
        return err
    }
})
```

### Fiber

```go
app := fiber.New()
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{}})
app.Use(func(c *fiber.Ctx) error {
    start := time.Now()
    err := c.Next()
    latency := time.Since(start)
    logger.WithFields(map[string]any{
        "method": c.Method(),
        "path":   c.Path(),
        "status": c.Response().StatusCode(),
        "latency": latency.String(),
    }).Info("request completed")
    return err
})
```

---

## Testes e Benchmark

Para rodar os testes:

```sh
go test ./...
```

---

## Para mais exemplos, veja a pasta `examples`

---

## Licença

MIT
