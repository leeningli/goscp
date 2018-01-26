package leescp

import(
	"strings"
    "fmt"
    "leeconfig"
    "strconv"
    "flag"
    "os"
    "utils"
	"io"
	"path/filepath"
	"github.com/pkg/sftp"
	"log"
	"path"
) 

func LeeSftp(){
	fmt.Println("----------------------------")
    TOPIC := leeconfig.GetConfig("kcxpservice")
    ips := TOPIC["ip"]
    ip_list := strings.Split(ips,",")
    fmt.Println("ip==", ips)
    fmt.Println("ip_list==", ip_list)
    port_int, err := strconv.Atoi(TOPIC["port"])
    
    if TOPIC["dpath"] == "" ||  TOPIC["spath"] == ""{
    	flag.PrintDefaults()
    	os.Exit(1)
    }
    srcFile, err1 := os.Open(TOPIC["spath"])
    if err1 != nil {
    	fmt.Println("open file failed:", err1)
    	os.Exit(1)
    }
    defer srcFile.Close()

    var (
    	sftpClient *sftp.Client

    )
    var filename = path.Base(TOPIC["spath"])


    if err == nil {
    	for _, ip := range ip_list {
			sftpClient, err = utils.SftpConnect(TOPIC["user"], TOPIC["pwd"], ip, port_int)
			if err != nil {
				log.Fatal(err)
			}
			defer sftpClient.Close()
			dstFile, err := sftpClient.Create(path.Join(TOPIC["dpath"], filename))
			if err != nil {
				log.Fatal(err)	
			}
			defer dstFile.Close()
			buf := make([]byte, 1024)
			for {
				n, _ := srcFile.Read(buf)
				if n == 0 {
					break
				}
				dstFile.Write(buf)
			}
			fmt.Printf("copy file to remote:%s is finish success.\n", ip)
			//scp(TOPIC["user"], TOPIC["pwd"], ip, port_int, File, info.Size(), TOPIC["dpath"])
			utils.RemoteExec(TOPIC["user"], TOPIC["pwd"], ip, TOPIC["cmd"], port_int)
    	}
    }
}

func LeeScpExecute(appname string){
    TOPIC := leeconfig.GetConfig(appname)
    ips := TOPIC["ip"]
    ip_list := strings.Split(ips,",")
    fmt.Println("ip_list==", ip_list)
    port_int, _ := strconv.Atoi(TOPIC["port"])
    
    if TOPIC["dpath"] == "" ||  TOPIC["spath"] == ""{
        flag.PrintDefaults()
        os.Exit(1)
    }
    
    for _, ip := range ip_list {
		File, err := os.Open(TOPIC["spath"])
	    if err != nil {
	        fmt.Println("open file failed:", err)
	        os.Exit(1)
	    }
	    info, _ := File.Stat()
	    defer File.Close()
        scp(TOPIC["user"], TOPIC["pwd"], ip, port_int, File, info.Size(), TOPIC["dpath"])
        utils.RemoteExec(TOPIC["user"], TOPIC["pwd"], ip, TOPIC["cmd"], port_int)
    }
}

func LeeScp(appname string){
    TOPIC := leeconfig.GetConfig(appname)
    ips := TOPIC["ip"]
    ip_list := strings.Split(ips,",")
    fmt.Println("ip_list==", ip_list)
    port_int, _ := strconv.Atoi(TOPIC["port"])
    
    if TOPIC["dpath"] == "" ||  TOPIC["spath"] == ""{
    	flag.PrintDefaults()
    	os.Exit(1)
    }
    	for _, ip := range ip_list {
		File, err := os.Open(TOPIC["spath"])
		    if err != nil {
		        fmt.Println("open file failed:", err)
		        os.Exit(1)
		    }
		    info, _ := File.Stat()
		    defer File.Close()
		scp(TOPIC["user"], TOPIC["pwd"], ip, port_int, File, info.Size(), TOPIC["dpath"])
    	}
}
    
func scp(user, pwd, ip string, port int, File io.Reader, size int64, path string) {
	fmt.Println("path==", path)
	filename := filepath.Base(path)
	dirname := strings.Replace(filepath.Dir(path), "\\", "/", -1)
	fmt.Println("filename==", filename)
	fmt.Println("dirname==", dirname)

	session, err := utils.Connect(user, pwd, ip, port)
	if err != nil {
		fmt.Println("create session is failed:", err)
		return
	}
	go func() {
		w, _ := session.StdinPipe()
		fmt.Fprintln(w, "C0644", size, filename)
		io.CopyN(w, File, size)
		fmt.Fprint(w, "\x00")
		w.Close()
	} ()

	if err := session.Run(fmt.Sprintf("/usr/bin/scp -qrt %s", dirname)); err != nil {
		fmt.Println("execute scp is failed:", err)
		return
	} else {
		fmt.Printf("send to %s is success.\n", ip)
		session.Close()
	}

	buf, err := session.Output(fmt.Sprintf("/usr/bin/md5sum %s", path))
	if err != nil {
			fmt.Println("check md5 is failed:", err)
			return
	}
	fmt.Printf("%s md5 is:\n%s\n", ip, string(buf))
}
