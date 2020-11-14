<template>
  <main class="wrapper">

    <a-page-header
        title="问题市场"
        sub-title="挑选您喜爱的问题并加入答题列表"
        @back="() => $router.push({path: '/home'})"
    >
      <template slot="extra">
        <a-button type="primary">
          <router-link to="/answer/tea/list">答题列表</router-link>
        </a-button>
      </template>
    </a-page-header>

    <div class="list">

      <a-radio-group
          v-model="subject"
          default-value="1"
          button-style="solid"
          @change="onSubjectChange"
      >
        <a-radio-button value="1">语文</a-radio-button>
        <a-radio-button value="2">数学</a-radio-button>
        <a-radio-button value="3">英语</a-radio-button>
        <a-radio-button value="4">物理</a-radio-button>
        <a-radio-button value="5">历史</a-radio-button>
        <a-radio-button value="6">化学</a-radio-button>
        <a-radio-button value="7">生物</a-radio-button>
        <a-radio-button value="8">地理</a-radio-button>
        <a-radio-button value="9">政治</a-radio-button>
        <a-radio-button value="10">其他</a-radio-button>
      </a-radio-group>

      <a-divider/>

      <a-list item-layout="horizontal" :data-source="questionList" :loading="questionListLoad">
        <a-list-item slot="renderItem" slot-scope="item">

          <div>
            <div v-for="q in JSON.parse(item.question)" :key="q.id">
              <template v-if="q.type === 'text'">{{ q.text }}</template>
              <img
                  v-else :src="'/assets/question/pictures/'+q.path"
                  class="question-img" alt="问题图片"
              >
            </div>
          </div>
          <a-button class="btn-add">添加</a-button>

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
      subject: 0,
      questionList: [],
      questionListLoad: true
    }
  },
  methods: {
    fetchMarket() {
      Axios.get(`/apis/market/${this.subject}/list`).then(res => {

        if (res.data === null) {
          throw '数据为空'
        }
        console.log('成功获取问题市场数据：')
        console.log(res.data)

        this.questionList = res.data

      }).catch(err => {
        console.error('获取问题市场失败：' + err)
        this.$notification.error({message: '这个领域还没有人哦'})
        this.questionList = []
      }).finally(() => {
        this.questionListLoad = false
      })

    },
    onSubjectChange() {
      this.fetchMarket()
    }
  },
  mounted() {
    this.fetchMarket()
  }
}
</script>

<style scoped>
.list {
  padding: 0 24px;
}

.question-img {
  max-width: 100%;
  max-height: 320px;
}

.btn-add {
  margin: 0 8px;
}
</style>