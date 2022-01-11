# kubefileborwser

kubernetes container file browser

## 启动可选环境变量

| 名称 | 类型 | 默认值 | 说明 |
| ---- | ---- | ---- | ---- |
| RUN_MODE | string | debug | 运行模式 |
| HTTP_ADDR | string | 0.0.0.0 | 监听地址 |
| HTTP_PORT | string | 9999 | 监听端口 |
| IP_WHITE_LIST | []string | * | 访问白名单 |
| KUBECONFIG | string | ~/.kube/config | k8s连接文件路径 |

+ 部署在k8s内创建使用管理员权限的serviceaccount即可

## Index.html
![kubefilebrowser_index_html](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/img/index_html.jpg)

## file_browser
![kubefilebrowser](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/img/file_browser.jpg)

## terminal
![terminal](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/img/terminal.jpg)

## Swagger

![kubefilebrowser swagger image](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/img/swagger_index.jpg)

## Reference documents

+ [golang 1.16 gin static embed](https://mojotv.cn/golang/golang-html5-websocket-remote-desktop)
+ [vue](https://cli.vuejs.org/config/)
+ [kubectl copy & shell 原理讲解](https://www.yfdou.com/archives/kuberneteszhi-kubectlexeczhi-ling-gong-zuo-yuan-li-shi-xian-copyhe-webshellyi-ji-filebrowser.html)
