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
				fmt.Println("请输入openid，例如: login 123")
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
				fmt.Printf("成功发送登录请求，cid：%d,uid：%d\n", conn.ID(), conn.UID())
			}
		case "push":
			if len(parts) < 2 {
				fmt.Println("请输入要发送的信息，例如: push hello")
				continue
			}
			logics.PushMessage(conn, parts[1])
		case "forceclose":
			logics.ForceCloseConn(conn)
		case "grececlose":
			logics.GreceCloseConn(conn)
		case "testLimiter":
			logics.TestLimiter(conn)
		case "testLimiterCriticalMessage":
			logics.TestLimiterCriticalMessage(conn)
		case "testStatefulRoute":
			logics.TestStaefulRoute(conn)
		case "retryConnect":
			logics.RetryConnect(conn, proxy)
		case "testBenchMark":
			benchtest.TestBenchMark(proxy)
		case "help":
			printhelper()
		case "exit":
			return
		default:
			fmt.Println("无效命令，请重新输入")
		}

		//case "ClientContinuePush":
		//	logics.ClientContinuePush(conn)
		//case "TestConnectionLimit":
		//	logics.TestConnectionLimit(proxy)
		//case "testForwardMessage":
		//	logics.TestForwardMessage(conn)

	}
}

func printhelper() {
	fmt.Println("==================================================")
	fmt.Println("欢迎使用客户端功能测试工具")
	fmt.Println("首次运行请先执行 login <uid> 命令连接网关服务器并鉴权")
	fmt.Println("==================================================")
	fmt.Println("\n请输入命令，例如: login 123")
	fmt.Println("\n[测试命令列表]")
	fmt.Println("  login <uid>: 连接服务器并鉴权（例如: login 123）")
	fmt.Println("  push <message>: 发送消息（例如: push hello）")
	fmt.Println("  forceclose: 强制关闭连接")
	fmt.Println("  grececlose: 优雅关闭连接")
	fmt.Println("  testLimiter: 测试限流功能")
	fmt.Println("  testLimiterCriticalMessage: 测试限流关键消息")
	fmt.Println("  testStatefulRoute: 测试有状态路由")
	fmt.Println("  retryConnect: 重连服务器")
	fmt.Println("  testBenchMark: 性能基准测试")
	fmt.Println("  help: 显示帮助信息")
	fmt.Println("  exit: 退出程序")
	fmt.Println("==================================================")
}
