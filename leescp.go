package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/leeningli/utils"

	"github.com/Unknwon/goconfig"
)

func LeeScpExecute(appname string) {

	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ips, err := cfg.GetValue(appname, "ip")
	if err != nil {
		fmt.Println(err)
	}
	ip_list := strings.Split(ips, ",")
	fmt.Println("ip_list==", ip_list)

	port, err := cfg.GetValue(appname, "port")
	if err != nil {
		fmt.Println(err)
	}
	port_int, _ := strconv.Atoi(port)

	dpath, _ := cfg.GetValue(appname, "dpath")
	spath, _ := cfg.GetValue(appname, "spath")

	scp_flag := 1
	cmd_flag := 1

	if dpath == "" || spath == "" {
		scp_flag = 0
	}

	user, _ := cfg.GetValue(appname, "user")
	pwd, _ := cfg.GetValue(appname, "pwd")
	cmd, _ := cfg.GetValue(appname, "cmd")

	if cmd == "" {
		cmd_flag = 0
	}

	for _, ip := range ip_list {
		if scp_flag == 1 {
			File, err := os.Open(spath)
			if err != nil {
				fmt.Println("open file failed:", err)
				os.Exit(1)
			}
			info, _ := File.Stat()
			defer File.Close()
			fmt.Printf("\r\n---------------%s--------------\r\n", ip)
			scp(user, pwd, ip, port_int, File, info.Size(), dpath)
			fmt.Printf("\r\n-------------------------------\r\n")
		}

		if cmd_flag == 1 {
			fmt.Printf("\r\n--------------%s---------------\r\n", ip)
			utils.RemoteExec(user, pwd, ip, cmd, port_int)
			fmt.Printf("\r\n-------------------------------\r\n")
		}
	}
}

func main() {
	LeeScpExecute("test")
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
	}()

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
