# July
![Go](https://github.com/JackyCZJ/July/workflows/Go/badge.svg)
An api server,base on echo . With mongodb store and cache with redis.

我的毕业设计，想的很丰满现实很骨感的破实现。目前还在整后端，前端有点难受。

thanks to [yowko](https://github.com/yowko/Docker-Compose-MongoDB-Replica-Set) 's docker-compose

## Usage
```bash
git clone https://github.com/JackyCZJ/July.git

go get

make ca        #if you want to use it with https


make all      #fmt , test , and build

./db_init.sh #db init

./july      #run
```

## TODO

### Buyer
- [ ] 商品浏览
- [x] 商品搜索
- [x] 用户登陆注册
- [ ] 购物车增删改查
- [ ] 下单购买
- [x] 用户登录注册
### Seller
- [x] 商家登录
- [ ] 商品上架
- [ ] 订单管理
- [ ] 商品发货
Admin-server
### 
- [x] 管理员登陆
- [ ] 商家审核
- [ ] 商家封禁