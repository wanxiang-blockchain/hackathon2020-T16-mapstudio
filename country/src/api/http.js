import Axios from 'axios'

let baseURL = 'http://47.92.225.192:8547/'; //生产
let baseURLChain = 'https://chainservice.digitalchainall.com/'; //区块查询相关
Axios.defaults.headers['Content-Type'] = 'application/json';
Axios.defaults.timeout = 15000;

let request = function(path, params, isPost, domain) {
  if(domain == 'chain'){
    Axios.defaults.baseURL = baseURLChain;
  }
  else{
    Axios.defaults.baseURL = baseURL;
  }
  if (params === undefined) {
    params = {}
  }
  let headers = {
    'language': 'zh'
  };
  return new Promise(resolve => {
    if (isPost) {
      return Axios.post(path, params, {
        headers: headers
      }).then(
          res => {
            resolve(res.data);
          },
          err => {
            resolve({ code: -1, message: '请求错误',text: JSON.stringify(err) });
          });
    } else {
      return Axios.get(path, {
        params: params,
        headers: headers
      }).then(
          res => {
            resolve(res.data);
          },
          err => {
            resolve({ code: -1, message: '请求错误',text: JSON.stringify(err) });
          });
    }
  });
};


export default {
  post(path, params,isChain) {
    console.log(`path`,path,`params`,params)
    return request(path, params, true,isChain);
  }
}




