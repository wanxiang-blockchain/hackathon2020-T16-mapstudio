import Vue from 'vue'
import App from './App.vue'
import router from './router/' //引入路由
import api from './api/api'

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  api,
  render: h => h(App),
}).$mount('#app')
