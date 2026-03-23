
---

# Honeytoken System

Este projeto consiste em um sistema de decepção (Honeypot) e monitoramento de ativos (SOC Dashboard) desenvolvido em Go e Python. O objetivo é simular serviços vulneráveis para capturar tentativas de intrusão, analisar payloads em tempo real e visualizar ameaças em um painel de controle centralizado.

## Arquitetura do Sistema

O sistema é dividido em duas camadas principais:

### 1. Engine de Captura (Go)
O binário `SecurityMonitor.exe` gerencia a superfície de ataque:
* **Serviços TCP:** Escuta nas portas 21 (FTP), 22 (SSH) e 3306 (MySQL), enviando banners de serviço falsos para ferramentas de varredura.
* **Serviço Web:** Escuta na porta 8080, simulando um servidor Apache 2.4.41. Monitora rotas críticas como `/.env` e `/admin`.
* **Mecanismo de Tarpit:** Implementa atrasos intencionais na resposta para scanners automatizados.
* **Persistência de Log:** Utiliza escrita com `f.Sync()` para garantir que os dados de auditoria sejam disponibilizados imediatamente no Windows para leitura compartilhada.

### 2. Dashboard SOC (Python)
O script `monitor.py` atua como a interface de análise:
* **Leitura não bloqueante:** Processa o arquivo `system_audit.log` sem interromper o serviço de captura.
* **Filtragem por Regex:** Extrai metadados como Timestamp, Porta, Protocolo, IP de origem e Payload/Detalhes.
* **Interface Visual:** Utiliza a biblioteca Rich para renderização de tabelas dinâmicas com atualização em tempo real (4 FPS).

## Estrutura de Arquivos

* `server.go`: Código fonte do servidor de captura em Go.
* `monitor.py`: Dashboard de visualização em Python.
* `system_audit.log`: Arquivo de log unificado (gerado em tempo real).
* `.gitignore`: Configurado para ignorar binários e arquivos de log sensíveis.

## Requisitos e Instalação

### Pré-requisitos
* Go 1.20 ou superior.
* Python 3.10 ou superior.
* Bibliotecas Python: `rich`.

### Compilação (Windows)
Para gerar o executável em modo oculto (sem janela de console):
```powershell
go build -ldflags "-H=windowsgui -s -w" -o SecurityMonitor.exe server.go
```

### Execução
1. Inicie o serviço de captura:
   ```powershell
   .\SecurityMonitor.exe
   ```
2. Inicie o dashboard de monitoramento:
   ```powershell
   python monitor.py
   ```

## Logs e Auditoria
Os logs seguem o formato unificado para integração com outras ferramentas de SIEM:
`[HH:MM:SS] PORT:XX | PROTO:XX | IP:XX | DETAIL:XX`

## Isenção de Responsabilidade
Este projeto foi desenvolvido para fins de pesquisa em cibersegurança e laboratórios controlados. O uso desta ferramenta em redes sem autorização prévia é de inteira responsabilidade do usuário.

---

---

## Integrações

Para complementar o ecossistema de segurança e tornar o **Honeytoken System** uma ferramenta de defesa ativa, as seguintes integrações técnicas podem ser implementadas:

### 1. Defesa Ativa: Bloqueio Automático (Firewall)
Esta integração transforma o sistema de monitoramento em uma ferramenta de prevenção. Um script auxiliar pode ler o arquivo `system_audit.log` e, ao detectar múltiplas tentativas de um mesmo IP, criar automaticamente uma regra no Firewall do Windows para banir o endereço de origem.
* **Mecanismo**: Monitoramento de frequência de logs por IP.
* **Ação**: Execução do comando `New-NetFirewallRule` para bloqueio imediato.

### 2. Alertas Remotos (Telegram / Discord)
Essencial para monitoramento em tempo real sem a necessidade de estar diante do terminal. O dashboard Python pode enviar um *Webhook* sempre que uma rota crítica (como `/.env`) for acessada.
* **Trigger**: Captura do padrão `DETAIL:GET /.env` via Regex.
* **Payload**: Envio de Timestamp, IP e Porta para um canal de SOC privado.

### 3. SIEM e Análise de Dados (Elasticsearch / Splunk)
O formato unificado de log foi projetado para ser facilmente processado por ferramentas de análise de dados profissionais.
* **Ingestão**: Uso de Filebeat para leitura contínua do arquivo `system_audit.log`.
* **Visualização**: Criação de mapas de calor (Heatmaps) para identificar a origem geográfica dos ataques.

### 4. Inteligência de Ameaças (OSINT)
Integração com APIs como VirusTotal ou AbuseIPDB para enriquecer os dados capturados.
* **Enriquecimento**: O monitor consulta o "Score de Abuso" do IP capturado.
* **Contexto**: Identificação instantânea se o invasor é um bot conhecido ou um novo host infectado.

---

