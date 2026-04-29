# Weather MCP Server 🌤️

Um servidor MCP (Model Context Protocol) robusto e profissional desenvolvido em Go para fornecer informações meteorológicas precisas e detalhadas através da API [HG Brasil Weather](https://hgbrasil.com/status/weather).

Este servidor permite que modelos de IA (como o Claude) acessem dados em tempo real sobre condições climáticas, previsões estendidas e até informações astronômicas (fases da lua) para qualquer localização geográfica.

## 🚀 Funcionalidades

- **Previsão por Coordenadas (`get_forecast_lat_lon`)**: Fornece dados meteorológicos completos baseados em latitude e longitude.
- **Previsão por Cidade (`get_forecast_city`)**: Permite buscar o clima diretamente pelo nome da cidade (ex: "São Paulo, SP").
- **Dados em Tempo Real**: Temperatura atual, umidade, velocidade do vento e nebulosidade.
- **Previsão Estendida**: Próximos dias com temperaturas máximas/mínimas e probabilidade de chuva.
- **Informações Visuais**: Links dinâmicos para ícones de condições climáticas e fases da lua.
- **Suporte Multi-Plataforma**: Binário compilado para execução eficiente em diversos sistemas.

## 🛠️ Tecnologias Utilizadas

- **Linguagem**: Go 1.25.3
- **SDK**: [Model Context Protocol Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- **API Externa**: HG Brasil Weather API

## 📋 Pré-requisitos

1.  **Go**: Certifique-se de ter o Go instalado (versão 1.25+ recomendada).
2.  **Node.js**: Necessário para executar o **MCP Inspector** via `npx`.
3.  **Chave de API**: Este projeto utiliza a API da HG Brasil. Para obter sua chave:
    - Acesse o [Guia de Criação de Chave](https://hgbrasil.com/docs/guide/key).
    - Cadastre-se ou faça login no portal.
    - Crie uma nova aplicação para obter sua `API_KEY` gratuita ou premium.

## 🔨 Compilação e Instalação

Para compilar o projeto e gerar o binário executável:

```bash
go build -o weather .
```

## 💻 Como Executar

O servidor pode ser executado passando a chave da API via flag ou variável de ambiente.

### Executando via Código Fonte
```bash
go run main.go -API_KEY {sua_api_key}
```

### Executando o Binário
```bash
./weather -API_KEY {sua_api_key}
```

> **Nota:** Você também pode definir a variável de ambiente `API_KEY` no seu sistema.

## 🔍 Debug e Desenvolvimento

### No VS Code

Como a pasta `.vscode/` é ignorada pelo controle de versão, você deve configurá-la manualmente para habilitar o debug:

1. Crie uma pasta chamada `.vscode` na raiz do projeto.
2. Crie um arquivo chamado `launch.json` dentro dessa pasta.
3. Adicione o seguinte conteúdo ao arquivo:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Weather Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}",
      "args": [
        "-API_KEY",
        "SUA_API_KEY_AQUI"
      ]
    }
  ]
}
```

4. Pressione `F5` para iniciar o debug. A saída aparecerá no **Debug Console**.

### Com MCP Inspector
Para testar e depurar a comunicação do protocolo MCP de forma interativa:

```bash
npx @modelcontextprotocol/inspector go run . -API_KEY {sua_api_key}
```

## 📊 Visualização Gráfica (HTML)

O projeto inclui um arquivo `template_clima.html` que pode ser usado para gerar relatórios meteorológicos visuais e elegantes.

### Como utilizar
Você pode pedir para o modelo de IA (Claude) ler o template e preencher com os dados obtidos pelo servidor MCP.

**Exemplo de Prompt:**
> "Com o modelo `template_clima.html`, poderia verificar o clima e previsão para os próximos dias para a **lat -23.550196048274127** e **lon -46.63396252883539**? Mostre as informações geradas em um arquivo `.html` a partir do template informado."

O resultado será um arquivo HTML estilizado com glassmorphism, gradientes modernos e ícones dinâmicos.

## ⚙️ Configuração no Claude Desktop

Para utilizar este servidor no Claude Desktop, adicione o seguinte ao seu arquivo `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "weather-server": {
      "command": "/caminho/para/o/projeto/weather",
      "args": ["-API_KEY", "SUA_API_KEY"]
    }
  }
}
```

## 📂 Estrutura do Projeto

- `main.go`: Ponto de entrada do servidor, definição de ferramentas e integração com a API.
- `go.mod`: Gerenciamento de dependências.

---
**Desenvolvido por [Cleberson Pauluci](https://github.com/cpauluci)**
