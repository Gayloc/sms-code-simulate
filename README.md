# sms-code-simulate
### 模拟验证码发送与登录的过程

## 已经实现的功能
- 能够输出正确的提⽰信息
- 限制单个⼿机号当⽇最多发送 5 次验证码
- 包含字⺟和数字的验证码
- 对验证码进⾏加密后再存储(使用md5加密)
- 对数据(登录信息)进⾏持久化

## 目前存在的小问题
- 主菜单中输入多个字符(如aaa)会重复出现多个错误信息