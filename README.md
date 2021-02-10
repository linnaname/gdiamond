
## Not Ready For Production Env ##
分布式配置管理平台gdiamod

## 工程架构 ##
1.分布式配置开源系统选型比较 TODO
2.gdiamond的架构 TODO
3.各个分布式配置系统如何感知配置变更并推送 TODO

## 功能列表 ##
*配置修改实时生效

用户在修改完配置并发布后，客户端能实时（1秒）接收到最新的配置，并通知到应用程序

* 部署简单

作为基础服务，为了提高可用性和降低运维的难度，目前唯一的外部依赖是MySQL，部署非常简单，只要安装好Golang和MySQL就可以跑起来

*统一管理 TODO

提供统一界面集中式管理配置。

## 快速开始 ##

#### 下载源码 ####
git clone https://github.com/linnaname/gdiamond.git
</br>

#### 部署NameServer ####
1.`cd gdiamond/namesrv/cmd & go build`

2.`./cmd`

看到console有如下输出表示启动成功
``` 
2021/02/02 21:04:11 NameServer is listening on :9000 (multi-cores: true, loops: 4)
2021/02/02 21:04:11 Starting  httpserver
```

#### 部署Server ####

1.启动mysql并创建库名为diamond的库，然后执行server/mysql目录下的init.sql创建config_info表

2.数据库配置在server/etc/gdiamond.toml,按照自己的情况进行修改

3.`cd gdiamond/server/cmd & go build`

4.`./cmd -n 127.0.0.1 -c ../configs/` -n后面是namesrv的地址，多个使用分号分割，-c指定配置文件目录目录内必须有gdiamond.toml

控制台看到如下输出代表启动成功
```2021/02/02 22:00:44 Starting  httpserver```

5.`curl 127.0.0.1:8080:/namesrv/addrs`  查看server是否注册成功， 127.0.0.1是namesrv的ip地址

6.发布配置 `curl -X POST "http://127.0.0.1:1210/diamond-server/publishConfig?dataId=linna&group=DEFAULT_GROUP&content=helloWorld"`

7.获取配置 `curl -X GET "http://127.0.0.1:8848//diamond-server/config?dataId=linna&group=DEFAULT_GROUP"`

#### client使用 ####

1.修改etc/hosts

增加 `127.0.0.1 gdiamond.namesrv.net`
127.0.0.1 为namesrv的ip地址，gdiamond.namesrv.net指无状态的namesrv http服务器，client代码中写死了为gdiamond.namesrv.net，这里也用这个域名

进行集群部署时，可修改为您真正的域名地址，[代码地址](https://github.com/linnaname/gdiamond/blob/master/client/internal/processor/serveraddress.go)

2.发布或更新配置
```golang
cli := client.NewClient()
b := cli.PublishConfig(dataId, group, content)
```
通过b的真假值判断是否发布或修改成功

3.读取配置
```golang
cli := client.NewClient()
content := cli.GetConfig(dataId, group, 1000)
```
content就是配置内容

4.读取配置并设置监听器

实现监听器 ManagerListener
```golang
type A struct {
}

func (a A) ReceiveConfigInfo(configInfo string) {
	println("ReceiveConfigInfo:", configInfo)
}
```

读取并监听
```golang
cli := client.NewClient()
content := cli.GetConfigAndSetListener(dataId, group, 1000, A{})
```
content为配置内容，当然需要程序常驻才可以监听


## bench ##
性能暂时还不能满足生产环境需求，具体bench可以看[本地bench](https://gdiamond.yuque.com/staff-xkx7zo/twhd01/wltq22)

## 一点小背景 ##

diamond在分布式配置管理系统在阿里内部使用非常广泛，而之前网上流传下来的diamond源码是很多年前的难以使用，而[nacos-config](https://nacos.io/zh-cn/index.html)
基本可以看成dimaond在阿里内部不断升级之后的开源版。 如果读过diamond和nacos-config的源码可以发现两者在大致架构和实现思路上并没有太多区别。

如果对[携程的Apollo](https://github.com/ctripcorp/apollo)
有所了解可也可以看得出来它和diamond在架构上有很多相似的地方。

diamond和[disconf](https://github.com/knightliao/disconf)
走的是完全不同的两条路，或者业务场景不同，disconf采用推的方式更适用于对配置更新实时感知的场景，而diamond采用的长轮询拉的方式。

从我的使用经历来看大部的配置更改要求实时性要求并没有那么高（大部分场景1s感知和30ms感知差别并不大），而diamond在架构上的简单和多层的可用性设计却给高可用和维护带来了很大的便利。
比如disconf基于zookeeper满足的是实时性要求比较高的场景，但是引入zookeeper确实也给系统带来了更大的复杂度和维护上的困难，而实时性配置推送并没有那么强的需求。

这似乎有点像AP和CP的选择，没有对错好坏之分，diamond在阿里内部的广泛使用从某种程度上也暗示着diamond走的必然是AP的路。


## Reference ##
* [挺老的一份diamond代码](https://github.com/takeseem/diamond)
* [diamond升级版nacos-config](https://nacos.io/en-us/)
* [Apollo架构分析](https://mp.weixin.qq.com/s/-hUaQPzfsl9Lm3IqQW3VDQ)


## Contributing ##
虽然这是我自己为了把Golang重新捡起来写的一个轮子，但是还是欢迎感兴趣的小伙参与项目贡献！但Golang的生态确实还不太完善，如果您感兴趣欢迎联系我，也可以额提交PR修复一个bug，或者新建 Issue 讨论新特性或者变更。

## TODO List ##
* localfile测试和localfile变更监听测试  done
* client 和 server从namesrv获取全量可用server逻辑  done
* client user api简化     done
* 日志系统完善 done
* namesrv可使用vip + 域名代替 done
* 配置变更实时感知测试（秒级） done
* 租户 增加 6
* daily/pre/production环境增加 7
* 集群环境性能测试  1 done 性能不理想，需要优化
* mysql连接、gnet连接优化  2
* namesrv和server优雅关闭 3
* 用户权限 5
* 管理页面（宜搭） 4
* 历史记录，回滚  8
* 增加fallback机制，在长轮询同时失效时也不会的导致一直无法拉取变更 9
* Docker部署 10