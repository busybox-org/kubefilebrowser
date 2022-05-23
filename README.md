# kubefileborwser

kubernetes container file browser. Is a simple web application that allows you to browse and edit files in a kubernetes container. 

## Parameters

+ `RUN_MODE`: run mode, options: `dev`, `prod`
+ `HTTP_ADDR`: listen address, default: `:8080`
+ `HTTP_PORT`: listen port, default: `8080`
+ `IP_WHITE_LIST`: access ip white list, default: `*` (all).
+ `KUBECONFIG`: k8s config file path, default: `$HOME/.kube/config`

## Run In docker

```shell
docker pull xmapst/kubefilebrowser:latest
docker run -d --restart=always -p 9999:9999 -e RUN_MODE=debug -v /path/to/kubeconfig:/root/.kube/config xmapst/kubefilebrowser:latest
```

## Deploy in kubernetes

```bash
kubectl apply -f deploy/kubefilebrowser.yaml
```

## Index.html
![kubefilebrowser_index_html](https://raw.githubusercontent.com/xmapst/kubefilebrowser/main/img/index_html.jpg)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fxmapst%2Fkubefilebrowser.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fxmapst%2Fkubefilebrowser?ref=badge_shield)

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


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fxmapst%2Fkubefilebrowser.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fxmapst%2Fkubefilebrowser?ref=badge_large)