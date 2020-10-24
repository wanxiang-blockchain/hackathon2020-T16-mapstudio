import Vue from 'vue'
import App from './App.vue'
import router from './router/' //引入路由


Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  render: h => h(App),
}).$mount('#app')
