package models

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"os"
	"strconv"
)

// Snowflake 노드 초기화
var node *snowflake.Node

func init() {
	nodeID := int64(0)
	if nodeIDstring := os.Getenv("SNOWFLAKE_NODE_ID"); nodeIDstring == "" {
		nodeID = int64(os.Getpid()) % 1024
	} else {
		if n, err := strconv.Atoi(nodeIDstring); err != nil {
			panic(fmt.Sprintf("Invalid SNOWFLAKE_NODE_ID: %s", nodeIDstring))
		} else {
			nodeID = int64(n)
		}
	}
	var err error
	// 노드 번호 1로 Snowflake 노드 생성
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		panic(err)
	}
}
