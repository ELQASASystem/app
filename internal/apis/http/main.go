package http

import (
	"github.com/ELQASASystem/app/internal/app"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ELQASASystem/app/internal/app/database"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
	"github.com/unidoc/unioffice/document"
)

// StartAPI 开启 API 服务
func StartAPI() {

	app := iris.New()
	API := app.Party("apis/")
	{

		// 握手
		Hello := API.Party("hello")
		{

			Hello.Get("/", func(c *context.Context) {
				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(iris.Map{"message": "hello"})
			})

		}

		Login := API.Party("sign")
		{

			// 登录
			Login.Get("in/{u}/{p}", func(c *context.Context) {

				c.Header("Access-Control-Allow-Origin", "*")

				pa := c.Params()
				res, err := database.Class.Account.ReadAccountsList(pa.Get("u"))
				if err != nil {
					log.Error().Err(err).Msg("校验密码失败")
					_, _ = c.JSON(iris.Map{"message": "no"})
					return
				}

				if pa.Get("p") != res.Password {
					_, _ = c.JSON(iris.Map{"message": "no"})
					return
				}

				_, _ = c.JSON(iris.Map{"message": "yes"})

			})

		}

		Group := API.Party("group")
		{

			// 获取群列表
			Group.Get("/list", func(c *context.Context) {

				err := class.Bot.C.ReloadGroupList()
				if err != nil {
					log.Error().Err(err).Msg("重新载入群列表失败")
					return
				}

				type groupList struct {
					ID       uint64 `json:"id"`        // 群号
					Name     string `json:"name"`      // 群名
					MemCount uint16 `json:"mem_count"` // 群成员数
				}

				var data []groupList
				for _, v := range class.Bot.C.GroupList {
					data = append(data, groupList{uint64(v.Uin), v.Name, v.MemberCount})
				}

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(data)

			})

			// 获取群成员
			Group.Get("/mem/{i}", func(c *context.Context) {

				i, err := c.Params().GetInt64("i")
				if err != nil {
					log.Error().Err(err).Msg("解析群号失败")
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

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(data)

			})

			// 激励
			Group.Get("/praise", func(c *context.Context) {
				target := c.URLParam("target")
				mem := c.URLParam("mem")

				c.Header("Access-Control-Allow-Origin", "*")

				if len(target) <= 0 && len(mem) <= 0 {
					_, _ = c.JSON(iris.Map{"message": "no"})
					return
				} else {
					if gid, err := strconv.ParseUint(mem, 10, 64); err != nil {
						_, _ = c.JSON(iris.Map{"message": "no"})
						return
					} else {
						mem = strings.TrimPrefix(mem, "[")
						mem = strings.TrimSuffix(mem, "]")
						ids := strings.Split(mem, ",")

						if class.Bot.C.FindGroup(int64(gid)) == nil {
							_, _ = c.JSON(iris.Map{"message": "no"})
							return
						} else {
							msgToSend := class.Bot.NewMsg().AddText("表扬以下答对的同学:\n")
							for _, id := range ids {
								if id, err := strconv.ParseInt(id, 10, 64); err == nil {
									msgToSend.AddAt(id)
								}
							}
							class.Bot.SendGroupMsg(msgToSend.AddText("\n希望同学们再接再厉!").To(gid))
							_, _ = c.JSON(iris.Map{"message": "yes"})
							return
						}
					}
				}
			})

		}

		Question := API.Party("question")
		{

			// 获取问题列表
			Question.Get("/list/{u}", func(c *context.Context) {

				res, err := database.Class.Question.ReadQuestionList(c.Params().Get("u"))
				if err != nil {
					log.Error().Err(err).Msg("读取问题列表失败")
					return
				}

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(res)

			})

			// 获取问题
			Question.Get("/a/{i}", func(c *context.Context) {

				i, err := c.Params().GetUint32("i")
				if err != nil {
					log.Error().Err(err).Msg("解析问题ID失败")
					return
				}

				type Question struct {
					*database.QuestionListTab
					Answer []*database.AnswerListTab `json:"answer"`
				}

				res, err := database.Class.Question.ReadQuestion(i)
				if err != nil {
					log.Error().Err(err).Msg("读取问题失败")
					return
				}

				res2, err := database.Class.Answer.ReadAnswerList(i)
				if err != nil {
					log.Error().Err(err).Msg("读取回答失败")
					return
				}

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(Question{res, res2})

			})

			// 新增问题
			Question.Get("/add/{question}/{creator_id}/{market}", func(c *context.Context) {

				pa := c.Params()
				err := database.Class.Question.WriteQuestionList(&database.QuestionListTab{
					Question:  pa.Get("question"),
					CreatorID: pa.Get("creator_id"),
					Market:    pa.GetBoolDefault("market", false),
				})
				if err != nil {
					log.Error().Err(err).Msg("新增答题失败")
					return
				}

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(iris.Map{"message": "yes"})

			})

			// 发布问题
			Question.Get("/publish/{question_id}", func(c *context.Context) {

				pa := c.Params()
				qid, err := pa.GetUint32("question_id")
				if err != nil {
					log.Error().Err(err).Msg("解析问题失败")
				}

				if data, err := database.Class.Question.ReadQuestion(qid); data != nil {
					if err != nil {
						log.Error().Err(err).Msg("读取问题失败")
						return
					}

					class.PublishQuestion(class.MakeQuestion(uint64(data.ID), data.Target, data.Question))

					c.Header("Access-Control-Allow-Origin", "*")
					_, _ = c.JSON(iris.Map{"message": "yes"})
				} else {
					c.Header("Access-Control-Allow-Origin", "*")
					_, _ = c.JSON(iris.Map{"message": "no"})
				}

			})

			// 删除问题
			Question.Get("/delete/{question_id}", func(c *context.Context) {

				pa := c.Params()

				qid, err := pa.GetUint64("question_id")
				if err != nil {
					log.Error().Err(err).Msg("解析问题ID失败")
				}

				// TODO: 调用数据库删除 QJNKSM:这个先咕咕
				class.ExpiredQuestion(qid)

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(iris.Map{"message": "yes"})

			})

			// 获取问题市场
			Question.Get("/market", func(c *context.Context) {

				res, err := database.Class.Question.ReadQuestionMarket()
				if err != nil {
					log.Error().Err(err).Msg("读取问题列表失败")
					return
				}

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(res)

			})

			// 准备 / 终止答题
			Question.Get("/control/{question_id}/{status}", func(c *context.Context) {
				pa := c.Params()

				qid, err := pa.GetUint64("question_id")
				if err != nil {
					log.Error().Err(err).Msg("解析问题ID失败")
					return
				}

				status, err := pa.GetUint("status")

				switch status {
				case 0:
					{
						// 执行准备操作
						_, _ = c.JSON(iris.Map{"message": "yes"})
					}
				case 1:
					{
						// 执行终止操作
						_, ok := class.ExpiredQuestion(qid)
						_, _ = c.JSON(iris.Map{"message": ok})
					}
				}
			})

			// 发布答题
			Question.Get("/post/{question_id}", func(c *context.Context) {
				pa := c.Params()

				qid, err := pa.GetUint64("question_id")
				if err != nil {
					log.Error().Err(err).Msg("解析问题ID失败")
					return
				}

				if q, _, ok := class.GetQuestionByID(qid); ok {
					class.PublishQuestion(q)
				}
			})

		}

		Upload := API.Party("upload")
		{

			// 上传 Docx 前预检请求
			Upload.Options("/docx", func(c *context.Context) {

				c.Header("Access-Control-Allow-Origin", "*")
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

				encodedName := class.HashSHA1(fileHeader.Filename+strconv.FormatInt(time.Now().Unix(), 10)) + ".docx"
				dest := filepath.Join("assets/temp/userUpload/", encodedName)

				log.Info().Str("文件名", encodedName).Msg("API：上传文件")

				if _, err := c.SaveFormFile(fileHeader, dest); err != nil {
					log.Error().Err(err).Msg("保存上传文件失败")
					return
				}

				c.Header("Access-Control-Allow-Origin", "*")
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
			Upload.Get("/docx/parse/{p}", func(c *context.Context) {

				doc, err := document.Open("assets/temp/userUpload/" + c.Params().Get("p"))
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

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(data)

			})

		}

	}

	if err := app.Listen(":4040"); err != nil {
		log.Panic().Err(err).Msg("启动 API 服务失败")
	}

}
