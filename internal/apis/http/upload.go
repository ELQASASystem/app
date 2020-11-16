package http

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ELQASASystem/app/internal/app"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
	"github.com/unidoc/unioffice/document"
)

type upload struct{}

// Upload 上传
func Upload() *upload { return new(upload) }

// options 上传预检
func (u *upload) options(c *context.Context) {
	c.Header("Access-Control-Allow-Headers", "x-requested-with")
	c.Header("Access-Control-Allow-Methods", "POST")
}

// docx 上传 Docx
func (u *upload) docx(c *context.Context) {

	c.SetMaxRequestBodySize(10485760) // 限制最大上传大小为 10MiB

	_, fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("文件上传失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	encodedName := class.HashForSHA1(fileHeader.Filename+strconv.FormatInt(time.Now().Unix(), 10)) + ".docx"
	dest := filepath.Join("web/assets/temp/docx/", encodedName)

	log.Info().Str("文件名", encodedName).Msg("API：上传文件")

	if _, err := c.SaveFormFile(fileHeader, dest); err != nil {
		log.Error().Err(err).Msg("保存上传文件失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	// 在一分钟后删除该文件
	time.AfterFunc(time.Minute, func() {
		log.Info().Str("文件名", encodedName).Msg("API：删除上传的文件")
		if err := os.Remove(dest); err != nil {
			log.Error().Err(err).Msg("删除文件失败")
		}
	})

	_, _ = c.JSON(iris.Map{"fileName": encodedName})
}

// parseDocx 解析 Docx
func (u *upload) parseDocx(c *context.Context) {

	doc, err := document.Open("web/assets/temp/docx/" + c.Params().Get("p"))
	if err != nil {
		log.Error().Err(err).Msg("打开 Docx 失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
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

}

// picture 上传图片
func (u *upload) picture(c *context.Context) {

	c.SetMaxRequestBodySize(4194304) // 限制最大上传大小为 4MiB

	_, fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("文件上传失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	encodedName := class.HashForSHA1(fileHeader.Filename+strconv.FormatInt(time.Now().Unix(), 10)) + "-" + fileHeader.Filename
	dest := filepath.Join("web/assets/question/pictures/", encodedName)

	log.Info().Str("文件名", fileHeader.Filename).Msg("API：上传文件")

	if _, err := c.SaveFormFile(fileHeader, dest); err != nil {
		log.Error().Err(err).Msg("保存上传文件失败")
		_, _ = c.JSON(iris.Map{"message": "no"})
		return
	}

	_, _ = c.JSON(iris.Map{"fileName": fileHeader.Filename})
}

// split 分词
func (u *upload) split(c *context.Context) {

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
}
