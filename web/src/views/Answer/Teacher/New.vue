<template>
  <div class="wrapper">

    <a-page-header
        title="新建答题"
        sub-title="新建一个答题"
        @back="() => history.back()"
    >
      <template slot="extra">
        <a-button type="primary">答题列表</a-button>
      </template>
    </a-page-header>

    <div class="basic">


    </div>


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


  </div>


</template>


<script>
export default {
  name: "New",
  methods: {
    uploadDocx(info) {

      console.log(info)

    }
  }
}
</script>

<style scoped>

</style>