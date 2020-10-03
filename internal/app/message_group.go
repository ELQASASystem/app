package class

import "strings"

// monitorGroup 监听群消息
func monitorGroup() {
	for {
		go processGroup(<-*classBot.msgChan)
	}
}

// processGroup 处理群消息
func processGroup(m *QQMsg) {

	if block(m) {
		return
	}

	if m.Chain[0].Text == ".hello" {
		classBot.SendGroupMsg(NewText("Hello, Client!").To(m.Group.ID))
		return
	}

	if strings.HasPrefix(m.Chain[0].Text, ".tts") {
		textSlice := strings.Split(m.Chain[0].Text, "")

		if len(textSlice) > 0 {
			ttsWord := textSlice[1]
			classBot.SendGroupMsg(NewMsg().AddTTSAudio(ttsWord).To(m.Group.ID))
		}

		return
	}

	// 处理答案
	handleAnswer(m)

}

// block 阻止可能的意外
func block(m *QQMsg) bool {

	// 当长度小于1时消息无法获取
	if len(m.Chain) < 1 {
		return true
	}

	// 匿名用户禁止
	if m.User.ID == 80000000 {
		return true
	}

	return false

}
