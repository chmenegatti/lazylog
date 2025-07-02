# lazylog 🚀

[![Go Reference](https://pkg.go.dev/badge/github.com/chmenegatti/lazylog.svg)](https://pkg.go.dev/github.com/chmenegatti/lazylog)
[![Go Report Card](https://goreportcard.com/badge/github.com/chmenegatti/lazylog)](https://goreportcard.com/report/github.com/chmenegatti/lazylog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build](https://github.com/chmenegatti/lazylog/actions/workflows/go.yml/badge.svg)](https://github.com/chmenegatti/lazylog/actions)

Uma biblioteca de logging para Go inspirada na Winston do NodeJS, com foco em flexibilidade, extensibilidade e facilidade de uso.

---

## ✨ Principais Features

- 🛣️ **Múltiplos transportes**: console, arquivo, rotação de arquivo (lumberjack), syslog, customizáveis
- 🏷️ **Níveis de log customizáveis**: registre seus próprios níveis além de DEBUG, INFO, WARN, ERROR
- 🎨 **Formatadores customizáveis**: texto, JSON, ou implemente o seu
- 🧩 **Metadata/contexto extra**: adicione campos extras (ex: user, request_id, etc)
- 🪝 **Hooks**: execute funções antes/depois de cada log, ou em caso de erro de transporte
- 🧹 **Filtros por transporte**: lógica customizada para decidir se um log será aceito
- ⚙️ **Configuração dinâmica via struct/map/arquivo (JSON/YAML)**
- 🧵 **Suporte a context.Context**: integração fácil com tracing/distribuição
- 🖌️ **Formatação customizada por mensagem**
- 🗂️ **Campos aninhados/estruturados**
- ➕ **Remoção/adicionamento dinâmico de transportes**
- 🪓 **Stacktrace automático**
- 👶 **API de child loggers**
- 💥 **Métodos Fatal/Panic**
- 🏎️ **Benchmarks e testes automatizados**
- 🌐 **Exemplos de integração com frameworks web (Gin, Echo, Fiber)**
- 🖥️ **Transporte para syslog**

---

## 🚦 Exemplos de Integração com Frameworks Web

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
cfg, _ := lazylog.LoadLoggerConfigJSON("logger_config.json")
logger, _ := lazylog.NewLoggerFromConfig(cfg)
logger.Info("Logger configurado via JSON!")
```

---

## 👶 Uso de Child Logger

```go
child := logger.WithFields(map[string]any{"service": "auth", "env": os.Getenv("ENV")})
child.Info("Log do serviço de autenticação")
child.Error("Erro no serviço de autenticação", map[string]any{"code": 401})
```

---

## 🖥️ Envio para Syslog

```go
syslogTransport, _ := lazylog.NewSyslogTransport(syslog.LOG_INFO|syslog.LOG_LOCAL0, "myapp", lazylog.INFO, &lazylog.TextFormatter{})
logger := lazylog.NewLogger(syslogTransport)
logger.Info("Log enviado para o syslog!")
```

---

## 😃 Logs com Emojis (EmojiFormatter)

O `EmojiFormatter` adiciona emojis automaticamente conforme o nível do log, tornando a leitura mais divertida e visual:

```go
logger := lazylog.NewLogger(&lazylog.ConsoleTransport{
    Level: lazylog.DEBUG,
    Formatter: &lazylog.EmojiFormatter{},
})

logger.Debug("Debugando...")   // 🐛 Debugando...
logger.Info("Tudo certo!")     // ℹ️ Tudo certo!
logger.Warn("Atenção!")        // ⚠️ Atenção!
logger.Error("Deu ruim!")      // ❌ Deu ruim!
```

---

## 🏎️ Benchmarks

Execute:

```sh
go test -bench=. -benchmem
```

---

## 📚 Para mais exemplos, veja a pasta `examples`

---

## 📝 Licença

MIT License. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

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

5. **Rode os testes e benchmarks**

   ```sh
   go test ./... -v
   go test -bench=. -benchmem
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
