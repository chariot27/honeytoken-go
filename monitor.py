import time, os, re
from rich.console import Console
from rich.table import Table
from rich.live import Live
from rich.panel import Panel

console = Console()

def parse_log_line(line):
    # Regex robusta: aceita espaços extras entre os pipes |
    pattern = r"\[(.*?)\]\s*PORT:(\d+)\s*\|\s*PROTO:(.*?)\s*\|\s*IP:(.*?)\s*\|\s*DETAIL:(.*)"
    match = re.search(pattern, line)
    return match.groups() if match else None

def generate_dashboard():
    table = Table(title="[bold red]GHOST LISTENER v4.0 - SOC DASHBOARD[/bold red]", expand=True)
    table.add_column("Horário", style="cyan", justify="center")
    table.add_column("Porta", style="bold white", justify="center")
    table.add_column("Serviço", justify="center")
    table.add_column("IP Atacante", style="bold red")
    table.add_column("Detalhes/Payload", style="green")

    log_path = "system_audit.log"
    
    if os.path.exists(log_path):
        try:
            # Abre, lê e fecha rápido para não bloquear o arquivo
            with open(log_path, "r", encoding="utf-8", errors="ignore") as f:
                lines = f.readlines()
                for line in lines[-12:]:
                    data = parse_log_line(line)
                    if data:
                        h, p, s, ip, d = data
                        # Limpa IPv6/Portas efêmeras do IP (ex: [::1]:54321 -> ::1)
                        clean_ip = ip.replace("[", "").replace("]", "").split(":")[0]
                        if not clean_ip or clean_ip == " ": clean_ip = "localhost"
                        
                        color = "red" if p in ["22", "3306"] else "yellow"
                        table.add_row(h, f"[{color}]{p}[/{color}]", s, clean_ip, d.strip())
        except Exception:
            pass
    
    if table.row_count == 0:
        table.add_row("--:--", "!!", "LISTENING", "0.0.0.0", "[dim]Aguardando eventos...[/dim]")

    return Panel(table, border_style="bright_blue", title="[ MONITORAMENTO ATIVO ]", subtitle="Foco: Cybersecurity Research")

if __name__ == "__main__":
    try:
        with Live(generate_dashboard(), refresh_per_second=5, screen=True) as live:
            while True:
                time.sleep(0.1)
                live.update(generate_dashboard())
    except KeyboardInterrupt:
        console.print("\n[bold red]Encerrado.[/bold red]")