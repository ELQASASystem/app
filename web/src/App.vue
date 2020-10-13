<template>
  <div id="app">
    <router-view/>
    <footer class="footer">
      Copyright ©2020 ELQASASystem. All Rights Reserved. 侵权必究.<br>
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
  text-align: center;
  word-break: break-all;
}
</style>
