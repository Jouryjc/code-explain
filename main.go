package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {

	// 加载 .env 文件
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error load config file:", err)
		os.Exit(1)
	}

	// 解析命令行参数
	inputPath := flag.String("input", "", "Path to the code file")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Please provide the --input flag with a path to the code file")
		os.Exit(1)
	}

	// 读取代码文件内容
	codeBytes, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	code := string(codeBytes)

	// 初始化OpenAI客户端
	config := openai.DefaultConfig(os.Getenv("OPENAI_KEY"))
	config.BaseURL = os.Getenv("OPENAI_URL")

	c := openai.NewClientWithConfig(config)

	str := fmt.Sprintf("解释以下go代码：\n\n```go\n%s\n```\n要求对每一行代码添加注释说明，注释和实际代码换行显示。只返回注释后的代码不需要其他的辅助说明", code)
	log.Print(str)

	// 通过OpenAI的ChatGPT API获取代码注释
	response, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: str,
				},
			},
		},
	)

	if err != nil {
		fmt.Println("Error getting code explanation:", err)
		os.Exit(1)
	}

	// 解析并格式化API响应以获取代码注释
	fmt.Println(response)
	result := response.Choices[0].Message.Content

	// 将result的结果写入到输出文件中，输出文件是 result/output.go
	outputFile, err := os.Create("result/output.go")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer outputFile.Close()
	_, err = outputFile.WriteString(result)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}
	fmt.Println("Successfully wrote code explanation to output file")
}
