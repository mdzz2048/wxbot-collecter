package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
	"wxbot/cubox"
	"wxbot/database"
	"wxbot/simpread"
	"wxbot/utils"

	"github.com/eatmoreapple/openwechat"
)

var config = utils.Config{}.GetConfig()
var collectType = utils.GetCollectType()

func SaveURL(url string) utils.Article {
	switch collectType {
	case "cubox":
		return SaveToCubox(url)
	case "simpread":
		return SaveToSimpRead(url)
	default:
		log.Fatalln("未定义的保存方式: ", collectType)
		return utils.Article{}
	}
}

func SaveToCubox(url string) utils.Article {
	webInfo, err := cubox.SearchEngineWebInfo(url)
	if err != nil {
		log.Fatalln("获取网页信息失败: ", err)
		return utils.Article{}
	}
	bookmark, err := cubox.SearchEngineNew(url, webInfo)
	if err != nil {
		log.Fatalln("网页收藏失败: ", err)
		return utils.Article{}
	}

	article := database.ConvertBookMarkToArticle(bookmark)
	db, _ := database.ConnectDatabase()
	database.ArticleAdd(db, &article)

	return article
}

func SaveToSimpRead(url string) utils.Article {
	webInfo, err := cubox.SearchEngineWebInfo(url)
	if err != nil {
		log.Fatalln("获取网页信息失败: ", err)
		return utils.Article{}
	}
	urlInfo := simpread.WebInfo{
		URL:         webInfo.URL,
		Title:       webInfo.Title,
		Description: webInfo.Description,
	}
	success := simpread.AddURL(url, &urlInfo)
	if !success {
		log.Fatalln("网页收藏失败: ", err)
		return utils.Article{}
	}

	article := utils.Article{
		Title:       urlInfo.Title,
		URL:         urlInfo.URL,
		Description: urlInfo.Description,
	}
	db, _ := database.ConnectDatabase()
	database.ArticleAdd(db, &article)

	return article
}

func IsOwner(ctx *openwechat.MessageContext) bool {
	sender, _ := ctx.Sender()
	// RemarkName: 备注的名字, NickName: 用户昵称
	fmt.Println("收到消息, 来自: ", sender.NickName)
	return sender.NickName == config.WeChat.WeChatName
}

func main() {
	database.InitDatabase()
	// 桌面模式
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 设置最大请求次数
	bot.Caller.Client.MaxRetryTimes = 10

	// 心跳包回调函数
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		// 当返回的 RetCode 不为0时打印输出。
		if !resp.Success() {
			println("RetCode:%s  Selector:%s", resp.RetCode, resp.Selector)
		}
	}

	// 日志输出到文件: 创建、追加、读写，777，所有权限
	file, err := os.OpenFile("wxbot.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		file.Close()
	}()
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 获取消息发生错误的 handle, 返回 nil 则尝试继续监听
	var err_count int32
	bot.MessageErrorHandler = func(err error) error {
		// 日志记录
		log.Fatalln("Message error: ", err)

		// 错误计数
		atomic.AddInt32(&err_count, 1)
		if err_count > 3 {
			// 直接退出程序
			log.Println("Logout.")
			bot.Logout()
		}
		return nil
	}

	// 免扫码登录
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
	log.Println("Get Start.")

	// 消息处理
	dispatcher := openwechat.NewMessageMatchDispatcher()
	dispatcher.OnMedia(func(ctx *openwechat.MessageContext) {
		if !IsOwner(ctx) {
			return
		}

		mediaData, err := ctx.MediaData()
		if err == nil {
			title := mediaData.AppMsg.Title
			url := mediaData.AppMsg.URL
			fmt.Printf("%v 收到消息: %s\n", ctx.CreateTime, url)

			article := SaveURL(url)
			if article.Title == "" {
				replyText := fmt.Sprintln("文章保存失败!")
				ctx.ReplyText(replyText)
			} else {
				replyText := fmt.Sprintf("文章已保存至 %s: \n标题: %s\n链接: %s\n", collectType, title, url)
				ctx.ReplyText(replyText)
			}
		}
	})
	dispatcher.OnText(func(ctx *openwechat.MessageContext) {
		if !IsOwner(ctx) {
			return
		}

		textData := ctx.Message.Content
		if utils.IsURL(textData) {
			fmt.Printf("%v 收到消息: %s", ctx.CreateTime, textData)

			article := SaveURL(textData)

			if article.Title == "" {
				replyText := fmt.Sprintln("尝试保存链接失败!")
				ctx.ReplyText(replyText)
			} else {
				replyText := fmt.Sprintf("文章已保存至 %s: \n标题: %s\n链接: %s\n", collectType, article.Title, article.URL)
				ctx.ReplyText(replyText)
			}
		}
	})

	bot.MessageHandler = dispatcher.AsMessageHandler()

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
