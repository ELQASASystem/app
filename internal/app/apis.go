package class

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/ELQASASystem/app/internal/app/database"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
	"github.com/unidoc/unioffice/document"
)

// startAPI 开启 API 服务
func startAPI() {

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

				pa := c.Params()
				res, err := database.Class.Account.ReadAccountsList(pa.Get("u"))
				if err != nil {
					log.Error().Err(err).Msg("校验密码失败")
					return
				}

				if pa.Get("p") != res.Password {
					_, _ = c.JSON(iris.Map{"message": "no"})
					return
				}

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(iris.Map{"message": "yes"})

			})

		}

		Group := API.Party("group")
		{

			// 获取群列表
			Group.Get("/list", func(c *context.Context) {

				err := classBot.c.ReloadGroupList()
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
				for _, v := range classBot.c.GroupList {
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
				for _, v := range classBot.c.FindGroupByUin(i).Members {

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

					publishQuestion(uint64(data.ID), data.Question)

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
				expiredQuestion(qid)

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

		}

		Upload := API.Party("upload")
		{

			// 上传 Docx 前预检请求
			Upload.Options("/docx", func(c *context.Context) {

				c.Header("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Headers", "Origin, x-requested-with, Content-Type, Accept")
				c.Header("Access-Control-Allow-Methods", "POST")

				_, _ = c.HTML("")

			})

			// 上传 Docx
			Upload.Post("/docx", func(c *context.Context) {

				c.SetMaxRequestBodySize(10485760) // 10MiB

				_, fileHeader, err := c.FormFile("file")

				encodedName := ""

				if err != nil {
					log.Error().Err(err).Msg("处理文件上传失败")
					return
				}

				uploadTime := time.Now()

				encodedName = string(hashSHA1(fileHeader.Filename + " " + uploadTime.String()))

				dest := filepath.Join("assets/temp/userUpload/", encodedName)

				if _, err := c.SaveFormFile(fileHeader, dest); err != nil {
					log.Error().Err(err).Msg("保存上传文件失败")
					return
				}

				c.Header("Access-Control-Allow-Origin", "*")
				_, _ = c.JSON(iris.Map{"fileName": encodedName})
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

	if err := app.Listen(":8080"); err != nil {
		log.Panic().Err(err).Msg("启动 API 服务失败")
	}

}
