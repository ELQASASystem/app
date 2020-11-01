<template>
  <main class="wrapper">

    <a-page-header
        title="答题详情"
        sub-title="查看答题详细数据"
        @back="() => $router.push({path: '/answer/tea/list'})"
    >

      <template slot="extra">
        <router-link to="/answer/tea/list">
          <a-button type="primary">答题列表</a-button>
        </router-link>
      </template>

      <a-spin :spinning="Question.loading">

        <div class="question-info">
          <h2>题目 <span class="question-type">[{{ Question.type }}]</span></h2>

          <h3>
            问题：
            <span v-for="item in Question.text" :key="item.id">
              <span v-if="item.type === 'text'">{{ item.text }}</span>
              <img
                  v-else :src="'/assets/question/pictures/'+item.path"
                  class="question-img" alt="问题图片"
              >
            </span>
          </h3>
          <h3 v-if="Question.optionsDisplay">
            选项：
            <a-tag
                color="blue"
                v-for="item in Question.options"
                :key="item.type"
            >
              {{ item.type + "：" + item.body }}
            </a-tag>
          </h3>
          <h3>答案：{{ Question.key }}</h3>

        </div>

        <div class="question-status">
          <h2>状态</h2>

          <div>
            <a-tag :color="Status.Tab[Status.status].color">
              <a-badge :status="Status.Tab[Status.status].badge"/>
              <span>{{ Status.Tab[Status.status].label }}</span>
            </a-tag>
            <a-tag color="blue">已作答人数：{{ Status.answererCount }}</a-tag>
          </div>

          <div class="question-status-change">
            <a-popconfirm
                title="你确定更改答题状态吗"
                ok-text="确定"
                cancel-text="取消"
                @confirm="changeStatus"
                @cancel="cancelChangeStatus"
            >

              <a-slider
                  :max="2"
                  :marks="Status.sliderLabel"
                  :step="null"
                  :default-value="0"
                  v-model="Status.sliderValue"
              />

            </a-popconfirm>
          </div>
        </div>

      </a-spin>
    </a-page-header>

    <div class="detail">

      <a-tabs default-active-key="1">
        <a-tab-pane key="1" tab="统计">

          <div class="detail-block">
            <h2>正确</h2>
            <div class="data-chart">

              <a-progress
                  type="circle"
                  :percent="Statistics.rightRate"
                  status="success"
                  :format="p => { return '正确率 ' + p + '%' }"
              />

              <div id="data-chart-right_count"/>

            </div>
            <a-collapse :bordered="false">
              <a-collapse-panel header="回答正确的学生">
                <a-button slot="extra" @click="praise">
                  <a-icon type="like"/>
                  表扬他们
                </a-button>

                <a-list item-layout="horizontal" :data-source="Statistics.rightStus">
                  <a-list-item slot="renderItem" slot-scope="item">
                    <span>{{ groupMemList[item].name }}</span>
                  </a-list-item>
                </a-list>


              </a-collapse-panel>

            </a-collapse>
          </div>

          <div class="detail-block">
            <h2>错误</h2>
            <div class="data-chart">

              <a-progress
                  type="circle"
                  :percent="Statistics.wrongRate"
                  status="exception"
                  :format="p => { return '错误率 ' + p + '%' }"
              />

              <div id="data-chart-wrong_count"/>
            </div>

            <a-collapse :bordered="false">
              <a-collapse-panel header="回答错误的学生">

                <a-collapse v-for="(item, key) in Statistics.wrongStus" :key="item.length" :bordered="false">
                  <a-collapse-panel :header="key">

                    <a-list item-layout="horizontal" :data-source="item">
                      <a-list-item slot="renderItem" slot-scope="item">
                        <span>{{ groupMemList[item].name }}</span>
                      </a-list-item>
                    </a-list>

                  </a-collapse-panel>
                </a-collapse>

              </a-collapse-panel>
            </a-collapse>

          </div>

          <div class="detail-block data-chart">
            <a-card
                size="small" title="易错选项"
                class="question-easy_wrong"
            >
              {{ this.Statistics.wrongList[0].type }}
            </a-card>

            <div id="data-chart-easy_wrong"/>
          </div>

        </a-tab-pane>

        <a-tab-pane key="2" tab="全部信息">

          <a-list item-layout="horizontal" :data-source="Question.object.answer">
            <a-list-item slot="renderItem" slot-scope="item">
              <div>
                <p>
                  {{ groupMemList[item.answerer_id].name }}
                </p>
              </div>
            </a-list-item>
          </a-list>

        </a-tab-pane>


      </a-tabs>

    </div>

  </main>
</template>

<script src="./Detail.js"/>

<style scoped src="./Detail.css"/>