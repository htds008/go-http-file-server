# Go HTTP File Server
基于命令行的HTTP文件共享服务器。

![Go HTTP File Server pages](doc/ghfs.gif)

## 软件特色
- 比Apache/Nginx更友好的目录列表
- 适配移动设备显示
- 单一可执行文件
- 可以将当前浏览目录内容打包下载
- 可以开启某个目录的上传权限
- 可以指定自定义模板来渲染页面
- 支持目录别名（将另一个目录挂载到某个URL路径）

## 编译
至少需要Go 1.9版本。
```bash
go build src/main.go
```
会在当前目录生成"main"可执行文件。

如果修改过`src/tpl`下的默认html模板，需要将其重新嵌入go文件：
```bash
cd src
make tpls
```
然后再像上面那样编译。

## 举例
在8080端口启动服务器，根目录为当前工作目录：
```sh
server -l 8080
``` 

在8080端口启动服务器，根目录为 /usr/share/doc：
```sh
server -l 8080 -r /usr/share/doc
```

在默认端口启动服务器，根目录为/tmp，并允许上传文件到/tmp/data：
```sh
server -r /tmp -u /data
```

共享/etc下的文件，同时把/usr/share/doc挂载到URL路径/doc下：
```sh
server -r /etc -a :/doc:/usr/share/doc
```

在8080端口启动服务器，使用HTTPS协议：
```sh
server -k /path/to/certificate/key -c /path/to/certificate/file -l 8080
```

不显示`.`开头的unix隐藏目录和文件。提示：用引号括起通配符以避免shell展开：
```sh
server -H '.*'
```

在命令行显示访问日志：
```sh
server -L -
```

http基本验证：
- 对URL /files 启用验证
- 用户名：user1，密码：pass1
- 用户名：user2，密码：pass2
```sh
server --auth /files --user user1:pass1 --user-sha1 user2:8be52126a6fde450a7162a3651d589bb51e9579d
```

启动2台虚拟主机：
- 服务器1
    - 在80端口提供http服务
    - 在443端口提供https服务
        - 证书文件：/cert/server1.pem
        - 私钥文件：/cert/server1.key
    - 主机名：server1.example.com
    - 根目录：/var/www/server1
- 服务器2
    - 在80端口提供http服务
    - 在443端口提供https服务
        - 证书文件：/cert/server2.pem
        - 私钥文件：/cert/server2.key
    - 主机名：server2.example.com
    - 根目录：/var/www/server2
```sh
server --listen-plain 80 --listen-tls 443 -c /cert/server1.pem -k /cert/server1.key --hostname server1.example.com -r /var/www/server1 ,, --listen-plain 80 --listen-tls 443 -c /cert/server2.pem -k /cert/server2.key --hostname server2.example.com -r /var/www/server2
```

## 使用方法
```
server [选项]

-l|--listen <IP|端口|:端口|IP:端口|socket> ...
    指定服务器要侦听的IP和端口，例如“:80”或“127.0.0.1:80”。
    如果指定了--cert和--key，端口接受TLS连接。
    如果未指定端口，则在纯HTTP模式下使用80端口，TLS模式下使用443端口。
    如果值中包含“/”，则将其当作unix socket路径。
    标志“-l”或“--listen”可以省略。
--listen-plain <IP|端口|:端口|IP:端口|socket> ...
    与--listen类似，但强制使用非TLS模式。
--listen-tls <IP|端口|:端口|IP:端口|socket> ...
    与--listen类似，但强制使用TLS模式。若未指定证书和私钥，则启动失败。

--hostname <主机名> ...
    指定与当前虚拟主机关联的主机名。
    如果值以“.”开头，则将其当作后缀，匹配该域下的所有子域名，例如“.example.com”。
    如果值以“.”结尾，则将其当作前缀，匹配所有域名后缀。

-r|--root <目录>
    服务器的根目录。
    默认为当前目录。

-R|--empty-root
    使用空的虚拟目录作为根目录。
    在仅需挂载别名的情况下较实用。

--default-sort <排序规则>
    指定文件和目录的默认排序规则。
    可用的排序key：
    - `n` 按名称递增排序
    - `N` 按名称递减排序
    - `e` 按类型（后缀）递增排序
    - `E` 按类型（后缀）递减排序
    - `s` 按大小递增排序
    - `S` 按大小递减排序
    - `t` 按修改时间递增排序
    - `T` 按修改时间递减排序
    - `_` 不排序
    目录顺序：
    - `/<key>` 目录在文件之前
    - `<key>/` 目录在文件之后
    - `<key>` 目录与文件混合

-I|--dir-index <文件> ...
    指定目录默认页面文件。

-a|--alias <分隔符><URL路径><分隔符><文件系统路径> ...
    设置路径别名。
    将某个文件系统路径挂载到URL路径下。
    例如：“:/doc:/usr/share/doc”。

-U|--global-upload
    对所有URL路径开启上传权限。
    请谨慎使用。
-u|--upload <URL路径> ...
    设置允许上传的URL路径（及子路径）。
    请谨慎使用。
--upload-dir <文件系统路径> ...
    与--upload类似，但指定的是文件系统路径，而不是URL路径。

    上传选项注意事项：
        如果名称已存在且是常规文件，
        若已启用删除（例如--delete选项），则尝试先删除文件，
        否则尝试添加或递增数字后缀。

--global-mkdir
    对所有URL路径开启创建子目录权限。
--mkdir <url-path> ...
    设置允许创建子目录的URL路径（及子路径）。
--mkdir-dir <fs-path> ...
    与--mkdir类似，但指定的是文件系统路径，而不是URL路径。

    创建子目录选项注意事项：
        为避免歧义，被别名遮蔽的目录名不能被创建。

--global-delete
    对所有URL路径开启删除子项权限。
--delete <url-path> ...
    设置允许删除子项的URL路径（及子路径）。
--delete-dir <fs-path> ...
    与--delete类似，但指定的是文件系统路径，而不是URL路径。

    删除选项注意事项：
        为避免歧义，URL路径下挂载的别名不能被删除。
        别名下的非别名文件/目录仍然可以被删除。
        为避免歧义，被别名遮蔽的正常文件/目录不能被删除。

-A|--global-archive
    对所有URL路径开启打包下载当前目录内容的功能。
    页面顶部会出现下载链接。
    请确保符号链接没有循环引用。
--archive <URL路径> ...
    对指定URL路径（及子路径）开启打包下载当前目录内容的功能。
--archive-dir <文件系统路径> ...
    与--archive类似，但指定的是文件系统路径，而不是URL路径。

--global-cors
    接受所有URL路径的CORS跨域请求。
--cors <url-path> ...
    接受指定URL路径（及子路径）的CORS跨域请求。
--cors-dir <fs-path> ...
    接受指定文件系统路径（及子路径）的CORS跨域请求。

--global-auth
    对所有URL路径启用http基本验证(Basic Auth)。
--auth <url-path> ...
    对指定URL路径（及子路径）启用http基本验证。
--auth-dir <fs-path> ...
    对指定文件系统路径（及子路径）启用http基本验证。
--user [<用户名>]:[<密码>] ...
    为当前虚拟主机指定用于http基本验证的用户，允许空的用户名和/或密码。
--user-base64 [<用户名>]:[<base64密码>] ...
--user-md5 [<用户名>]:<md5密码> ...
--user-sha1 [<用户名>]:<sha1密码> ...
--user-sha256 [<用户名>]:<sha256密码> ...
--user-sha512 [<用户名>]:<sha512密码> ...
    指定http基本验证的用户，对密码使用特定的编码。

-c|--cert <证书文件>
    指定TLS证书文件。

-k|--key <私钥文件>
    指定TLS私钥文件。

-t|--template <模板文件>
    指定用于渲染页面的自定义模板，代替内建模板。

-S|--show <通配符> ...
-SD|--show-dir <通配符> ...
-SF|--show-file <通配符> ...
    如果指定该选项，只有匹配通配符的目录或文件（除了被hide选项隐藏的）才会显示出来。

-H|--hide <wildcard> ...
-HD|--hide-dir <wildcard> ...
-HF|--hide-file <wildcard> ...
    如果指定该选项，匹配通配符的目录或文件不会显示出来。

-L|--access-log <文件>
    访问日志。
    使用“-”指定为标准输出。
    设为空来禁用。

-E|--error-log <文件>
    错误日志。
    使用“-”指定为标准错误输出。
    设为空来禁用。
    默认为“-”。

--config <文件>
    为当前虚拟主机指定外部配置文件。

    其内容为任何其他选项，
    与在命令行指定的形式相同，
    用空白符分割。

    外部配置的优先级低于命令行选项。
    如果在命令行指定了某个选项，则其外部配置被忽略。

,,
    要指定多台虚拟主机的选项，用此符号分割每台主机的选项。
    可以为每台虚拟主机分别指定以上选项。

    如果多个虚拟主机共享相同的IP和端口，
    使用--hostname指定主机名，根据请求头中的主机名来区分虚拟主机。
    如果请求的主机名不匹配任何虚拟主机，
    服务器尝试使用第一个没有指定主机名的虚拟主机，
    如果失败则使用第一个虚拟主机。
```
