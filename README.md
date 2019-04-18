# socks5--proxy

簡單翻牆代理：/n
通訊層加密/n
tls鏈接/n
網絡編程學習-----------------------------------------------------------/n
author: Dende----------------ps------------------------------------ps/n
需要開啓本地socks主機/n
終端代理需要添加規則/n
在~/.bashrc文件中增加以下两句，表示bash终端中的http和https的请求也通过socks5协议进行代理转发。/n
export http_proxy="socks5://127.0.0.1:2553"/n
export https_proxy="socks5://127.0.0.1:2553"/n
