## GoMailer
轻量电子邮件推送服务(A lightly email sending service for Go)

通过form提交用户输入的数据，GoMailer会将这些数据填入预先定义好的邮件内容模板中，并帮你把内容投递到指定的邮箱。
也可以选择把邮件暂存在GoMailer中，另外选择时间手动触发投递。

额外的可选配置: 
1 支持开启reCaptcha验证，避免恶意投递<br>
2 配置请求成功或失败时的重定向地址，相应事件发生时用户将被重定向到指定页面

## Release Note
- v0.1.0

## 使用说明
提供三个接口与GoMailer进行交互，epKey获取接口(epKey唯一标识一个服务接入点)，邮件发送接口，邮件查询接口。

#### 1. 获取 epKey
API: `POST /api/shortcut`

将如下json作为 request body, 发送POST请求到`/api/shortcut`接口获取epKey。
```json
{
  "user": {
    "username": "A",
    "password": "123456"
  },
  "app": {
    "name": "sample",
    "host": "sample.com"
  },
  "endpoint": {
    "name": "sample用户反馈",
    "dialer": {
      "host": "smtp.qq.com",
      "port": 465,
      "authUsername": "666@qq.com",
      "authPassword": "xxx",
      "name": "sample用户反馈专用"
    },
    "receiver": [
      {
        "address": "xxx@163.com",
        "receiverType": "To"
      },
      {
        "address": "xxx@gmail.com",
        "receiverType": "Cc"
      },
      {
        "address": "xxx1@gmail.com",
        "receiverType": "Bcc"
      }
    ],
    "template": {
      "contentType": "text/html",
      "template": "<div>来自用户[{{name}}]的反馈, 用户联系方式: {{contact}}, 反馈内容如下:<hr><p>{{content}}</p><div/>"
    },
    "preference": {
      "deliverStrategy": "DELIVER_IMMEDIATELY",
      "enableReCaptcha": 1,
      "successRedirect": "http://www.sample.com/feedback-success.html",
      "failRedirect": "http://www.sample.com/feedback-fail.html"
    }
  }
}
```

字段说明:
- dialer: 邮件发件人配置，需到自己的邮箱网页端自行获取，参考[QQ邮箱的获取方式](https://service.mail.qq.com/cgi-bin/help?subtype=1&id=28&no=1001256)
- dialer.name: 发件人名称
- receiver: 收件人配置，To: 接收人，Cc: 抄送人，Bcc: 密送人
- template: 邮件内容模板配置，类似{{contact}}的部分最终会被form中相应字段值替换
- preference.deliverStrategy: 邮件投递策略: DELIVER_IMMEDIATELY: 立即发送，STAGING: 保存但不发送
- preference.enableReCaptcha: 是否启用reCaptcha验证，1: 启用，2: 不启用
- preference.successRedirect: 邮件发送成功时的重定向地址

请求成功时的返回结果
```text
{
  "epKey": "xxx"
}
``` 
epKey唯一标识一个服务接入点，后续请求都需要将该参数拼接在url中(或form中)传递给服务器。借用上面的请求示例进行说明，如有用户A，是
sample网站的管理员，sample网站有两个地方接入了GoMailer的邮件服务，一个是用户反馈功能，一个是质量投诉功能，那么用户反馈就为一个
接入点，质量投诉为另一个接入点，epKey不同，拥有独立的配置。

可将上述request body保存为文件，后使用register.sh脚本快捷获取epKey:
```shell script
# http://localhost:8080/api/shortcut 为接口地址，注意进行替换
# sample.json 为接入点配置文件
./register.sh http://localhost:6060/api/shortcut sample.json
```

当接入点配置更新时亦可通过该接口进行更新。

#### 2. 在网站中集成
API: `POST /api/mail/send`

将第一步获取到的epKey放入url参数中。
```html
<form action="http://localhost:6060/api/mail/send?epKey=xxxxx" method="post">
    <input name="name" placeholder="该怎么称呼您"/><br>
    <input name="contact" placeholder="联系方式(可不填)"/><br>
    <textarea name="content" placeholder="反馈内容"></textarea><br>
    <input type="hidden" name="grecaptcha_token" value="xxx">
    <input type="submit">
</form>
```
若第一步选择(或后续进行更新)启用reCaptcha，应将reCaptcha token放入`grecaptcha_token`字段提交到服务器，放在form中或拼接在url都可以。
reCaptcha的集成可参考[这里](https://www.cnblogs.com/dulinan/p/12033018.html)

#### 3. 查询邮件
API: `GET /api/mail/list`

请求示例: /api/mail/list?uid=1&pn=1&ps=10
- pn: 分页页码，可不传，从1开始，默认1
- ps：分页页大小，可不传，默认10条

响应数据:
```json
{
    "PageNum":1,
    "PageSize":10,
    "Total":1,
    "List": [
        {
            "InsertTime":"2020-03-21T16:37:58+08:00",  
            "State":"STAGING",
            "DeliveryTime":"2020-03-21T16:37:58+08:00",
            "Content":"<div>来自用户[小马]的反馈, 用户电话号码: 1999999999, 反馈内容如下:<hr><p>不错</p><div/>",
            "Raw":{
                "name":"小马",
                "contact": "1999999999",
                "content": "不错"
            }
        }
    ]
}
```

字段说明:
- InsertTime: 创建时间
- State: 邮件状态，STAGING: 只保存未投递 DELIVER_SUCCESS: 投递成功 DELIVER_FAILED: 投递失败 
- DeliveryTime: 邮件投递时间
- Raw: 对应form中的数据

License
============
```text
                  GNU LESSER GENERAL PUBLIC LICENSE
                       Version 2.1, February 1999

 Copyright (C) 1991, 1999 Free Software Foundation, Inc.
 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.

[This is the first released version of the Lesser GPL.  It also counts
 as the successor of the GNU Library Public License, version 2, hence
 the version number 2.1.]
```
