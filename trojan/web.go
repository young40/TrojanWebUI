package trojan

import (
	"crypto/sha256"
	"fmt"
	"trojan/core"
	"trojan/util"
)

// ShowWebPath 显示Web路径
func ShowWebPath() {
	path, _ := core.GetValue("web_path")
	if path == "" {
		path = "默认路径(未设置)"
	}
	fmt.Println("当前Web路径:", util.Green(path))
}

// SetWebPath 设置Web路径
func SetWebPath() {
	currentPath, _ := core.GetValue("web_path")
	if currentPath != "" {
		fmt.Println("当前Web路径:", util.Green(currentPath))
	}
	newPath := util.Input("请输入新的Web路径: ", "")
	if newPath == "" {
		fmt.Println("已取消修改")
		return
	}
	err := core.SetValue("web_path", newPath)
	if err != nil {
		fmt.Println(util.Red("设置失败:"), err)
	} else {
		fmt.Println(util.Green("Web路径已更新"))
	}
}

// ShowRealm 显示Realm值
func ShowRealm() {
	realm, _ := core.GetValue("realm")
	if realm == "" {
	realm = "test"
	}
	fmt.Printf("当前Realm: %s\n", util.Green(realm))
}

// SetRealm 设置Realm值
func SetRealm() {
	realm, _ := core.GetValue("realm")
	if realm == "" {
	realm = "test"
	}
	newRealm := util.Input(fmt.Sprintf("请输入新的Realm值 [%s]: ", realm), realm)
	if newRealm == "" {
		fmt.Println("未作更改!")
		return
	}
	err := core.SetValue("realm", newRealm)
	if err != nil {
		fmt.Println(util.Red("设置失败:"), err)
	} else {
		fmt.Println(util.Green("Realm已更新"))
	}
}

// WebMenu web管理菜单
func WebMenu() {
	fmt.Println()
	menu := []string{
		"重置web管理员密码",
		"修改显示的域名(非申请证书)",
		"显示Web路径",
		"修改Web路径",
		"显示Realm值",
		"修改Realm值",
	}
	switch util.LoopInput("请选择: ", menu, true) {
	case 1:
		ResetAdminPass()
	case 2:
		SetDomain("")
	case 3:
		ShowWebPath()
	case 4:
		SetWebPath()
	case 5:
		ShowRealm()
	case 6:
		SetRealm()
	}
}

// ResetAdminPass 重置管理员密码
func ResetAdminPass() {
	inputPass := util.Input("请输入admin用户密码: ", "")
	if inputPass == "" {
		fmt.Println("撤销更改!")
	} else {
		encryPass := sha256.Sum224([]byte(inputPass))
		err := core.SetValue("admin_pass", fmt.Sprintf("%x", encryPass))
		if err == nil {
			fmt.Println(util.Green("重置admin密码成功!"))
		} else {
			fmt.Println(err)
		}
	}
}

// SetDomain 设置显示的域名
func SetDomain(domain string) {
	if domain == "" {
		domain = util.Input("请输入要显示的域名地址: ", "")
	}
	if domain == "" {
		fmt.Println("撤销更改!")
	} else {
		core.WriteDomain(domain)
		Restart()
		fmt.Println("修改domain成功!")
	}
}

// GetDomainAndPort 获取域名和端口
func GetDomainAndPort() (string, int) {
	config := core.GetConfig()
	return config.SSl.Sni, config.LocalPort
}
