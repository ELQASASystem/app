<template>
  <main class="wrapper">

    <a-page-header
        title="答题列表"
        sub-title="存储您所有新建过的答题"
        @back="() => $router.push({path: '/answer/home'})"
    >
      <template slot="extra">
        <a-button type="primary">
          <router-link to="/market/list">答题市场</router-link>
        </a-button>
        <a-button type="primary">
          <router-link to="/answer/tea/new">新增答题</router-link>
        </a-button>
      </template>
    </a-page-header>

    <div class="list">
      <a-list item-layout="horizontal" :data-source="questionList" :loading="questionListLoad">
        <a-list-item slot="renderItem" slot-scope="item">
          <a-list-item-meta>

            <router-link slot="title" :to="questionAddr(item.id)">
              <a-tag :color="tipColor(item.status)">{{ tipText(item.status) }}</a-tag>
              <a-tag>{{ groupNameTip(item.target) }}</a-tag>
              {{ questionText(item.question) }}
            </router-link>

          </a-list-item-meta>
        </a-list-item>
      </a-list>
    </div>

  </main>
</template>

<script>
import Axios from 'axios'

export default {
  name: "List",
  data() {
    return {
      groupList: {},
      questionListLoad: true,
      questionList: []
    }
  },
  methods: {
    fetchQuestionList() {

      Axios.get('/apis/questions/list?u=' + this.$cookies.get('account')).then(res => {

        for (let i = 0; i < res.data.groups.length; i++) {
          this.groupList[res.data.groups[i].id] = res.data.groups[i].name
        }

        console.log("成功拉取群列表：", this.groupList)

        this.questionList = res.data.questions
        this.questionListLoad = false

        console.log("成功拉取答题数据：", res.data.questions)

      }).catch(err => {
        console.error("获取数据失败：", err)
        this.questionListLoad = false
      })
    },

    questionText(a) {

      a = JSON.parse(a)
      let text = ""

      for (let i = 0; i < a.length; i++) {
        if (a[i].type === "img") {
          text += "[图片]"
          continue
        }

        text += a[i].text
      }

      return text
    },
    questionAddr(s) {
      return "/answer/tea/a/" + s
    },
    tipText(s) {
      return {0: '未开始', 1: '答题中', 2: '已结束'}[s]
    },
    tipColor(s) {
      return {0: '#87d068', 1: '#2db7f5', 2: '#f50'}[s]
    },
    groupNameTip(s) {
      return this.groupList[s]
    }
  },
  mounted() {
    this.fetchQuestionList()
  }
}
</script>

<style scoped>
.list {
  padding: 0 24px;
}
</style>
