package main

import(
	"fmt"
	"net/http"
	"net"
	"os"
	"flag"
	"log"
	"bytes"
	"os/exec"
	"io"
)

func exec_shell(s string) (string){
    cmd := exec.Command("/bin/bash", "-c", s)
    var out bytes.Buffer

    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    return out.String()
}

func ExporterHandler(w http.ResponseWriter, r *http.Request){
	cmd := `uptime|awk '{print $11}'|awk -F"," '{print $1}'`
	io.WriteString(w, exec_shell(cmd))
}
func main(){
	port := flag.String("port", "30083", "Input your exporter port")
	flag.Parse()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var host string
	for _, a := range addrs {
		if hostip, ok := a.(*net.IPNet); ok && !hostip.IP.IsLoopback() {
			if hostip.IP.To4() != nil {
					fmt.Println(hostip.IP.String())
					host = hostip.IP.String()
				}
		}
	}
	http.HandleFunc("/metrics", ExporterHandler)
	url := fmt.Sprintf("%s:%s", host, *port)
	fmt.Println("url=", url)
	err = http.ListenAndServe(url, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
