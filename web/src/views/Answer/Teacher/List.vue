<template>
  <div class="wrapper">

    <a-page-header
        title="答题列表"
        sub-title="存储您所有新建过的答题"
        @back="() => null"
    >
      <template slot="extra">
        <a-button type="primary">
          <router-link to="/answer/tea/new">新增答题</router-link>
        </a-button>
      </template>
    </a-page-header>

    <div class="list">
      <a-list item-layout="vertical" :data-source="questionList" :loading="questionListLoad">
        <a-list-item slot="renderItem" slot-scope="item">
          <a-list-item-meta>

            <router-link slot="title" :to="questionAddr(item.id)">
              <a-tag :color="tipColor(item.status)">{{ tipText(item.status) }}</a-tag>
              {{ item.question }}
            </router-link>

          </a-list-item-meta>
        </a-list-item>
      </a-list>
    </div>

  </div>
</template>

<script>
import Axios from 'axios'

export default {
  name: "List",
  data() {
    return {
      questionListLoad: true,
      questionList: []
    }
  },
  methods: {
    fetchQuestionList() {

      Axios.get(`http://localhost:4040/apis/question/list/${this.$cookies.get('account')}`).then(res => {

        this.questionList = res.data
        this.questionListLoad = false

        console.log("成功拉取答题数据：")
        console.log(res.data)

      }).catch(err => {
        console.log("拉取答题数据失败：" + err)
      })

    },
    questionAddr(s) {
      return "/answer/tea/a/" + s
    },
    tipText(s) {
      return {0: '未开始', 1: '答题中', 2: '已结束'}[s]
    },
    tipColor(s) {
      return {0: '#87d068', 1: '#2db7f5', 2: '#f50'}[s]
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

.ant-list-item-meta {
  margin: 0;
}
</style>
