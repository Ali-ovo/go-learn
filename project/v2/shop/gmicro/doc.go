package gmicro

/*
为什么我们不用现有的go-zero、kratos而是要自己开发一个框架
1. go语言还没有像java的spring一样由一个大一统的框架，生态越来越成熟，目前来说到底哪个框架最终能发展成spring一样
2. 这个时候我们的系统依赖了kratos或者go-zero，在公司内可能会面临不同的开发人员选择不同的框架
3. 我们使用了kratos就会受到kratos的限制，gorm，监控等大量的组件， kratos如果提供了这些了功能还好， go-zero
4. kratos这些框架可能无法提供一些特殊的功能
5. 我们可以使用kratos但是不能丧失自己搭建底层框架的能力
6. 熟悉底层各个细节的最好方式是自己搭建底层框架，然后使用自己的框架， 我希望把日志打印到jaeger
7. 很多同学使用kratos的时候，不去理解细节不懂得看源码，逐渐就成为了crud boy
基于上述这些我们决定自己搭建一套微服务框架 gin grpc
*/
