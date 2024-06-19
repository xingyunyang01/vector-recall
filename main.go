package main

import (
	"fmt"
	"log"

	"github.com/xingyunyang01/vector-recall/pkg/aihelper"
	"github.com/xingyunyang01/vector-recall/pkg/strhelper"
	"github.com/xingyunyang01/vector-recall/pkg/tairhelper"
)

func main() {
	//0.创建向量数据库索引
	indexName := "Higress"
	var err error
	//err = tairhelper.GetIndex(indexName)
	//if err != nil {
	// err = tairhelper.CreateIndex(indexName, 1536, "HNSW", "L2", "FLOAT32", true, 16)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//}

	//1.准备prompt
	prompt := "什么是Higress?"
	// // prompt := "总结一下这份文档的内容"
	// //prompt := "汉堡好吃吗"
	chat := aihelper.Chat(prompt)
	fmt.Println(chat)

	//2.向量化
	textInputs := []string{prompt, chat}
	vector, err := aihelper.GetVec(textInputs)
	if err != nil {
		log.Fatal(err)
	}

	//3.插入向量数据库
	svector := strhelper.FloatSliceToString(vector)
	fields := make(map[string]interface{})
	fields["VECTOR"] = svector
	fields["question"] = prompt
	fields["answer"] = chat
	tairhelper.Insert(indexName, prompt, fields)

	fmt.Println("//////////////////////////////////////////////////")
	//第二轮提问
	prompt2 := "介绍一下Higress?"

	//2.向量化
	vector, err = aihelper.SimpleGetVec(prompt2)
	if err != nil {
		log.Fatal(err)
	}

	svector = strhelper.FloatSliceToString(vector)

	//3.在向量数据库中搜索
	key, err := tairhelper.VectorSearch(indexName, svector)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("key: " + key)

	answer, err := tairhelper.KeySearch(indexName, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("answer: " + answer)
	//4.如果找到了(距离<0.1)

}
