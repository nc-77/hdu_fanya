## 简介

本仓库实现了一个HDU泛雅平台学习通上未完成作业的定时邮件提示系统

## 配置

1. Fork本[仓库](https://github.com/nc-77/hdu_fanya)

2. 设置QQ邮箱权限

   QQ邮箱需要开启SMTP服务，并取得授权码作为之后的邮箱登录密码。

   该邮箱将作为定时提示服务的发送方。

   具体操作步骤可参考 [官方教程](https://service.mail.qq.com/cgi-bin/help?subtype=1&&no=1001256&&id=28)

3. 添加Secrets

   在Fork后的仓库内的Settings->Secrets 添加以下字段。

   | 字段名        | 说明                    |
   | ------------- | ----------------------- |
   | MAIL_ADDRESS  | 接收邮件地址            |
   | MAIL_USERNAME | 发送邮件地址            |
   | MAIL_PASSWORD | 发送邮件登陆密码/授权码 |
   | HDU_USERNAME  | 数字杭电账号            |
   | HDU_PASSWORD  | 数字杭电密码            |


4. 验证

   在Fork后的仓库内的Actions启用workflow功能。

   启动workflow功能后可手动run进行验证。

   ![](https://img.nc-77.top/20210423220739.png)

   workflow 运行成功后，该系统会在每天早上7点通过邮件发送学习通上未完成的作业进行提醒。