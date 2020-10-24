import Vue from 'vue'
import http from './http'

//接口方法定义
let requests = {
  getDetail:(params)=>http.post('',params)
};

Vue.prototype.$http=requests;
export default requests;
