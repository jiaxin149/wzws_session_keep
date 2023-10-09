# wzws_session_keep
奇安信网站卫士session保持程序

网站卫士的所有操作都要携带PHPSESSID

但是由于有效期很短，只有30分钟左右，所以此程序用于保持PHPSESSID处于活跃状态，使开发者能随时、自由的调用网站卫士WEB的所有接口

# 使用方法
1. 初次启动程序会在运行目录生成``go_conf.json``文件，填写对应参数再次2. 运行程序即可

3. 成功启动后访问``127.0.0.1:14911/get_phpsessid``即可获取到PHPSESSID数据

4. 最后就可以开始愉快的使用了！

## 使用示例 
1. 获取网站卫士域名列表
```php
<?php
$phpsessid = file_get_contents("127.0.0.1:14911/get_phpsessid");
$url = 'https://wangzhan.qianxin.com/domain/getDomainList';
// 将Cookie添加到GET请求的头部
$options = array(
    'http' => array(
        'header'  => "Cookie: PHPSESSID=".$phpsessid,
        'method'  => 'GET',
    ),
);
// 携带PHPSESSID发送HTTP GET请求并获取响应内容
$response = file_get_contents($url, false, stream_context_create($options));
// 输出响应内容
echo $response;
```
………………
