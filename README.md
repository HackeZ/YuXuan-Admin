# YuXuan-Admin

## Intro

YuXuan Shop Admin. Base on Beego.

## Usage
```shell
$ bee run
```

Brower : [http://127.0.0.1:8080/public/login](http://127.0.0.1:8080/public/login)

## Dev Log

1. error : must have one register DataBase alias named `default`

    You miss `orm.RegisterDataBase("default", db_type, dns)`