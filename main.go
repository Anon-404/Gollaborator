package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"regexp"
	"time"
)

const (
	Reset   = "\033[0m"

	// Bright Colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Styles
	Bold      = "\033[1m"
)

var requestCount int

func printBanner() {
	fmt.Println(Bold + BrightWhite + "____________________________________________________________________________" + Reset)
	fmt.Println(Bold + BrightCyan + `
   .aMMMMP .aMMMb  dMP     dMP     .aMMMb  dMMMMb 
  dMP"    dMP"dMP dMP     dMP     dMP"dMP dMP"dMP 
 dMP MMP"dMP dMP dMP     dMP     dMMMMMP dMMMMK"  
dMP.dMP dMP.aMP dMP     dMP     dMP dMP dMP.aMF   
VMMMP"  VMMMP" dMMMMMP dMMMMMP dMP dMP dMMMMP"    
                                                  
   .aMMMb  dMMMMb  .aMMMb dMMMMMMP .aMMMb  dMMMMb 
  dMP"dMP dMP.dMP dMP"dMP   dMP   dMP"dMP dMP.dMP 
 dMP dMP dMMMMK" dMMMMMP   dMP   dMP dMP dMMMMK"  
dMP.aMP dMP"AMF dMP dMP   dMP   dMP.aMP dMP"AMF   
VMMMP" dMP dMP dMP dMP   dMP    VMMMP" dMP dMP    ` + Reset)
	fmt.Println(Bold + BrightGreen + "                        William Steven (Anon404)" + Reset)
	fmt.Println(Bold + BrightWhite + "____________________________________________________________________________" + Reset)
}

func startTunnel(port string) {
	cmd := exec.Command("cloudflared", "tunnel", "--url", "http://localhost:"+port)
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	err := cmd.Start()
	if err != nil {
		fmt.Println(Bold + BrightRed + "❌ Failed to start Cloudflared:", err, Reset)
		return
	}

	regex := regexp.MustCompile(`https://[a-zA-Z0-9\-\.]+\.com`)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			match := regex.FindString(line)
			if match != "" && match != "https://www.cloudflare.com" {
				if match == "https://api.trycloudflare.com" {
					fmt.Println(Bold + BrightRed + "❌ Seems like you're offline\nAborting.........." + Reset)
					os.Exit(1)
				} else {
					fmt.Println(Bold + BrightGreen + "🌍 Public URL  :" + Bold + BrightCyan + " " + match + Reset)
				}
			}
		}
	}()

	fmt.Println(Bold + BrightBlue + "🌐 Cloudflare tunnel running..." + Reset)
	fmt.Println(Bold + BrightGreen + "🔒 Private URL :" + Bold + BrightCyan + " http://localhost:" + port + Reset)
}

func checkAndInstallCloudflared() {
	_, err := exec.LookPath("cloudflared")
	if err == nil {
		fmt.Println(Bold + BrightGreen + "✔️  Cloudflared found" + Reset)
		return
	}

	fmt.Println(Bold + BrightYellow + "⚠️   Cloudflared not found. Installing..." + Reset)

	// Try installing with package managers
	if _, err := exec.LookPath("pkg"); err == nil {
		// Termux
		cmd := exec.Command("pkg", "install", "-y", "cloudflared")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else if _, err := exec.LookPath("pacman"); err == nil {
		// Arch/Artix
		cmd := exec.Command("sudo", "pacman", "-Sy", "--noconfirm", "cloudflared")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else if _, err := exec.LookPath("apt"); err == nil {
		// Debian/Ubuntu
		cmd := exec.Command("sudo", "apt", "install", "-y", "cloudflared")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		fmt.Println(Bold + BrightRed + "❌ Could not detect package manager. Please install cloudflared manually." + Reset)
		os.Exit(1)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	requestCount++
	id := requestCount

	html := `
<!DOCTYPE html>
<html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Collaborator</title>
</head>
<body>
  <h1>Collaborator</h1>
  <p>Connection Logged.</p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	remoteIp := r.Header.Get("X-Forwarded-For")
	if remoteIp == "" {
		remoteIp = r.RemoteAddr
	}
	fmt.Println(Bold + BrightWhite + "┌───────────────────────────────────────────────────────────────┐" + Reset)
	fmt.Printf(Bold+BrightGreen+"│ 🔢 Request    │ "+Bold+BrightCyan+"%d\n"+Reset, id)
	fmt.Printf(Bold+BrightGreen+"│ 📅 Time       │ "+Bold+BrightCyan+"%s\n"+Reset, timestamp)
	fmt.Printf(Bold+BrightGreen+"│ 🌐 Host       │ "+Bold+BrightCyan+"%s\n"+Reset, r.Host)
	fmt.Printf(Bold+BrightGreen+"│ 📥 Method     │ "+Bold+BrightCyan+"%s\n"+Reset, r.Method)
	fmt.Printf(Bold+BrightGreen+"│ 🌍 Remote IP  │ "+Bold+BrightCyan+"%s\n"+Reset, remoteIp)
	fmt.Printf(Bold+BrightGreen+"│ 📂 Path       │ "+Bold+BrightCyan+"%s\n"+Reset, r.URL.Path)
	fmt.Printf(Bold+BrightGreen+"│ ❓ Query      │ "+Bold+BrightCyan+"%s\n"+Reset, r.URL.RawQuery)
	fmt.Println(Bold + BrightWhite + "├───────────────────────────────────────────────────────────────┤" + Reset)
	fmt.Println(Bold + BrightCyan + "│ 📦 Request Dump                                               │" + Reset)
	fmt.Println(Bold + BrightWhite + "├───────────────────────────────────────────────────────────────┤" + Reset)

	rawRequest, err := httputil.DumpRequest(r, true)
	if err == nil {
		fmt.Println(Bold + BrightGreen + string(rawRequest) + Reset)
	}

	fmt.Println(Bold + BrightWhite + "└───────────────────────────────────────────────────────────────┘" + Reset)
}

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	printBanner()
	checkAndInstallCloudflared()

	port := "4444"
	fmt.Println(Bold + BrightYellow + "🚀 Starting Collaborator server..." + Reset)
	startTunnel(port)

	http.HandleFunc("/", handle)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(Bold + BrightRed + "❌ Server failed to start:", err, Reset)
	}
}
