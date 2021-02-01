
## Not Ready For Production Env ##
gdiamond可以视为淘宝分布式配置管理diamond的Go实现

## 工程架构 ##


## 功能列表 ##

## 详细文档 ##

## 快速开始 ##

## 性能测试 ##


## 一点小背景 ##

diamond在分布式配置管理系统在阿里内部使用非常广泛，而之前网上流传下来的diamond源码是很多年前的难以使用，而[nacos-config](https://nacos.io/zh-cn/index.html)
基本可以看成dimaond在阿里内部不断升级之后的开源版。 如果读过diamond和nacos-config的源码可以发现两者在大致架构和实现思路上并没有太多区别。

如果对[携程的Apollo](https://github.com/ctripcorp/apollo)
有所了解可也可以看得出来它和diamond在架构上有很多相似的地方。

diamond和[disconf](https://github.com/knightliao/disconf)
走的是完全不同的两条路，或者业务场景不同，disconf采用推的方式更适用于对配置更新实时感知的场景，而diamond采用的长轮询拉的方式。

从我的使用经历来看大部的配置更改要求实时性要求并没有那么高（大部分场景1s感知和10s感知差别并不大），而diamond在架构上的简单和多层的可用性设计却给高可用和维护带来了很大的便利。
比如disconf基于zookeeper满足的是实时性要求比较高的场景，但是引入zookeeper确实也给系统带来了更大的复杂度和维护上的困难，而实时性配置推送并没有那么强的需求。

这似乎有点像AP和CP的选择，没有对错好坏之分，diamond在阿里内部的广泛使用从某种程度上也暗示着diamond走的必然是AP的路。


## Reference ##
* [挺老的一份diamond代码](https://github.com/takeseem/diamond)
* [diamond升级版nacos-config](https://nacos.io/en-us/)
* [Apollo架构分析](https://mp.weixin.qq.com/s/-hUaQPzfsl9Lm3IqQW3VDQ)


## Contributing ##
虽然这是我自己为了把Golang重新捡起来写的一个轮子，但是还是欢迎感兴趣的小伙参与项目贡献！但Golang的生态确实还不太完善，如果您感兴趣欢迎联系我，也可以额提交PR修复一个bug，或者新建 Issue 讨论新特性或者变更。