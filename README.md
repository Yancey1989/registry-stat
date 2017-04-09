# registry-stat
根据docker registry的日志分析每个docker image被pull的次数

## Prepare postgreSQL

启动postgreSQL

```bash
mkdir -p registry-stat/data
cd registry-stat
docker run -d -v --name registry-stat-pg $PWD/data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:9
```

创建postgreSQL用户，数据库，表:

```sql
CREATE USER registry_stat;
CREATE DATABASE registry;
GRANT ALL PRIVILEGES ON DATABASE registry TO registry_stat;
```

为用户`registry_stat`配置密码

```sql
ALTER user registry_stat with password '<your password>'
```

创建request数据表

```sql
CREATE TABLE REQUEST(
  requestID VARCHAR,
  timestamp INTEGER,
  remoteAddr VARCHAR(40),
  imageName VARCHAR(128),
  imageTag VARCHAR(128),
  PRIMARY KEY(requestID)
)
```


## Build registry-stat docker image

```bash
go build .
docker bulid -t registry-stat .
```

## Run

```bash
docker run --rm --name registry-stat \
  -e CONTAINER_NAME=registry \
  -e CONTAINER_PATH=/var/lib/docker/ \
  -e RECORD_FILE=/var/lib/docker/registry-stat.pos \
  -e DBCONNECT=registry_stat:paddle@172.31.64.183:5432/registry?sslmode=disable \
  -v /var/lib/docker/:/var/lib/docker/ \
  yancey1989/registry-stat

```
