# July
An api server,base on echo . With mongodb store and cache with redis.


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
