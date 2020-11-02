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
      <a-list item-layout="horizontal" :data-source="questionList" :loading="questionListLoad">
        <a-list-item slot="renderItem" slot-scope="item">

          <div v-for="q in JSON.parse(item.question)" :key="q.id">
            <span v-if="q.type === 'text'">{{ q.text }}</span>
            <img
                v-else :src="'/assets/question/pictures/'+q.path"
                class="question-img" alt="问题图片"
            >
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
      questionList: [],
      questionListLoad: true
    }
  },
  methods: {
    fetchMarket() {
      Axios.get('/apis/market/list').then(res => {

        console.log('成功获取问题市场数据：')
        console.log(res.data)

        this.questionListLoad = false
        this.questionList = res.data

      }).catch(err => {
        console.error('获取问题市场失败：' + err)
      })

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
}

.btn-add {
  margin: 0 8px;
}
</style>