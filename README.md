# GoMailer
轻量电子邮件推送服务(A lightly email sending service for Go)
<br>

### shortcut to register an endpoint for email delivery
request body
```json
{
  "user": {
    "username": "djn",
    "password": "123456"
  },
  "app": {
    "appName": "demo",
    "host": "demo.com"
  },
  "endPoint": {
    "name": "发送反馈",
    "dialer": {
      "host": "smtp.qq.com",
      "port": 465,
      "authUsername": "666@qq.com",
      "authPassword": "666aaa",
      "name": "XX公司"
    },
    "receiver": [
      {
        "address": "djn163<duan_jia_ning@163.com>",
        "receiverType": "TO"
      },
      {
        "address": "djn163<foo@163.com>",
        "receiverType": "CC"
      }
    ],
    "template": {
      "contentType": "text/html",
      "template": "<div><hr><h1>Test email{{msg}}</h1><div/>"
    },
    "preference": {
      "deliverStrategy": "IMMEDIATELY",
      "enableReCaptcha": true
    }
  }
}
```

response body
```json
{
  "AppKey": "MToxOmRlbW86MTU4NDAwMTI2NDY4NTQ3MjYwMA"
}
``` 