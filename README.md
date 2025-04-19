# FileBrowser - 轻量级文件浏览器

Filebrowser是一个轻量级的文件浏览器，支持 Windows / Linux 系统，及兼容 S3 的对象存储（如 Minio）。特别适用于服务器上文件的存取，可替代 FTP、vsFTP、SMB等文件服务器，以及 Docker / Kuberntes Pod 中 PVC 挂载盘上文件的存取。

# 功能
- Windows 设备和驱动器浏览、文件浏览、下载、搜索、创建目录、文件上传
- Linux 服务器上文件浏览、下载、搜索、创建目录、文件上传，支持启动时指定根路径
- 基于 JwtToken 的账号登录授权（管理员、只读、读写三种角色）
- 支持兼容 S3 的对象存储（页面可配置对接多个 S3）上的文件浏览、下载、搜索、创建目录、文件上传
- 因安全考虑，暂不支持文件删除操作

# 使用方法
## 二进制文件启动
直接在 [release](http://github.com/hongyuxuan/filebrowser/releases) 页面下载编译好的二进制文件，解压后命令行启动：

<b>windows</b>
```powershell
# 使用解压工具解压 filebrowser-windows-amd64-v1.1.0.tar.gz
# 打开 CMD 或 Powershell
.\filebrowser.exe

# 或指定端口、数据库文件启动
.\filebrowser.exe --port 9121 -d \path\to\your\filebrowser.db

# 或直接双击 filebrowser.exe 启动
```
<b>linux</b>
```shell
tar zxf filebrowser-linux-amd64-v1.1.0.tar.gz
./filebrowser 

# 或指定端口、数据库文件启动
./filebrowser --port 9121 -d /path/to/your/filebrowser.db
```

<b>docker</b>
```shell
docker run -v /data:/data -p 9121:9121 hongyuxuan/filebrowser:v1.1.0 --port 9121 -d /data/filebrowser.db
```

# 编译
在 linux 环境下，拉取源码到本地后，先编译前端（需要 node20 以上）：
```shell
/bin/bash ci/scripts/build-web.sh
```

再编译后端（需要 go1.23+、make）：
```shell
/bin/bash ci/scripts/build.sh
```