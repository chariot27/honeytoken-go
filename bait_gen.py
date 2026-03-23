import os

def create_baits():
    bait_path = "C:\\Users\\Public\\Documents\\Confidential"
    if not os.path.exists(bait_path):
        os.makedirs(bait_path)

    # Isca 1: TXT com credenciais fakes apontando para o seu IP
    with open(f"{bait_path}\\db_access.txt", "w") as f:
        f.write("--- INTERNAL DATABASE ACCESS ---\n")
        f.write("Admin Panel: http://localhost:8080/admin\n")
        f.write("User: root_maintenance\n")
        f.write("Pass: TempPass2026!\n")

    # Isca 2: .env fake em pastas de projeto
    with open(".env", "w") as f:
        f.write("DB_HOST=localhost\n")
        f.write("DB_PORT=3306\n") # Isso aciona sua porta 3306 em Go
        f.write("AWS_SECRET_LOG=http://localhost:8080/.aws/credentials\n")

    print(f"🪤 Iscas criadas em {bait_path}")

if __name__ == "__main__":
    create_baits()