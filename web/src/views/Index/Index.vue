<template>
  <div class="wrapper">
    <header class="header">
      <a href="/" class="logo">ELQASA</a>
      <div>
        <ul class="menu">
          <li><a href="/">产品官网</a></li>
          <li><a :href="GitHubAddr">GitHub Project</a></li>
          <li><a href="/">About us</a></li>
        </ul>
      </div>
      <div></div>
    </header>
    <main>

      <div class="banner">
        <div class="banner-text">
          <h1>在线教育问答统计深度学习分析系统</h1>
          <span>ELQASASystem 是帮助教师、家长和学生在线学习，易于使用的教育跟踪深度学习分析系统</span>
          <div class="banner-btn">
            <router-link to="/start">
              <a class="banner-btn1">开始使用</a>
            </router-link>
            <a :href="GitHubAddr">
              <a class="banner-btn2">了解更多</a>
            </a>
          </div>
        </div>
        <div class="banner-img">
          <img :src="BannerImg">
        </div>

      </div>

      <div class="login">

        <div class="login-form">
          <a-form
              layout="inline"
              :form="form"
              @submit="handleSubmit"
          >

            <a-form-item>
              <a-input
                  v-decorator="['userName', {rules: [{required: true, message: '请输入用户名!' }]}]"
                  placeholder="用户名"
              >
                <a-icon slot="prefix" type="user" style="color: rgba(0,0,0,.25)"/>
              </a-input>
            </a-form-item>
            <a-form-item>
              <a-input
                  v-decorator="['password', {rules: [{required: true, message: '请输入密码!'}]}]"
                  type="password"
                  placeholder="密码"
              >
                <a-icon slot="prefix" type="lock" style="color: rgba(0,0,0,.25)"/>
              </a-input>
            </a-form-item>

            <a-form-item>
              <a-button type="primary" html-type="submit">登录</a-button>
            </a-form-item>

          </a-form>
        </div>

        <div class="login-banner">登录后开始使用...</div>

      </div>


      <div class="features">
        我们的功能
      </div>
    </main>
  </div>
</template>

<script>
import Axios from 'axios'

export default {
  name: "Index",
  data() {
    return {
      GitHubAddr: 'https://github.com/ELQASASystem/app',
      BannerImg: require('../../assets/images/Index/banner.png'),
      form: this.$form.createForm(this, {})
    }
  },
  methods: {
    handleSubmit(e) {

      e.preventDefault()
      this.form.validateFields((err, values) => {

        if (err) {
          console.log("登录失败：校验失败")
          return
        }

        console.log("登录请求数据：")
        console.log(values)

        this.login(values.userName, values.password)

      })
    },

    login(u, p) {

      Axios.get(`/apis/sign/${u}/in?p=${p}`).then(res => {

        console.log("登录返回信息：")
        console.log(res.data)

        if (res.data.message === "yes") {

          this.$cookies.set('account', u, '30d')
          console.log("登录成功")
          this.$notification.success({message: '登录成功'})
          this.$router.push({path: '/answer/tea/list'})

        } else {
          console.log("登录失败")
          this.$notification.error({message: '登录失败，请检查帐号密码'})
        }

      }).catch(err => {
        console.log("登录请求错误：" + err)
      })
    }
  }
}
</script>

<style scoped src="./Index.css"/>
