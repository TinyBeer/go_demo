package main

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

var (
	id   int
	name string
	pwd  string
	err  error
)

const (
	MainPage     = "MainPage"
	LoginPage    = "LoginPage"
	RegisterPage = "RegisterPage"
	UserPage     = "UserPage"
	ExitPage     = "ExitPage"
)

func main() {
	app := tview.NewApplication()

	addePages(app)

	if err := app.Run(); err != nil {
		panic(err)
	}

}

func addMainPage(pages *tview.Pages) {
	pages.AddPage(
		MainPage,
		tview.NewList().
			AddItem("登录", "使用账号密码进行登录", '1', func() {
				// 跳转到登录页
				pages.SwitchToPage(LoginPage)
			}).
			AddItem("注册", "注册新的账号", '2', func() {}).
			AddItem("退出", "退出系统", 'q', func() {
				pages.SwitchToPage(ExitPage)
			}),
		false,
		true,
	)
}

func addLoginPage(pages *tview.Pages) {
	pages.AddPage(
		LoginPage,
		tview.NewForm().
			AddInputField("账号", "", 20,
				func(textToCheck string, lastChar rune) bool {
					_, err = strconv.Atoi(textToCheck)
					return err == nil
				},
				func(text string) {
					id, _ = strconv.Atoi(text)
				}).
			AddPasswordField("密码", "", 10, '*', func(text string) {
				pwd = text
			}).
			AddButton("登录", func() {
				fmt.Println("恭喜用户", id, "登录成功")
				fmt.Println("UserId:", id)
				fmt.Println("UserPwd:", pwd)
			}).
			AddButton("退出", func() {
				pages.SwitchToPage(MainPage)
			}),
		false,
		false,
	)
}

func addRegisterPage(pages *tview.Pages) {
	pages.AddPage(
		RegisterPage,
		tview.NewForm().
			AddDropDown("性别", []string{"男", "女"}, 0, nil).
			AddInputField("账号", "", 20,
				func(textToCheck string, lastChar rune) bool {
					_, err = strconv.Atoi(textToCheck)
					return err == nil
				},
				func(text string) {
					id, _ = strconv.Atoi(text)
				}).
			AddInputField("昵称", "", 20, nil, func(text string) {
				name = text
			}).
			AddPasswordField("密码", "", 10, '*', func(text string) {
				pwd = text
			}).
			AddButton("注册", func() {
				fmt.Println("注册成功")
				fmt.Println("UserId:", id)
				fmt.Println("UserName:", name)
				fmt.Println("UserPwd:", pwd)
				pages.SwitchToPage(MainPage)
			}).
			AddButton("退出", func() {
				pages.SwitchToPage(MainPage)
			}),
		false,
		false,
	)
}

func addUserPage(pages *tview.Pages) {
	pages.AddPage(
		UserPage,
		tview.NewList().
			AddItem("在线用户列表", "查看在线用户列表", '1', func() {
				fmt.Println("在线用户列表")
			}).
			AddItem("发送消息", "发送群聊消息", '2', func() {
				fmt.Println("发送群聊消息")
			}).
			AddItem("注销", "注销登录", 'q', func() {
				pages.SwitchToPage(MainPage)
			}).
			SetBorder(true),
		false,
		false,
	)
}

func addExitPage(pages *tview.Pages, exitFunc func()) {
	pages.AddPage(
		ExitPage,
		tview.NewModal().
			SetText("您是否确定要退出？").
			AddButtons([]string{"是", "否"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "否" {
					pages.SwitchToPage(MainPage)
				} else {
					exitFunc()
				}
			}),
		false,
		false,
	)
}

func addePages(app *tview.Application) {
	pages := tview.NewPages()

	// 添加主界面
	addMainPage(pages)
	// 添加登录页面
	addLoginPage(pages)
	// 添加注册页面
	addRegisterPage(pages)
	// 添加用户页面
	addUserPage(pages)
	// 添加退出页面
	addExitPage(pages, app.Stop)

	app.SetRoot(pages, true).SetFocus(pages)
}
