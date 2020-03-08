# GoMailer
轻量电子邮件推送服务(A lightly email sending service for Go)
<br>

### shortcut to register an endpoint for email delivery
request body
```json
{
  "username": "djn",
  "password": "123456",
  "app": {
    "app_name": "demo",
    "host": "demo.com"
  },
  "end_point": {
    "name": "发送反馈",
    "dialer": {
      "host": "smtp.qq.com",
      "port": 465,
      "auth_username": "666@qq.com",
      "auth_password": "666aaa",
      "name": "XX公司"
    },
    "receiver": [
      {
        "address": "djn163<duan_jia_ning@163.com>"
      },
      {
        "address": "djn163<foo@163.com>",
        "receiver_type": "CC"
      }
    ],
    "template": {
      "content_type": "text/html",
      "template": "<div><hr><h1>Test email{{msg}}</h1><div/>"
    },
    "preference": {
      "deliver_strategy": "IMMEDIATELY",
      "enable_re_captcha": true
    }
  }
}
``` 

response body
```json
{
  "delivery_key": "adfdsaV.sdaawerzvsdf"
}
``` 