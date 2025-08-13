package helper

import (
	"bufio"
	"fmt"
	"gatesvr/cluster/node"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/gamenode/logics"
	"gatesvr/log"
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
				fmt.Println("请输入广播内容，例如: broadcast 123 活动开始")
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
					fmt.Println("请输入待绑定用户id，例如: bindNode 123")
					continue
				}
				logics.BindNode(proxy, parts[1])
			}
		case "isOnline":
			{
				if len(parts) < 2 {
					fmt.Println("请输入待查询用户id，例如: isOnline 1")
					continue
				}
				online, err := logics.IsOnline(proxy, parts[1])
				if err != nil {
					log.Debugf("uid %d Not Online Or Not Exist", parts[1])
					continue
				}
				if online {
					fmt.Println("用户在线")
				} else {
					fmt.Println("用户离线")
				}

			}
		case "onlineUserCount":
			{
				count, err := logics.OnlineUserCount(proxy)
				if err != nil {
					log.Debugf("get online user count failed, err:%v", err)
					return
				} else {
					log.Debugf("在线用户数:", count)
				}
			}
		case "reconnectDemo":
			{
				if len(parts) < 3 {
					fmt.Println("请输入推送内容，例如: ReconnectDemo 111 请准备")
					continue
				}
				logics.ReconnectDemo(proxy, parts[1], parts[2])
			}
		case "help":
			{
				printhelper()
			}
		case "exit":
			return

		default:
			fmt.Println("无效命令，请重新输入")
		}
	}
}

func printhelper() {
	fmt.Println("==================================================")
	fmt.Println("欢迎使用节点控制台测试工具")
	fmt.Println("请输入命令，例如: broadcast 123 活动开始")
	fmt.Println("==================================================")
	fmt.Println("\n[功能命令列表]")
	fmt.Println("  broadcast <seq> <message>: 广播消息（例如: broadcast 324 活动开始）")
	fmt.Println("  push <uid> <message>: 推送消息给指定用户（例如: push 123 请准备）")
	fmt.Println("  bindNode <uid>: 绑定用户到当前节点（例如: bindNode 123）")
	fmt.Println("  isOnline <uid>: 查询用户是否在线（例如: isOnline 1）")
	fmt.Println("  onlineUserCount: 查询当前在线用户数")
	fmt.Println("  reconnectDemo <uid> <message>: 重连演示（例如: reconnectDemo 111 请准备）")
	fmt.Println("  help: 显示帮助信息")
	fmt.Println("  exit: 退出程序")
	fmt.Println("==================================================")
}
