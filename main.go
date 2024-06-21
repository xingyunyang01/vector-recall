package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xingyunyang01/vector-recall/pkg/aihelper"
	"github.com/xingyunyang01/vector-recall/pkg/strhelper"
	"github.com/xingyunyang01/vector-recall/pkg/tairhelper"
)

func main() {
	//0.创建向量数据库索引
	indexName := "Higress"

	result, err := tairhelper.GetIndex(indexName)
	if err != nil {
		log.Fatal(err)
	} else {
		if len(result) == 0 {
			err = tairhelper.CreateIndex(indexName, 1536, "HNSW", "L2", "FLOAT32", true, 16)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("-------------------Round 1------------------")
	//1.准备prompt
	prompt := "可以动态修改Higress的路由配置吗？"

	//2.转向量
	vector, err := aihelper.SimpleGetVec(prompt)
	if err != nil {
		log.Fatal(err)
	}

	svector := strhelper.FloatSliceToString(vector)

	//3.在向量数据库中搜索
	key, err := tairhelper.VectorSearch(indexName, svector)
	if err != nil {
		log.Fatal(err)
	}

	//4.没搜到
	if key == "" {
		//4.1 请求chatgpt
		chat := aihelper.Chat(prompt)

		//4.2 将返回的json中包含的question和answer拆解出来
		var output aihelper.OutPutFormat
		if chat != "" {
			err = json.Unmarshal([]byte(chat), &output)
			//返回的确实是json
			if err == nil {
				//4.3.向量化
				textInputs := []string{output.Question, output.Answer}
				vector, err := aihelper.GetVec(textInputs)
				if err != nil {
					log.Fatal(err)
				}

				//4.4.插入向量数据库
				svector := strhelper.FloatSliceToString(vector)
				fields := make(map[string]interface{})
				fields["VECTOR"] = svector
				fields["question"] = output.Question
				fields["answer"] = output.Answer
				tairhelper.Insert(indexName, output.Question, fields)
				fmt.Println("请求通义千问的结果是: " + chat)
			} else {
				fmt.Println("请求通义千问的结果是: " + chat)
			}
		}
	} else {
		//搜到了
		//根据key把answer取出来
		answer, err := tairhelper.KeySearch(indexName, key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("搜到了，答案是：" + answer)
	}

	fmt.Println("-------------------Round 2------------------")
	//第二轮提问
	prompt2 := "怎么操作呢"

	//2.向量化
	vector, err = aihelper.SimpleGetVec(prompt2)
	if err != nil {
		log.Fatal(err)
	}

	svector = strhelper.FloatSliceToString(vector)

	//3.在向量数据库中搜索
	key, err = tairhelper.VectorSearch(indexName, svector)
	if err != nil {
		log.Fatal(err)
	}

	//4.没搜到
	if key == "" {
		//4.1 请求chatgpt
		chat := aihelper.Chat(prompt)

		//4.2 将返回的json中包含的question和answer拆解出来
		var output aihelper.OutPutFormat
		if chat != "" {
			err = json.Unmarshal([]byte(chat), &output)
			//返回的确实是json
			if err == nil {
				//4.3.向量化
				textInputs := []string{output.Question, output.Answer}
				vector, err := aihelper.GetVec(textInputs)
				if err != nil {
					log.Fatal(err)
				}

				//4.4.插入向量数据库
				svector := strhelper.FloatSliceToString(vector)
				fields := make(map[string]interface{})
				fields["VECTOR"] = svector
				fields["question"] = output.Question
				fields["answer"] = output.Answer
				tairhelper.Insert(indexName, output.Question, fields)
				fmt.Println("请求通义千问的结果是: " + chat)
			} else {
				fmt.Println("请求通义千问的结果是: " + chat)
			}
		}
	} else {
		//搜到了
		//根据key把answer取出来
		answer, err := tairhelper.KeySearch(indexName, key)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("搜到了，答案是：" + answer)
	}
}
