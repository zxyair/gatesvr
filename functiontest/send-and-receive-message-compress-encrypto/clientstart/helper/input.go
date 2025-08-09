package helper

import (
	"bufio"
	"fmt"
	"gatesvr/cluster/client"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/benchtest"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/logics"
	"gatesvr/log"
	"os"
	"strconv"
	"strings"
)

func HandleConsoleInput(proxy *client.Proxy) {
	reader := bufio.NewReader(os.Stdin)
	var conn *client.Conn
	printhelper()
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {
		case "login":
			if len(parts) < 2 {
				fmt.Println("请输入uid，例如: login 123")
				continue
			}
			uid, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				fmt.Println("无效的uid格式，请输入数字")
				continue
			}
			if conn, err = proxy.Dial(); err != nil {
				log.Errorf("connect server failed: %v", err)
				return
			}
			addr, err := conn.LocalAddr()
			if err != nil {
				fmt.Println("获取本地地址失败:", err)
				return
			}
			remoteAddr, err := conn.RemoteAddr()
			if err != nil {
				fmt.Println("获取远程地址失败:", err)
				return
			}
			fmt.Printf("连接网关服务器成功,cid:%d，本机地址为：%v,网关地址为：%v\n", conn.ID(), addr, remoteAddr)
			conn.Bind(uid)
			succ := logics.Authorition(conn, uid)
			if succ {
				fmt.Printf("登录成功，cid：%d,uid：%d\n", conn.ID(), conn.UID())
			}
		case "push":
			logics.PushMessage(conn)
		case "RetryConnect":
			logics.RetryConnect(conn, proxy)
		case "ClientContinuePush":
			logics.ClientContinuePush(conn)
		case "forceclose":
			logics.ForceCloseConn(conn)
		case "grececlose":
			logics.GreceCloseConn(conn)
		case "TestLimiter":
			logics.TestLimiter(conn)
		case "TestLimiterCriticalMessage":
			logics.TestLimiterCriticalMessage(conn)
		//case "TestConnectionLimit":
		//	logics.TestConnectionLimit(proxy)
		case "TestForwardMessage":
			logics.TestForwardMessage(conn)
		case "TestStaefulRoute":
			logics.TestStaefulRoute(conn)
		case "TestBenchMark":
			benchtest.TestBenchMark(proxy)
		case "help":
			printhelper()
		case "exit":
			return
		default:
			fmt.Println("无效命令，请重新输入")
		}
	}
}

func printhelper() {
	fmt.Println("欢迎使用客户端功能测试工具——首次运行请先执行 dial 命令连接网关服务器, 之后输入auth <uid> 进行鉴权")
	fmt.Println("请输入命令，例如: dial")
	fmt.Println("功能测试命令列表:")
	fmt.Println("dial: 连接服务器")
	fmt.Println("auth <uid>: 鉴权（例如: auth 123）")
	fmt.Println("push: 验证通信功能以及路由功能")
	fmt.Println("RetryConnect: 重连服务器")
	fmt.Println("ClientContinuePush: 验证client发送消息的有序性")
	fmt.Println("forceclose: 强制关闭连接")
	fmt.Println("grececlose: 优雅关闭连接")
	fmt.Println("help: 帮助")
	fmt.Println("性能测试命令列表:")
}
