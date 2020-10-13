import Vue from 'vue'
import Antd from 'ant-design-vue'
import Cookie from 'vue-cookies'

import App from './App.vue'
import router from './router'

import 'ant-design-vue/dist/antd.min.css'

Vue.use(Antd)
Vue.use(Cookie)

Vue.config.productionTip = false

new Vue({
    router,
    render: h => h(App)
}).$mount('#app')
