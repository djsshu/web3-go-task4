package main

import (
	"fmt"
	. "go_task4/model"
	. "go_task4/router"
)

func main() {
	db, config := InitDB()
	r := SetupRouter(db)
	addr := fmt.Sprintf(":%d", config.Server.Port)
	r.Run(addr)
}
