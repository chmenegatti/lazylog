# lazylog 🚀

<p align="center">
  <img src="logo/lazylog.png" alt="lazylog logo" width="240"/>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/chmenegatti/lazylog.svg)](https://pkg.go.dev/github.com/chmenegatti/lazylog)
[![Go Report Card](https://goreportcard.com/badge/github.com/chmenegatti/lazylog)](https://goreportcard.com/report/github.com/chmenegatti/lazylog)
[![License: Apache 2.0](https://img.shields.io/badge/apache-2.0-yellow)](LICENSE.md)

Uma biblioteca de logging para Go inspirada na Winston do NodeJS, com foco em flexibilidade, extensibilidade e facilidade de uso.

---

## ✨ Principais Features

- � **Thread-safe**: uso seguro em goroutines concorrentes (protegido por `sync.RWMutex`)
- �🛣️ **Múltiplos transportes**: console, arquivo, rotação de arquivo (lumberjack), syslog, customizáveis
- 🏷️ **Níveis de log customizáveis**: registre seus próprios níveis além de DEBUG, INFO, WARN, ERROR
- 🎨 **Formatadores customizáveis**: texto, JSON, emojis ou implemente o seu
- 🧩 **Metadata/contexto extra**: adicione campos extras (ex: user, request_id, etc)
- 🪝 **Hooks**: execute funções antes/depois de cada log, ou em caso de erro de transporte
- 🧹 **Filtros por transporte**: lógica customizada para decidir se um log será aceito
- ⚙️ **Configuração dinâmica via struct/map/arquivo (JSON/YAML)**
- 🧵 **Suporte a context.Context**: integração fácil com tracing/distribuição
- 🖌️ **Formatação customizada por mensagem**
- 🗂️ **Campos aninhados/estruturados**
- ➕ **Remoção/adicionamento dinâmico de transportes**
- 🪓 **Stacktrace automático**
- 👶 **Child loggers** com contexto fixo
- 💥 **Métodos Fatal/Panic** com stacktrace
- 🔌 **Logger.Close()** fecha todos os transportes automaticamente
- 🏎️ **Benchmarks e testes automatizados com race detector**
- 🌐 **Exemplos de integração com frameworks web (Gin, Echo, Fiber)**
- 🖥️ **Transporte para syslog** com mapeamento correto de níveis

---

## � Instalação

```sh
go get github.com/chmenegatti/lazylog@latest
```

---

## 🚀 Início Rápido

```go
package main

import "github.com/chmenegatti/lazylog"

func main() {
    logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
        Level:     lazylog.INFO,
        Formatter: &lazylog.TextFormatter{},
    })
    defer logger.Close()

    logger.Info("Hello, world!")
    logger.Warn("Cuidado!")
    logger.Error("Algo deu errado!")
}
```

---

## 📚 Exemplos de Uso

### Múltiplos Transportes (Console + Arquivo)

```go
fileTransport, _ := lazylog.NewFileTransport("app.log", lazylog.INFO, &lazylog.JSONFormatter{})

logger := lazylog.NewLogger(
    &lazylog.ConsoleTransport{Level: lazylog.DEBUG, Formatter: &lazylog.TextFormatter{}},
    fileTransport,
)
defer logger.Close()

logger.Debug("Aparece só no console")
logger.Info("Vai para console e arquivo")
```

---

### Rotação de Arquivo (Lumberjack)

```go
lj := lazylog.NewLumberjackTransport("app.log", lazylog.INFO, &lazylog.TextFormatter{}, 10, 3, 7, true)

logger := lazylog.NewLogger(lj)
defer logger.Close()

logger.Info("Log com rotação automática!")
```

---

### Metadata/Contexto Extra (Fields)

```go
logger.ComFields(map[string]interface{}{
    "user":       "bob",
    "request_id": "abc123",
}).Info("Log com contexto")
```

---

### Campos Aninhados (JSON)

```go
logger.ComFields(map[string]interface{}{
    "user": "cesar",
    "request": map[string]interface{}{
        "id": 123,
        "ip": "1.2.3.4",
    },
}).Info("Log com campos aninhados")
```

---

### Child Logger (Contexto fixo)

Ideal para microserviços — cria loggers derivados com campos que são incluídos automaticamente em toda mensagem:

```go
baseLogger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{},
})

authLogger := baseLogger.WithFields(map[string]any{"service": "auth"})
paymentLogger := baseLogger.WithFields(map[string]any{"service": "payment"})

authLogger.Info("Usuário autenticado")
paymentLogger.Error("Falha no pagamento", map[string]any{"code": 500})
```

---

### Hooks (Before / After / Error)

```go
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{},
})

// Hook executado ANTES de cada log
logger.AddHook(func(e *lazylog.Entry) {
    if e.Fields == nil {
        e.Fields = make(map[string]interface{})
    }
    e.Fields["app"] = "meu-servico"
}, true) // true = before

// Hook executado DEPOIS de cada log
logger.AddHook(func(e *lazylog.Entry) {
    fmt.Println("Log registrado com sucesso!")
}, false) // false = after

// Hook para erros de transporte
logger.AddErrorHook(func(e *lazylog.Entry, t lazylog.Transport, err error) {
    fmt.Printf("Erro ao gravar log: %v\n", err)
})

logger.Info("Testando hooks!")
```

---

### Filtros por Transporte

```go
tr := &lazylog.WriterTransport{
    Writer:    os.Stdout,
    Level:     lazylog.INFO,
    Formatter: &lazylog.TextFormatter{},
}

filter := func(entry *lazylog.Entry) bool {
    return strings.Contains(entry.Message, "importante")
}

logger := lazylog.NewLogger(&lazylog.TransportWithFilter{
    Transport: tr,
    Filter:    filter,
})

logger.Info("este log não vai aparecer")
logger.Info("log importante!")  // ✅ este sim
```

---

### Formatação Customizada por Mensagem

```go
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{},
})

// Usa um formatter diferente apenas para esta mensagem
logger.WithFormatter(&lazylog.JSONFormatter{}).Info("Este log sai em JSON!")
```

---

### 😃 Logs com Emojis (EmojiFormatter)

O `EmojiFormatter` adiciona emojis automaticamente conforme o nível do log:

```go
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.DEBUG,
    Formatter: &lazylog.EmojiFormatter{Base: &lazylog.TextFormatter{}},
})

logger.Debug("Debugando...")   // 🐛 Debugando...
logger.Info("Tudo certo!")     // ℹ️ Tudo certo!
logger.Warn("Atenção!")        // ⚠️ Atenção!
logger.Error("Deu ruim!")      // ❌ Deu ruim!
```

---

### Stacktrace Automático

```go
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.DEBUG, Formatter: &lazylog.TextFormatter{},
})

// Ativa stacktrace para ERROR (inclui stack no campo "stacktrace")
logger.EnableStacktrace(lazylog.ERROR)

logger.Error("Erro grave — stacktrace será incluído automaticamente!")
```

---

### Métodos Fatal e Panic

Ambos incluem stacktrace automaticamente no log antes de encerrar/panic:

```go
// Fatal: loga com stacktrace e chama os.Exit(1)
logger.Fatal("Erro fatal!", map[string]any{"code": 500})

// Panic: loga com stacktrace e chama panic()
logger.Panic("Erro crítico!", map[string]any{"reason": "null pointer"})
```

---

### Níveis Customizados

```go
const TRACE lazylog.Level = 5
lazylog.RegisterLevel("TRACE", TRACE)

fmt.Println(lazylog.ParseLevel("TRACE")) // 5
```

---

### Suporte a Context (Tracing)

Extrai `trace_id` automaticamente do `context.Context`:

```go
// Usando string como chave
ctx := context.WithValue(context.Background(), "trace_id", "abc-123")
logger.InfoCtx(ctx, "Request recebida", nil)
// Output: ... trace_id=abc-123

// Usando o tipo exportado CtxKey (recomendado)
ctx = context.WithValue(context.Background(), lazylog.CtxKey("trace_id"), "xyz-789")
logger.InfoCtx(ctx, "Request processada", map[string]interface{}{
    "method": "GET",
    "path":   "/api/users",
})
```

---

## ⚙️ Configuração via Arquivo (JSON/YAML)

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
cfg, err := lazylog.LoadLoggerConfigJSON("logger_config.json")
if err != nil {
    log.Fatal(err)
}
logger, err := lazylog.NewLoggerFromConfig(cfg)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

logger.Info("Logger configurado via JSON!")
```

---

## 🖥️ Envio para Syslog

O `SyslogTransport` mapeia os níveis automaticamente para a severity correta do syslog (`Debug`, `Info`, `Warning`, `Err`):

```go
syslogTransport, err := lazylog.NewSyslogTransport(
    syslog.LOG_INFO|syslog.LOG_LOCAL0, "myapp",
    lazylog.DEBUG, &lazylog.TextFormatter{},
)
if err != nil {
    log.Fatal(err)
}

logger := lazylog.NewLogger(syslogTransport)
defer logger.Close()

logger.Debug("vai como syslog.Debug()")
logger.Info("vai como syslog.Info()")
logger.Warn("vai como syslog.Warning()")
logger.Error("vai como syslog.Err()")
```

---

## 🔌 Logger.Close()

Fecha todos os transportes que implementam `io.Closer` (FileTransport, LumberjackTransport, SyslogTransport):

```go
logger := lazylog.NewLogger(fileTransport, lumberjackTransport, consoleTransport)
defer logger.Close() // fecha file e lumberjack; console não precisa fechar
```

---

## 🚦 Integração com Frameworks Web

Os exemplos de integração com frameworks estão em módulos separados dentro de `examples/`:

### Gin

```go
r := gin.New()
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.DEBUG, Formatter: &lazylog.TextFormatter{},
})

r.Use(func(c *gin.Context) {
    start := time.Now()
    c.Next()
    latency := time.Since(start)
    logger.WithFields(map[string]any{
        "method":  c.Request.Method,
        "path":    c.Request.URL.Path,
        "status":  c.Writer.Status(),
        "latency": latency.String(),
    }).Info("request completed")
})
```

> Veja o exemplo completo em [`examples/05_gin`](examples/05_gin/main.go)

### Echo

```go
e := echo.New()
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{},
})

e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        err := next(c)
        latency := time.Since(start)
        logger.WithFields(map[string]any{
            "method":  c.Request().Method,
            "path":    c.Request().URL.Path,
            "status":  c.Response().Status,
            "latency": latency.String(),
        }).Info("request completed")
        return err
    }
})
```

> Veja o exemplo completo em [`examples/06_echo`](examples/06_echo/main.go)

### Fiber

```go
app := fiber.New()
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{},
})

app.Use(func(c *fiber.Ctx) error {
    start := time.Now()
    err := c.Next()
    latency := time.Since(start)
    logger.WithFields(map[string]any{
        "method":  c.Method(),
        "path":    c.Path(),
        "status":  c.Response().StatusCode(),
        "latency": latency.String(),
    }).Info("request completed")
    return err
})
```

> Veja o exemplo completo em [`examples/07_fiber`](examples/07_fiber/main.go)

> **Nota**: Os exemplos de frameworks são módulos Go separados. Para executá-los, entre na pasta do exemplo e rode `go mod tidy && go run .`

---

## 🏎️ Benchmarks

```sh
go test -bench=. -benchmem
```

Para validar thread-safety:

```sh
go test -race -v ./...
```

---

## 📚 Para mais exemplos, veja a pasta [`examples/`](examples/)

---

## 📝 Licença

Apache 2.0 License. Veja o arquivo [LICENSE](LICENSE.md) para mais detalhes.

## 🤝 Como Contribuir

Contribuições são muito bem-vindas! Siga as etapas abaixo para colaborar com o desenvolvimento do lazylog:

1. **Fork o repositório**
   - Clique em "Fork" no topo da página do GitHub para criar uma cópia do projeto no seu perfil.

2. **Clone o seu fork**

   ```sh
   git clone https://github.com/seu-usuario/lazylog.git
   cd lazylog
   ```

3. **Crie uma branch para sua feature/correção**

   ```sh
   git checkout -b minha-feature
   ```

4. **Implemente sua melhoria**
   - Siga o padrão de código e comentários do projeto.
   - Adicione testes automatizados para novas funcionalidades.
   - Atualize a documentação e exemplos, se necessário.

5. **Rode os testes com o race detector**

   ```sh
   go test -race -v ./...
   go test -race -bench=. -benchmem ./...
   ```

6. **Faça commit e push das alterações**

   ```sh
   git add .
   git commit -m "feat: descreva sua feature/correção"
   git push origin minha-feature
   ```

7. **Abra um Pull Request**
   - Acesse o repositório original e clique em "New Pull Request".
   - Descreva claramente sua contribuição.

### Dicas

- Use mensagens de commit claras e objetivas.
- Mantenha as dependências atualizadas e evite adicionar dependências desnecessárias.
- Para grandes mudanças, abra uma issue antes para discutir a proposta.

---
