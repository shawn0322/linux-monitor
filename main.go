package main

import "fmt"
import (
	"./mail"
	"strings"
	"bytes"
	"log"
	"strconv"
	"os/exec"
	"os"
)

type Process struct {
	pid     int
	cpu     float64
	mem     float64
	command string
}

type Memary struct {
	total int64
	used  int64
	free  int64
	shared  int64
	buffCache   int64
	available int64
}

func main() {

	hostname,err := os.Hostname()

	cmd := exec.Command("free", "")
	var out bytes.Buffer
	cmd.Stdout = &out
	err1 := cmd.Run()
	if err1 != nil {
		log.Fatal(err1)
	}
	out.ReadString('\n')
	line, err := out.ReadString('\n')
	if(err != nil){
		println(err)
	}

	var memStr string = "<tr>";
	tokens := strings.Split(line, " ")
	for _, t := range tokens {
		if t != "" && t != "\t" {
			memStr += "<td>"+t+"</td>"
		}
	}
	memStr += "</tr>"

	err, processes := getProcessInfo()

	var str string = ""
	var i = 0;
	for _, p := range (processes) {
		i++
		str += "<tr>" +
			"<td> " + strconv.Itoa(p.pid) + " </td>" +
			"<td>" + strconv.FormatFloat(p.cpu, 'f', -1, 32) + " % </td>" +
			"<td>" + strconv.FormatFloat(p.mem, 'f', -1, 32) + " %</td>" +
			"<td> " + p.command + "</td>" +
			"</tr>"
		if i > 9 {
			break;
		}
	}

	to := "gaoxun@loex.com"
	subject := "【LOEX服务器监控】"

	body := `
			<html>
			<body>
			<H1>系统信息</H1>
			<h5>主机名：`+hostname+`</h5>
			<h5>内存状态：</h5>
			<table border="1" style="width: 80%;">
	          <tr>
	             <th>type</th>
	             <th>total</th>
	             <th>used</th>
	             <th>free</th>
	             <th>shared</th>
	             <th>buffCache</th>
	             <th>available</th>
	          </tr>
	         
	            ` + memStr + `
	          
	        </table>
			<H1>进程信息(TOP 10)</H1>
			<table border="1" style="width: 80%;">
	          <tr>
	             <th>进程号</th>
	             <th>CPU使用率</th>
	             <th>内存使用率</th>
	             <th>命令</th>
	          </tr>
	         
	            ` + str + `
	          
	        </table>
			</body>
			</html>
			`
	fmt.Println("send email")
	err = mail.SendToMail(to, subject, body, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}

}

func getProcessInfo() (error, []*Process) {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]*Process, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break;
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range (tokens) {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		//log.Println(len(ft), ft)
		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		mem, err := strconv.ParseFloat(ft[3], 64);
		if err != nil {
			log.Fatal(err)
		}

		command := ft[10];
		processes = append(processes, &Process{pid, cpu, mem, command})
	}
	return err, processes
}
