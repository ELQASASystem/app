package http

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ELQASASystem/app/internal/app"
	"github.com/ELQASASystem/app/internal/app/database"

	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
	"github.com/unidoc/unioffice/document"
)

// StartAPI 开启 API 服务
func StartAPI() {

	app := iris.New()
	API := app.Party("apis/")

	// 握手
	API.Party("hello").Get("/", func(c *context.Context) {
		_, _ = c.JSON(iris.Map{"message": "hello"})
	})

	// 登录
	API.Party("sign").Get("/{u}/in", func(c *context.Context) {

		res, err := class.Database().Account().ReadAccountsList(c.Params().Get("u"))
		if err != nil {
			log.Error().Err(err).Msg("获取用户帐号失败")
			_, _ = c.JSON(iris.Map{"message": "no"})
			return
		}

		if c.URLParam("p") != res.Password {
			_, _ = c.JSON(iris.Map{"message": "no"})
			return
		}

		_, _ = c.JSON(iris.Map{"message": "yes"})
	})

	Group := API.Party("group")
	{

		// 获取群列表
		Group.Get("/list", func(c *context.Context) {

			type groupList struct {
				ID       uint64 `json:"id"`        // 群号
				Name     string `json:"name"`      // 群名
				MemCount uint16 `json:"mem_count"` // 群成员数
			}

			var data []groupList
			for _, v := range class.Bot.C.GroupList {
				data = append(data, groupList{uint64(v.Uin), v.Name, v.MemberCount})
			}

			_, _ = c.JSON(data)
		})

		// 获取群成员
		Group.Get("/{i}/mem", func(c *context.Context) {

			i, err := c.Params().GetInt64("i")
			if err != nil {
				log.Error().Err(err).Msg("解析群号失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			type memList struct {
				ID   uint64 `json:"id"`   // 群员帐号
				Name string `json:"name"` // 群员名片
			}

			var data []memList
			for _, v := range class.Bot.C.FindGroupByUin(i).Members {

				var name string
				if n := v.CardName; n != "" {
					name = n
				} else {
					name = v.Nickname
				}

				data = append(data, memList{uint64(v.Uin), name})

			}

			_, _ = c.JSON(data)
		})

		// 表扬
		Group.Get("/{i}/praise", func(c *context.Context) {

			i, err := c.Params().GetUint64("i")
			if err != nil {
				log.Error().Err(err).Msg("解析目标群失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			var ids []uint64
			err = jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(c.URLParam("mem"), &ids)
			if err != nil {
				log.Error().Err(err).Msg("解析目标成员失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			m := class.Bot.NewMsg().AddText("表扬以下答对的同学:\n")
			for _, id := range ids {
				m.AddAt(id)
			}
			class.Bot.SendGroupMsg(m.AddText("\n希望同学们再接再厉!").To(i))
			_, _ = c.JSON(iris.Map{"message": "yes"})
		})

	}

	Question := API.Party("question")
	{

		// 获取问题列表
		Question.Get("/{u}/list", func(c *context.Context) {

			res, err := class.Database().Question().ReadQuestionList(c.Params().Get("u"))
			if err != nil {
				log.Error().Err(err).Msg("读取问题列表失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			_, _ = c.JSON(res)
		})

		// 获取问题
		Question.Get("/a/{i}", func(c *context.Context) {

			i, err := c.Params().GetUint32("i")
			if err != nil {
				log.Error().Err(err).Msg("解析问题ID失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			res, err := class.ReadQuestion(i)
			if err != nil {
				log.Error().Err(err).Msg("获取答题失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			_, _ = c.JSON(res)
		})

		// 新增问题
		Question.Post("/add", func(c *context.Context) {

			// 初始化问题列表, 用于解析 JSON 后储存
			qlt := database.QuestionListTab{}

			if err := c.ReadJSON(&qlt); err != nil {
				log.Error().Err(err).Msg("解析传入 JSON 失败")
				return
			}

			err := class.Database().Question().WriteQuestionList(&qlt)
			if err != nil {
				log.Error().Err(err).Msg("新增答题失败")
				return
			}

			_, _ = c.JSON(iris.Map{"message": "yes"})
		})

		// 开始问答
		Question.Get("/{question_id}/start", func(c *context.Context) {

			qid, err := c.Params().GetUint32("question_id")
			if err != nil {
				log.Error().Err(err).Msg("解析问题 ID 失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
			}

			if err = class.StartQA(qid); err != nil {
				log.Error().Err(err).Msg("开启问答失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
			}
			_, _ = c.JSON(iris.Map{"message": "yes"})
		})

		// 停止问答
		Question.Get("/{question_id}/stop", func(c *context.Context) {

			qid, err := c.Params().GetUint32("question_id")
			if err != nil {
				log.Error().Err(err).Msg("解析问题 ID 失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
			}

			if err = class.StopQA(qid); err != nil {
				log.Error().Err(err).Msg("停止答题失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}
			_, _ = c.JSON(iris.Map{"message": "yes"})
		})

		// 准备问答
		Question.Get("/{question_id}/prepare", func(c *context.Context) {

			qid, err := c.Params().GetUint32("question_id")
			if err != nil {
				log.Error().Err(err).Msg("解析问题 ID 失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
			}

			if err = class.PrepareQA(qid); err != nil {
				log.Error().Err(err).Msg("准备答题失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}
			_, _ = c.JSON(iris.Map{"message": "yes"})
		})

		// 删除问题
		Question.Get("/{question_id}/delete", func(c *context.Context) {

			_, err := c.Params().GetUint64("question_id")
			if err != nil {
				log.Error().Err(err).Msg("解析问题ID失败")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			// TODO: 调用数据库删除 QJNKSM:这个先咕咕
			// class.Database().Question().RemoveQuestion(qid)

			_, _ = c.JSON(iris.Map{"message": "yes"})
		})

		// 获取问题市场
		Question.Get("/market", func(c *context.Context) {

			res, err := class.Database().Question().ReadQuestionMarket()
			if err != nil {
				log.Error().Err(err).Msg("读取问题列表失败")
				return
			}

			_, _ = c.JSON(res)
		})

	}

	Upload := API.Party("upload")
	{

		// 上传 Docx 前预检请求
		Upload.Options("/docx", func(c *context.Context) {

			c.Header("Access-Control-Allow-Headers", "x-requested-with")
			c.Header("Access-Control-Allow-Methods", "POST")
		})

		// 上传 Docx
		Upload.Post("/docx", func(c *context.Context) {

			c.SetMaxRequestBodySize(10485760) // 限制最大上传大小为 10MiB

			_, fileHeader, err := c.FormFile("file")
			if err != nil {
				log.Error().Err(err).Msg("文件上传失败")
				return
			}

			encodedName := class.HashForSHA1(fileHeader.Filename+strconv.FormatInt(time.Now().Unix(), 10)) + ".docx"
			dest := filepath.Join("web/assets/temp/docx/", encodedName)

			log.Info().Str("文件名", encodedName).Msg("API：上传文件")

			if _, err := c.SaveFormFile(fileHeader, dest); err != nil {
				log.Error().Err(err).Msg("保存上传文件失败")
				return
			}

			_, _ = c.JSON(iris.Map{"fileName": encodedName})

			// 在一分钟后删除该文件
			time.AfterFunc(time.Minute, func() {

				log.Info().Str("文件名", encodedName).Msg("API：删除上传的文件")
				if err := os.Remove(dest); err != nil {
					log.Error().Err(err).Msg("删除文件时发生了意外")
				}

			})
		})

		// 解析 docx 文件
		Upload.Get("/docx/{p}/parse", func(c *context.Context) {

			doc, err := document.Open("web/assets/temp/docx/" + c.Params().Get("p"))
			if err != nil {
				log.Error().Err(err).Msg("打开 Docx 失败")
				return
			}

			var data []string
			for _, v := range doc.Paragraphs() {

				var data0 []string
				for _, vv := range v.Runs() {
					data0 = append(data0, vv.Text())
				}

				data = append(data, strings.Join(data0, ""))

			}

			_, _ = c.JSON(data)
		})

		// 上传图片前预检请求
		Upload.Options("/picture", func(c *context.Context) {

			c.Header("Access-Control-Allow-Headers", "x-requested-with")
			c.Header("Access-Control-Allow-Methods", "POST")
		})

		// 上传图片
		Upload.Post("/picture", func(c *context.Context) {

			c.SetMaxRequestBodySize(4194304) // 限制最大上传大小为 4MiB

			_, fileHeader, err := c.FormFile("file")
			if err != nil {
				log.Error().Err(err).Msg("文件上传失败")
				return
			}

			encodedName := class.HashForSHA1(fileHeader.Filename+strconv.FormatInt(time.Now().Unix(), 10)) + "-" + fileHeader.Filename
			dest := filepath.Join("web/assets/question/pictures/", encodedName)

			log.Info().Str("文件名", fileHeader.Filename).Msg("API：上传文件")

			if _, err := c.SaveFormFile(fileHeader, dest); err != nil {
				log.Error().Err(err).Msg("保存上传文件失败")
				return
			}

			_, _ = c.JSON(iris.Map{"fileName": fileHeader.Filename})
		})

		Upload.Get("/{text}/split", func(c *context.Context) {
			words, err := class.Bot.C.GetWordSegmentation(c.Params().Get("text"))

			if err != nil {
				log.Error().Err(err).Msg("分词时出错")
				_, _ = c.JSON(iris.Map{"message": "no"})
				return
			}

			for k, v := range words {
				words[k] = strings.ReplaceAll(v, "\u0000", "")
			}

			_, _ = c.JSON(words)
		})

	}

	if err := app.Listen(":4040"); err != nil {
		log.Panic().Err(err).Msg("启动 API 服务失败")
	}

}
