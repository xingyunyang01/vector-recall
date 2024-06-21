package tairhelper

import (
	"context"
	"fmt"

	"github.com/alibaba/tair-go/tair"
	"github.com/go-redis/redis/v8"
)

var tairClient *tair.TairClient // 全局客户端

var ip = "r-bp1f6egs00szv63y1tpd.redis.rds.aliyuncs.com"

func init() {
	tairClient = tair.NewTairClient(&redis.Options{
		Addr:     ip + ":" + "6379",
		Username: "r-bp1f6egs00szv63y1t",
		Password: "flzx3qcYSYHL6T#", // no password set
		DB:       0,                 // use default DB
	})
}

// 创建索引
func CreateIndex(indexName string, dim int, algorithm string, distance_method string, data_type string, auto_gc bool, M int32) error {
	_, err := tairClient.TvsCreateIndex(context.Background(), indexName, dim, algorithm, distance_method,
		tair.TvsCreateIndexArgs{}.New().AutoGc(auto_gc).M(M).DataType(data_type)).Result()

	if err != nil {
		return err
	}

	return nil
}

// 获取索引
func GetIndex(indexName string) error {
	result := tairClient.TvsGetIndex(context.Background(), indexName)
	err := result.Err()
	if err != nil {
		return err
	}

	return nil
}

// 插入
// fields[VECTOR]存放向量
// fields[question]存放问题
// fields[answer]存放答案
func Insert(indexName string, key string, fields map[string]interface{}) error {
	_, err := tairClient.TvsHSet(context.Background(), indexName, key, tair.TvsHSetArgs{}.New().Fields(fields)).Result()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// 使用向量搜索相似的key
func VectorSearch(indexName string, vector string) (string, error) {
	result, err := tairClient.TvsKnnSearch(context.Background(), indexName, 1, vector, tair.TvsKnnSearchArgs{}.New().MaxDist(0.2)).Result()
	//result, err := tairClient.TvsKnnSearch(context.Background(), indexName, 1, vector, nil).Result()
	if err != nil {
		return "", err
	}

	fmt.Println(result)
	key, ok := result[0].(string)
	if ok {
		return key, nil
	}
	return "", nil
}

// 使用key拿到内容
func KeySearch(indexName string, key string) (string, error) {
	result, err := tairClient.TvsHGetAll(context.Background(), indexName, key).Result()
	if err != nil {
		return "", err
	}

	// 转换为字符串切片以方便处理
	strSlice := make([]string, len(result))
	for i, v := range result {
		strSlice[i], _ = v.(string) // 假设所有元素都是字符串
	}

	// 查找 "answer" 字符串的索引
	answerIndex := -1
	for i, v := range strSlice {
		if v == "answer" {
			answerIndex = i
			break
		}
	}

	// 检查是否找到了 "answer" 并获取其后的字符串
	if answerIndex != -1 && answerIndex+1 < len(strSlice) {
		nextString := strSlice[answerIndex+1]
		return nextString, nil
	} else {
		fmt.Println("未找到 'answer' 或其后没有字符串")
	}

	return "", nil
}
