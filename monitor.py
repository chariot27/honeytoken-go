import time, os, re
from rich.console import Console
from rich.table import Table
from rich.live import Live
from rich.panel import Panel

console = Console()

def parse_log_line(line):
    # Regex atualizada para capturar o novo campo de detalhes (incluindo User-Agent)
    pattern = r"\[(.*?)\]\s*PORT:(\d+)\s*\|\s*PROTO:(.*?)\s*\|\s*IP:(.*?)\s*\|\s*DETAIL:(.*)"
    match = re.search(pattern, line)
    return match.groups() if match else None

def get_last_lines(filepath, n=15):
    """Lê as últimas N linhas de forma eficiente."""
    if not os.path.exists(filepath):
        return []
    try:
        with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
            # Para arquivos pequenos, readlines é ok, mas pegamos apenas o final
            return f.readlines()[-n:]
    except Exception:
        return []

def generate_dashboard():
    table = Table(title="[bold red]GHOST LISTENER v4.1 - SOC DASHBOARD[/bold red]", expand=True, border_style="blue")
    
    table.add_column("Horário", style="cyan", justify="center", width=10)
    table.add_column("Porta", justify="center", width=8)
    table.add_column("Serviço", justify="center", width=10)
    table.add_column("IP Atacante", style="bold yellow", width=16)
    table.add_column("Payload / User-Agent", style="green", no_wrap=True)

    lines = get_last_lines("system_audit.log")
    
    for line in lines:
        data = parse_log_line(line)
        if data:
            h, p, s, ip, d = data
            
            # Lógica de Alerta Visual
            row_style = ""
            # Se for porta crítica ou contiver tentativa de path traversal/config
            if p in ["22", "3306"] or any(x in d.lower() for x in [".env", "admin", "config", "../"]):
                row_style = "bold red"
                p_display = f"[blink red]{p}[/blink red]"
            else:
                p_display = p

            table.add_row(
                h, 
                p_display, 
                s.strip(), 
                ip.strip(), 
                d.strip(),
                style=row_style
            )
    
    if table.row_count == 0:
        table.add_row("--:--", "!!", "WAITING", "0.0.0.0", "[dim]Aguardando telemetria...[/dim]")

    return Panel(
        table, 
        border_style="bright_blue", 
        title="[bold white] MONITORAMENTO DE ATIVOS [/bold white]", 
        subtitle="[bold yellow]Labs: Security Research[/bold yellow]"
    )

if __name__ == "__main__":
    try:
        # refresh_per_second=10 para uma sensação mais 'real-time'
        with Live(generate_dashboard(), refresh_per_second=12, screen=True) as live:
                while True:
                    time.sleep(0.05) # Loop interno mais rápido
                    live.update(generate_dashboard())
    except KeyboardInterrupt:
        console.print("\n[bold red]SOC Dashboard encerrado com segurança.[/bold red]")