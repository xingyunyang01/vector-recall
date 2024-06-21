package aihelper

import (
	"context"
	"fmt"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

var MessageStore ChatMessages

func init() {
	MessageStore = make(ChatMessages, 0)
	MessageStore.Clear() //清理和初始化

}

// chat对话
func Chat(prompt string) string {
	c := NewOpenAiClient()
	MessageStore.AddForUser(prompt)
	for k, _ := range MessageStore {
		fmt.Println(MessageStore[k].Msg.Role)
		fmt.Println(MessageStore[k].Msg.Content)
	}
	rsp, err := c.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		//Model: "qwen-long",
		Model:    "qwen2-57b-a14b-instruct",
		Messages: MessageStore.ToMessage(),
	})
	if err != nil {
		log.Println(err)
		return ""
	}

	// var output OutPutFormat
	// err = json.Unmarshal([]byte(rsp.Choices[0].Message.Content), &output)
	// if err != nil {
	// 	log.Println(err)
	// 	return ""
	// }

	//fmt.Println("output.Question: ", output.Question)
	//fmt.Println("output.Answer: ", output.Answer)

	//MessageStore.AddForUser(output.Question)

	//MessageStore.AddForAssistant(strings.Trim(output.Answer, "\n"))
	MessageStore.AddForAssistant(rsp.Choices[0].Message.Content)

	return rsp.Choices[0].Message.Content
}

// 定义chat模型
type ChatMessages []*ChatMessage
type ChatMessage struct {
	Msg openai.ChatCompletionMessage
}

// 枚举出角色
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

// 定义人设
func (cm *ChatMessages) Clear() {
	*cm = make([]*ChatMessage, 0) //重新初始化

	//cm.AddForSystem("你是一个非常有帮助的阿里云Higress产品助手，我会问你关于Higress产品的相关问题，请务必确保只回答Higress产品相关问题，如果问到其他问题，则回复抱歉，我无法回复此问题")
	//cm.AddForSystem("我们玩一个角色扮演的游戏，你是一个有帮助的美食助手，请给我回答美食相关的问题")
	cm.AddForSystem(instruction + outputFormat + example)
	//cm.AddForSystem("你是一个有帮助的Higress产品助手，我会问你Higress产品相关问题，请结合文档给我回答。注意：如果我问的问题，缺少主谓宾，但结合上下文，你能理解，则给我补全主谓宾后，按照如下格式返回答案。例如：我问：Higress可以替换kubernetes吗？你说：可以。我继续问：Spring Cloud Gateway呢？此时Spring Cloud Gateway呢？就缺少了主谓宾，且结合上下文能知道这句话的意思是Spring Cloud Gateway能替换Nginx Ingress吗？此时你给我回答的格式是：question:Spring Cloud Gateway能替换Nginx Ingress吗？answer:可以")
}

// 添加角色和对应的prompt
func (cm *ChatMessages) AddFor(msg string, role string) {
	*cm = append(*cm, &ChatMessage{
		Msg: openai.ChatCompletionMessage{
			Role:    role,
			Content: msg,
		},
	})
}

// 添加Assistant角色的prompt
func (cm *ChatMessages) AddForAssistant(msg string) {
	cm.AddFor(msg, RoleAssistant)

}

// 添加System角色的prompt
func (cm *ChatMessages) AddForSystem(msg string) {
	cm.AddFor(msg, RoleSystem)
}

// 添加User角色的prompt
func (cm *ChatMessages) AddForUser(msg string) {
	cm.AddFor(msg, RoleUser)
}

// 组装prompt
func (cm *ChatMessages) ToMessage() []openai.ChatCompletionMessage {
	ret := make([]openai.ChatCompletionMessage, len(*cm))
	for index, c := range *cm {
		ret[index] = c.Msg
	}
	return ret
}

// 得到返回的消息
func (cm *ChatMessages) GetLast() string {
	if len(*cm) == 0 {
		return "什么都没找到"
	}

	return (*cm)[len(*cm)-1].Msg.Content
}
