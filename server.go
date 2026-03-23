package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// logEvent escreve no log e força o Sync no disco para o Windows liberar a leitura
func logEvent(port, proto, ip, detail string) {
	msg := fmt.Sprintf("[%s] PORT:%s | PROTO:%s | IP:%s | DETAIL:%s\n", 
		time.Now().Format("15:04:05"), port, proto, ip, detail)
	
	f, err := os.OpenFile("system_audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(msg)
		f.Sync() // Essencial para o monitor Python ler em tempo real no Windows
	}
}

func startTCPLayer(port string, label string, banner string, wg *sync.WaitGroup) {
	defer wg.Done()
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		remoteIP := conn.RemoteAddr().String()
		go logEvent(port, label, remoteIP, "Tentativa de Conexão")
		
		if banner != "" {
			conn.Write([]byte(banner + "\n"))
		}
		time.Sleep(5 * time.Second) // Tarpit
		conn.Close()
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	// Camada Web (8080)
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			detail := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			logEvent("8080", "HTTP", r.RemoteAddr, detail)
			
			w.Header().Set("Server", "Apache/2.4.41 (Ubuntu)")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized access to %s recorded.", r.URL.Path)
		})
		http.ListenAndServe(":8080", mux)
	}()

	// Portas de Infraestrutura
	go startTCPLayer("21", "FTP", "220 vsFTPd 3.0.3", &wg)
	go startTCPLayer("22", "SSH", "SSH-2.0-OpenSSH_8.2p1", &wg)
	go startTCPLayer("3306", "SQL", "5.7.33-MySQL", &wg)

	wg.Wait()
}