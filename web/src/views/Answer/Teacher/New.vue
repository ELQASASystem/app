<template>
  <main class="wrapper">

    <a-page-header
        title="新建答题"
        sub-title="新建一个答题"
        @back="() => $router.push({path: '/answer/tea/list'})"
    >
      <template slot="extra">
        <router-link to="/answer/tea/list">
          <a-button type="primary">答题列表</a-button>
        </router-link>
      </template>
    </a-page-header>

    <div class="detail">

      <a-tabs default-active-key="1">
        <a-tab-pane key="1" tab="单个新建">

          <a-form-model :model="form" :label-col="labelCol" :wrapper-col="wrapperCol">

            <a-form-model-item label="作答群">
              <a-select mode="multiple" placeholder="请选择欲参与答题的群">
                <!--网络请求获取群列表-->
                <a-select-option value="1000">英语群</a-select-option>
                <a-select-option value="1001">物理群</a-select-option>
              </a-select>
            </a-form-model-item>

            <a-form-model-item label="题目">
              <a-input v-model="form.desc" type="textarea"/>
            </a-form-model-item>

            <a-form-model-item label="题目图片">
              <a-upload
                  action="https://www.mocky.io/v2/5cc8019d300000980a055e76"
                  list-type="picture-card"
              >
                <div>
                  <a-icon type="plus"/>
                  <div class="ant-upload-text">
                    Upload
                  </div>
                </div>
              </a-upload>
              <a-modal>
                <img alt="example" style="width: 100%"/>
              </a-modal>


            </a-form-model-item>

            <a-form-model-item label="题目类型">
              <a-radio-group v-model="form.resource">
                <a-radio value="1">选择题</a-radio>
                <a-radio value="2" disabled>简答题</a-radio>
                <a-radio value="3" disabled>多选题</a-radio>
              </a-radio-group>
            </a-form-model-item>

            <a-form-model-item label="选项">
              <a-button type="dashed" @click="addDomain">
                <a-icon type="plus"/>
                新增选项
              </a-button>
              <a-button type="dashed" @click="removeOp">
                <a-icon type="plus"/>
                删除选项
              </a-button>
            </a-form-model-item>
            <a-form-model-item
                v-for="(domain, i) in dynamicValidateForm.domains"
                :label="abc[i]"
                :key="domain.key"
            >
              <a-input placeholder="请输入选项答案"/>
            </a-form-model-item>

            <a-form-model-item label="答案">

              <a-select default-value="A">
                <a-select-option value="A">A</a-select-option>
                <a-select-option value="B">B</a-select-option>
              </a-select>

            </a-form-model-item>

            <a-form-model-item label="分享至问题市场">
              <a-switch v-model="form.delivery"/>
            </a-form-model-item>


            <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
              <a-button type="primary" @click="onSubmit">
                创建新题目
              </a-button>
            </a-form-model-item>
          </a-form-model>


        </a-tab-pane>

        <a-tab-pane key="2" tab="批量新建">

          <div class="upload">
            <a-upload-dragger
                accept=".docx"
                action="/apis/upload/docx"
                @change="uploadDocx"
            >
              <p class="ant-upload-drag-icon">
                <a-icon type="inbox"/>
              </p>
              <p class="ant-upload-text">
                点击或拖拽文件到这里上传
              </p>
              <p class="ant-upload-hint">
                仅支持单个（.docx）文件上传
              </p>
            </a-upload-dragger>
          </div>

        </a-tab-pane>


      </a-tabs>

    </div>


  </main>


</template>


<script>
export default {
  name: "New",
  data() {
    return {
      labelCol: {span: 4},
      wrapperCol: {span: 14},
      form: {
        name: '',
        region: undefined,
        date1: undefined,
        delivery: false,
        type: [],
        resource: '',
        desc: '',

      },


      abc: ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
        'W', 'X', 'Y', 'Z'],

      formItemLayout: {
        labelCol: {
          xs: {span: 24},
          sm: {span: 4},
        },
        wrapperCol: {
          xs: {span: 24},
          sm: {span: 20},
        },
      },
      formItemLayoutWithOutLabel: {
        wrapperCol: {
          xs: {span: 24, offset: 0},
          sm: {span: 20, offset: 4},
        },
      },
      dynamicValidateForm: {
        domains: [{
          value: '', key: Date.now(),
        }, {
          value: '', key: Date.now(),
        }],
      },

    }
  },
  methods: {
    uploadDocx(info) {
      console.log(info)
    },

    addDomain() {
      this.dynamicValidateForm.domains.push({
        value: '',
        key: Date.now(),
      });
    },
    removeOp() {
      this.dynamicValidateForm.domains.pop()
    },
    onSubmit(e) {
      console.log(e)
    }
  }
}
</script>

<style scoped>

</style>