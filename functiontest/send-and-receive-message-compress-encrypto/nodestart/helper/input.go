package helper

import (
	"bufio"
	"fmt"
	"gatesvr/cluster/node"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/nodestart/logics"
	"os"
	"strings"
)

func HandleConsoleInput(proxy *node.Proxy) {
	reader := bufio.NewReader(os.Stdin)
	printhelper()
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.Split(input, " ")
		command := parts[0]
		switch command {
		case "broadcast":
			if len(parts) < 3 {
				fmt.Println("请输入广播内容，例如: broadcast 活动开始")
				continue
			}
			logics.Broadcast(proxy, parts[1], parts[2])
		case "push":
			{
				if len(parts) < 3 {
					fmt.Println("请输入推送内容，例如: push 123 请准备")
					continue
				}
				logics.Push(proxy, parts[1], parts[2])
			}
		case "bindNode":
			{
				if len(parts) < 2 {
					fmt.Println("请输入待绑定用户id，例如: bindNode 1")
					continue
				}
				logics.BindNode(proxy, parts[1])
			}
		case "NodeContinuePush":
			{
				if len(parts) < 3 {
					fmt.Println("请输入推送内容，例如: push 123 请准备")
					continue
				}
				logics.NodeContinuePush(proxy, parts[1], parts[2])
			}
		case "IsOnline":
			{
				if len(parts) < 2 {
					fmt.Println("请输入待查询用户id，例如: isOnline 1")
					continue
				}
				logics.IsOnline(proxy, parts[1])
			}
		//	这两个关闭必须展示，有设计点的呢
		case "ForceClose":
			{
				if len(parts) < 2 {
					fmt.Println("请输入待关闭用户id，例如: ForceClose 1")
					continue
				}
				logics.ForceClose(proxy, parts[1])
			}
		case "GracefulClose":
			{
				if len(parts) < 2 {
					fmt.Println("请输入待关闭用户id，例如: GracefulClose 1")
					continue
				}
				logics.GracefulClose(proxy, parts[1])
			}
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
}
