<template>
  <div id="app">
    <router-view/>
    <footer class="footer">
      Copyright ©2020 <a href="https://github.com/ELQASASystem">ELQASASystem</a>. All Rights Reserved.<br>
      第 13 届龙岩市青少年科技创新大赛作品. 版权所有，侵权必究.
    </footer>
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 0;
}

#app {
  display: flex;
  flex-flow: column;
  align-items: center;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.wrapper {
  width: 100%;
  max-width: 1333px;
}

.footer {
  margin: 16px;
  text-align: center;
  word-break: break-all;
}
</style>
