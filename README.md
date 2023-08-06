### 大学网课/成人教育网课/湘潭大学网课全自动代看

#### 环境要求
- go 1.16
- go mod

#### 使用说明一
1. 把项目下载下来 `git clone https://github.com/qinjintian/superchutou.git`

1. 进入到`superchutou`目录
   
1. 把`configs/config.yaml.dtst`命名为`configs/config.yaml`

1. 构建项目 `go build -o superchutou.exe .\cmd\`（在此操作前必须先安装go环境）后会在当前目录下生成一个 superchutou.exe 可执行文件(windows 10 环境，其他环境自行百度构建)

1. 双击运行 `superchutou.exe` 即可

PS: 目前已支持手机号+验证码方式登录，双击运行`superchutou.exe`后，根据终端提示输入手机号和验证码即可登录

#### 使用说明二
#### _**对于不想看代码实现过程的朋友可以直接下载 `superchutou.exe` 和 `configs/config.yaml.dtst`到本地电脑（windows），并把`configs/config.yaml.dtst`命名为`configs/config.yaml`，需要注意的是`configs/config.yaml`和`superchutou.exe`必须在同一级目录，最后双击运行`superchutou.exe`根据终端提示操作即可**_

#### 代看原理
_**该程序自动代看的原理与浏览器多窗口播放是一样的，程序内实现以守护进程方式在后台自动模拟数据提交到视频接口，忽略掉了视频播放界面，从而达到了观看视频进度更新**_

#### 截图

![image](https://github.com/qinjintian/superchutou/blob/main/运行截图.png?raw=true)