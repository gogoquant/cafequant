const PROXY_CONFIG = [
  {
    context: [
      '/api',
      '/callback',
      '/login',
      '/server'
    ],
    target: 'http://localhost:1323',
    secure: false
  }
  // {
  //   context: [
  //     '/proxy/server/*'
  //   ],
  //   target: 'http://localhost:1323',
  //   secure: false,
  //   pathRewrite: {
  //     '^/proxy': ''
  //   }
  // }
]
module.exports = PROXY_CONFIG
