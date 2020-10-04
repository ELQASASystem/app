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
        <div class="bannerText">
          <h1>在线教育问答统计深度学习分析系统</h1>
          <span>ELQASASystem 是帮助教师、家长和学生在线学习，易于使用的教育跟踪深度学习分析系统</span>
          <div class="bannerBtn">
            <router-link to="/start">
              <a class="bannerBtn1">开始使用</a>
            </router-link>
            <a :href="GitHubAddr">
              <a class="bannerBtn2">了解更多</a>
            </a>
          </div>
        </div>
        <div class="bannerImg">
          <img :src="BannerImg">
        </div>

      </div>
      <div class="features">
        我们的功能
      </div>
    </main>
  </div>
</template>

<script>
export default {
  name: "Index",
  data() {
    return {
      GitHubAddr: 'https://github.com/ELQASASystem/app',
      BannerImg: require('../../assets/images/Index/banner.png'),
    }
  }
}
</script>

<style scoped src="./Index.css"/>
