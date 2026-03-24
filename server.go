package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// Global logger channel para evitar concorrência no arquivo
var logChan = make(chan string, 100)

func loggerWorker() {
	f, err := os.OpenFile("system_audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Erro ao abrir log: %v\n", err)
		return
	}
	defer f.Close()

	for msg := range logChan {
		f.WriteString(msg)
		f.Sync() // Mantém o monitor Python atualizado em tempo real
	}
}

func logEvent(port, proto, remoteAddr, detail string) {
	// Limpa o IP (remove a porta do atacante e colchetes de IPv6)
	ip, _, _ := net.SplitHostPort(remoteAddr)

	msg := fmt.Sprintf("[%s] PORT:%s | PROTO:%s | IP:%s | DETAIL:%s\n",
		time.Now().Format("15:04:05"), port, proto, ip, detail)

	// Envia para o worker de log de forma não-bloqueante
	select {
	case logChan <- msg:
	default:
		// Se o canal estiver cheio, descarta ou imprime no console
	}
}

func startTCPLayer(port string, label string, banner string, wg *sync.WaitGroup) {
	defer wg.Done()
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Erro ao abrir porta %s: %v\n", port, err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		// A mágica acontece aqui: Iniciamos uma Goroutine para CADA conexão.
		// Isso libera o loop 'for' para aceitar o próximo ataque imediatamente.
		go func(c net.Conn) {
			defer c.Close()

			// Loga no exato milissegundo em que o Handshake ocorre
			logEvent(port, label, c.RemoteAddr().String(), "TCP Handshake Detectado")

			if banner != "" {
				// Pequeno delay apenas para parecer um servidor real processando
				time.Sleep(200 * time.Millisecond)
				c.Write([]byte(banner + "\r\n"))
			}

			// Tarpit: Segura o atacante por 10 segundos sem travar os outros
			time.Sleep(10 * time.Second)

			// O logEvent de encerramento é opcional, mas ajuda a ver o fim da conexão
			// logEvent(port, label, c.RemoteAddr().String(), "Conexão encerrada pelo Tarpit")
		}(conn)
	}
}

func main() {
	// Inicia o worker de log em background
	go loggerWorker()

	var wg sync.WaitGroup
	wg.Add(4)

	// Camada Web (8080)
	go func() {
		defer wg.Done()
		server := &http.Server{
			Addr: ":8080",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				detail := fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, r.UserAgent())
				logEvent("8080", "HTTP", r.RemoteAddr, detail)

				// Ofuscação de Server Header
				w.Header().Set("Server", "Apache/2.4.41 (Ubuntu)")
				w.Header().Set("X-Powered-By", "PHP/7.4.3")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "<html><head><title>401 Authorization Required</title></head><body><h1>401 Unauthorized</h1><p>Access to %s is restricted.</p></body></html>", r.URL.Path)
			}),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		server.ListenAndServe()
	}()

	// Portas de Infraestrutura
	go startTCPLayer("21", "FTP", "220 vsFTPd 3.0.3", &wg)
	go startTCPLayer("22", "SSH", "SSH-2.0-OpenSSH_8.2p1 Ubuntu-4ubuntu0.1", &wg)
	go startTCPLayer("3306", "SQL", "5.7.33-0ubuntu0.20.04.1", &wg)

	fmt.Println("Honeytoken System v4.1 rodando... Pressione Ctrl+C para encerrar.")
	wg.Wait()
}
