# socks5--proxy

##author: Dende

###簡單翻牆代理：

####通訊層加密

####tls鏈接

####網絡編程學習

####科学上网

##SERVER服务端

####使用docker运行：

#####<拉取镜像>

````
docker pull godende/go-proxy:0.1
````

#####<运行镜像>

````
sudo docker run -id -p 11080:11080 --name go-proxy godende/go-proxy:0.1
````

##CLIENT客户端

````
需要開啓本地socks主機
終端代理需要添加規則
在~/.bashrc文件中增加以下两句，表示bash终端中的http和https的请求也通过socks5协议进行代理转发。
export http_proxy="socks5://127.0.0.1:2553"
export https_proxy="socks5://127.0.0.1:2553"
````
