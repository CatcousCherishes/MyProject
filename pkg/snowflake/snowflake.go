package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

//调用snowflake 雪花分布式生成ID  	 基于雪花算法分布式生成用户ID

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

// 返回 int64 64位的ID值
func GenID() int64 {
	return node.Generate().Int64()
}
