package main

import (
    "github.com/DGHeroin/finder"
    "log"
)

func main() {
    e := finder.NewEngine("my-data")
    defer e.CloseAll()

    r := e.Bucket("my-bucket")

    r.Index("100", "所以, 你好, 再见")
    r.Index("101", "Hd God, 你好")
    r.Index("102", "ACDC", "你", "好")
    r.Index("103", "Beyond", "你好")
    r.Index("104", "こんにちは世界")
    r.Index("105", "你好世界")

    r.Index("200", "google is nice search engine")
    r.Index("201", "nignx is best web service engine")
    r.Flush()
    // r.Remove("200")
    // r.Flush()
    log.Println("find:", r.Find("你好"))
    log.Println("find:", r.Find("世界"))
    log.Println("find:", r.Find("engine"))


}
