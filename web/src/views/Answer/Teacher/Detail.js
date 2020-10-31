import Axios from 'axios'
import {Chart} from '@antv/g2'

export default {
    name: "Detail",
    data() {
        return {
            Question: { // 题目信息
                loading: true, // 加载状态
                id: this.$route.params.id, // ID
                object: {}, // 对象
                type: '', // 类型
                text: [], // 题干
                optionsDisplay: true, // 是否显示选项
                options: [], // 选项[仅选择题]
                key: '' // 答案
            },
            groupMemList: [], // 群成员
            Status: { // 状态信息
                Tab: {
                    0: {label: '等待发布答题', color: 'green', badge: 'success'},
                    1: {label: '正在答题中', color: 'blue', badge: 'processing'},
                    2: {label: '答题已结束', color: 'red', badge: 'error'}
                },

                status: 0, // 状态值
                answererCount: 0, // 答题人数
                sliderValue: 0, // 状态条值
                sliderLabel: {0: '准备作答', 1: '允许作答', 2: '停止作答'}, // 说明标签

            },
            Statistics: {
                rightRate: 0, // 正确率
                rightCount: 0, // 正确人数
                wrongRate: 0, // 错误率
                wrongCount: 0, // 错误人数
                rightStus: [], // 回答正确的学生
                wrongStus: {} // 回答错误的学生
            },

            CHARTData: {}
        }
    },
    methods: {
        fetchData() { // API 获取数据

            Axios.get('/apis/question/a/' + this.Question.id).then(res => {

                console.log('成功获取答题数据：')
                console.log(res.data)

                Axios.get(`/apis/group/${res.data.target}/mem`).then(res => {

                    console.log('成功获取群成员：')
                    console.log(res.data)

                    let list = {}

                    for (let i = 0; i < res.data.length; i++) {
                        list[res.data[i].id] = res.data[i]
                    }

                    console.log(list)
                    this.groupMemList = list

                }).catch(err => {
                    console.error('获取群成员失败：' + err)
                })

                if (res.data.answer === null) {
                    console.log('没有作答')
                    res.data.answer = []
                }
                this.Question.object = res.data

                try {
                    this.displayQuestion()
                    console.log('数据初始化成功')
                } catch (e) {
                    console.error('数据初始化失败：' + e)
                }

                this.Question.loading = false

            }).catch(err => {
                console.error('获取答题数据失败：' + err)
            })

            this.openws()

        },
        openws() { // 开启 WS 连接

            let ws = new WebSocket(`ws://${location.host}/question`)
            ws.onopen = () => {
                ws.send(String(this.Question.id))
                setInterval(() => {
                    ws.send('keep heart');
                    console.log('WS保持连接');
                }, 50000);
            }
            ws.onmessage = msg => {

                const data = JSON.parse(msg.data)

                console.log('服务器推送问答数据：')
                console.log(data)

                if (data.answer === null) {
                    console.log('没有作答')
                    data.answer = []
                }
                this.Question.object = data

                try {
                    this.uploadQuestion()
                    console.log('数据更新成功')
                } catch (e) {
                    console.error('数据更新失败：' + e)
                }

            }
            ws.onclose = () => {
                console.error('WS连接已关闭')
            }

        },

        displayQuestion() { // 显示问题数据

            { // 题目
                this.Question.type = {0: '选择题', 1: '简答题'}[this.Question.object.type] // 类型

                this.Question.text = JSON.parse(this.Question.object.question) // 题目
                this.Question.key = this.Question.object.key // 选项

                if (this.Question.object.type === 1) {
                    this.Question.optionsDisplay = false
                } else {
                    this.Question.options = JSON.parse(this.Question.object.options)
                }
            }

            this.updateHeader()
            this.calc()

            { // 图表
                this.histogram(0, 'data-chart-right_count',
                    this.Statistics.rightCount, '正确人数 ' + this.Statistics.rightCount)
                this.histogram(1, 'data-chart-wrong_count',
                    this.Statistics.wrongCount, '错误人数 ' + this.Statistics.wrongCount)
            }

        },
        uploadQuestion() {

            this.updateHeader()
            this.calc()

            { // 图表
                this.updateHistogram(0, this.Statistics.rightCount, '正确人数 ' + this.Statistics.rightCount)
                this.updateHistogram(1, this.Statistics.wrongCount, '错误人数 ' + this.Statistics.wrongCount)
            }

        },
        updateHeader() {
            { // 状态
                this.Status.status = this.Question.object.status
                this.Status.answererCount = this.Question.object.answer.length
                this.Status.sliderValue = this.Question.object.status
            }
        },

        calc() {

            const answer = this.Question.object.answer, options = this.Question.options
            let rightCount = 0, wrongCount = 0, rightStus = [], wrongStus = {}

            for (let i = 0; i < answer.length; i++) {

                if (answer[i].answer === this.Question.key) {
                    rightCount++
                    rightStus.push(answer[i].answerer_id)
                } else {
                    wrongCount++

                    // 寻找错误的选项
                    for (let ii = 0; ii < options.length; ii++) {
                        if (options[ii].type !== answer[i].answer) {
                            continue
                        }

                        let list = wrongStus[options[ii].type]
                        if (list === undefined) {
                            list = []
                        }

                        list.push(answer[i].answerer_id)
                        wrongStus[options[ii].type] = list
                    }
                }
            }

            this.Statistics.rightCount = rightCount
            this.Statistics.wrongCount = wrongCount
            this.Statistics.rightStus = rightStus
            this.Statistics.wrongStus = wrongStus

            this.Statistics.rightRate = Math.floor(rightCount / this.Status.answererCount * 100)
            this.Statistics.wrongRate = 100 - this.Statistics.rightRate

        },

        histogram(id, elm, data, text) {

            const chart = new Chart({
                container: elm,
                autoFit: true,
                width: 240
            })

            chart.data([{type: text, value: data}])
            chart.scale('sales', {nice: true})
            chart.interval().position('type*value')

            chart.render()
            this.CHARTData[id] = chart

        },
        updateHistogram(id, data, text) {
            this.CHARTData[id].changeData([{type: text, value: data}])
        },

        changeStatus() {

            const text = ['准备答题', '发布答题', '停止答题']
            const code = ['prepare', 'start', 'stop']
            const c = this.Status.sliderValue

            console.log(text[c])

            Axios.get(`/apis/question/${this.Question.id}/${code[c]}`).then(res => {
                if (res.data.message === 'yes') {
                    this.$notification.success({message: `成功${text[c]}`})
                } else {
                    this.$notification.error({message: `${text[c]}失败`})
                }
            }).catch(err => {
                console.error(`${text[c]}失败：` + err)
                this.$notification.error({message: `${text[c]}失败`})
            })

            this.$notification.info({
                message: `正在${text[c]}中...`
            })

            this.Status.status = this.Status.sliderValue

        },
        cancelChangeStatus() {
            this.Status.sliderValue = this.Status.status
        },

        praise() {
            Axios.get(`/apis/group/${this.Question.object.target}/praise?mem=${JSON.stringify(this.Statistics.rightStus)}`).then(() => {
                this.$notification.success({message: '表扬成功'})
            }).catch(err => {
                this.$notification.success({message: '表扬失败'})
                console.error('表扬失败：' + err)
            })
        }

    },
    mounted() {
        this.fetchData()
    }
}