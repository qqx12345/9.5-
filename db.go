package main

import (
    "encoding/json"
    "log"
    "gorm.io/datatypes"
)

func Query(id string) ([]string, error) {
    cli, err := GlobalMySQLPool.Get()
    if err != nil {
        log.Println("连接失败:", err)
        return nil, err
    }
    defer GlobalMySQLPool.Put(cli)

    var user User
    result := cli.Where("open_id = ?", id).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }

    var comments []string
    json.Unmarshal([]byte(user.Comments), &comments);

    return comments, nil
}

func Insert(id string, comments []string) error {
    // 1. 转成 JSON
    data, _ := json.Marshal(comments)

    user := User{
        OpenID:   id,
        Comments: datatypes.JSON(data),
    }

    cli, err := GlobalMySQLPool.Get()
    if err != nil {
        log.Println("连接失败:", err)
        return err
    }
    defer GlobalMySQLPool.Put(cli)

    result := cli.Create(&user)
    return result.Error
}