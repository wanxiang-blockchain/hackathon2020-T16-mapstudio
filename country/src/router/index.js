 /*
 2 * 路由对象模块
 3 * */
import Vue from 'vue'
import VueRouter from 'vue-router'

/*引入*/
const Home = ()=>import('../components/Home');
const Detail = ()=>import('../components/Detail');


 //申明使用插件
 Vue.use(VueRouter)

export default new VueRouter({
      routes:[
          {
              path: '/',
              redirect: '/home' //系统默认页
         },
        {
              path:'/home',
          component: Home,
        },
          {
              path: '/detail',
              component:Detail
          }
   ]
 })
