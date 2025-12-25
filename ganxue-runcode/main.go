package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-redis/redis/v8"
)

var db *redis.Client

type OutputData struct {
	Status int    `json:"status"`
	Output string `json:"output"`
}

func main() {
	db = redis.NewClient(&redis.Options{
		Addr:     RedisCfg.Addr,
		Password: RedisCfg.Password,
		DB:       RedisCfg.DB,
	})
	defer db.Close()
	// 测试连接
	str, _err := db.Ping(context.Background()).Result()
	if _err != nil {
		log.Println("redis连接失败", str, _err)
		os.Exit(1)
	}
	log.Println("redis连接成功", str)

	go runGoCode()
	go runCppCode()
	go runCCode()

	select {}
}

// getCodeAndInput 获取用户代码和输入
func getCodeAndInput(id string) (string, string, error) {
	userCode, err := db.Get(context.Background(), "userCode"+id).Result()
	if err != nil {
		return "", "", fmt.Errorf("get userCode失败: %w", err)
	}
	input, err := db.Get(context.Background(), "input"+id).Result()
	if err != nil {
		return "", "", fmt.Errorf("get input失败: %w", err)
	}
	return userCode, input, nil
}

// runCode 运行用户代码
func runCode(cmd *exec.Cmd, id string, input string, timeoutCTX context.Context) {
	// 交互式输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println("创建stdin失败", err)
		output, err := json.Marshal(OutputData{
			Status: 1,
			Output: "创建stdin失败",
		})
		if err != nil {
			log.Println("marshal output失败", err)
			return
		}
		err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
		if err != nil {
			log.Println("RPush失败", err)
			return
		}
		err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
		if err != nil {
			log.Println("设置过期时间失败", err)
			return
		}
		return
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Println("执行用户代码失败", err)
		output, err := json.Marshal(OutputData{
			Status: 1,
			Output: "执行用户代码失败",
		})
		if err != nil {
			log.Println("marshal output失败", err)
			return
		}
		err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
		if err != nil {
			log.Println("RPush失败", err)
			return
		}
		err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
		if err != nil {
			log.Println("设置过期时间失败", err)
			return
		}
		return
	}
	if _, err := stdin.Write([]byte(input)); err != nil {
		log.Println("写入标准输入失败", err)
		output, err := json.Marshal(OutputData{
			Status: 1,
			Output: "写入标准输入失败",
		})
		if err != nil {
			log.Println("marshal output失败", err)
			return
		}
		err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
		if err != nil {
			log.Println("RPush失败", err)
			return
		}
		err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
		if err != nil {
			log.Println("设置过期时间失败", err)
			return
		}
		return
	}
	if err := stdin.Close(); err != nil {
		log.Println("关闭标准输入失败", err)
		output, err := json.Marshal(OutputData{
			Status: 1,
			Output: "关闭标准输入失败",
		})
		if err != nil {
			log.Println("marshal output失败", err)
			return
		}
		err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
		if err != nil {
			log.Println("RPush失败", err)
			return
		}
		err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
		if err != nil {
			log.Println("设置过期时间失败", err)
			return
		}
		return
	}

	if err := cmd.Wait(); err != nil {
		if errors.Is(timeoutCTX.Err(), context.DeadlineExceeded) {
			log.Println("执行超时", err)
			output, err := json.Marshal(OutputData{
				Status: 1,
				Output: "执行超时",
			})
			if err != nil {
				log.Println("marshal output失败", err)
				return
			}
			err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
			if err != nil {
				log.Println("RPush失败", err)
				return
			}
			err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				return
			}
			return
		}
		log.Println("执行出错", err)
		output, err := json.Marshal(OutputData{
			Status: 1,
			Output: "执行出错",
		})
		if err != nil {
			log.Println("marshal output失败", err)
			return
		}
		err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
		if err != nil {
			log.Println("RPush失败", err)
			return
		}
		err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
		if err != nil {
			log.Println("设置过期时间失败", err)
			return
		}
		return
	}

	if stderr.String() != "" {
		log.Println("执行出错", err)
		output, err := json.Marshal(OutputData{
			Status: 2,
			Output: stderr.String(),
		})
		if err != nil {
			log.Println("marshal output失败", err)
			return
		}
		err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
		if err != nil {
			log.Println("RPush失败", err)
			return
		}
		err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
		if err != nil {
			log.Println("设置过期时间失败", err)
			return
		}
		return
	}

	actualOut := stdout.String()
	log.Println("执行成功", actualOut)
	output, err := json.Marshal(OutputData{
		Status: 0,
		Output: actualOut,
	})
	if err != nil {
		log.Println("marshal output失败", err)
		return
	}
	err = db.RPush(context.Background(), "runcode-finish"+id, output).Err()
	if err != nil {
		log.Println("RPush失败", err)
		return
	}
	err = db.Expire(context.Background(), "runcode-finish"+id, 15*time.Minute).Err()
	if err != nil {
		log.Println("设置过期时间失败", err)
		return
	}
}

func runGoCode() {
	for {
		id, err := db.BLPop(context.Background(), 0*time.Second, "runcode-go").Result()
		if err != nil {
			fmt.Println("BLPop失败", err)
			continue
		}
		userCode, input, err := getCodeAndInput(id[1])
		if err != nil {
			log.Println("getCodeAndInput失败, ", err)
			output, err := json.Marshal(OutputData{
				Status: 1,
				Output: err.Error(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}

		// 将用户代码写入文件中
		if err := os.WriteFile("/home/ganxue/ganxue-runcode/go/main.go", []byte(userCode), 0644); err != nil {
			log.Println("写入文件失败", err)
			output, err := json.Marshal(OutputData{
				Status: 1,
				Output: err.Error(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
				continue
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}
		timeoutCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// 编译代码
		buildCmd := exec.CommandContext(timeoutCTX, "go", "build", "-o", "/home/ganxue/ganxue-runcode/go/main", "/home/ganxue/ganxue-runcode/go/main.go")

		// 设置编译的输出捕获
		var buildOutput bytes.Buffer
		buildCmd.Stdout = &buildOutput
		buildCmd.Stderr = &buildOutput

		// 执行编译
		if err := buildCmd.Run(); err != nil {
			// 编译失败，记录错误信息
			log.Printf("编译失败: %v, 输出: %s", err, buildOutput.String())
			// 处理编译失败逻辑
			output, err := json.Marshal(OutputData{
				Status: 2,
				Output: buildOutput.String(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
				continue
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}

		// 执行编译后的程序
		cmd := exec.CommandContext(timeoutCTX, "/home/ganxue/ganxue-runcode/go/main")
		runCode(cmd, id[1], input, timeoutCTX)
		cancel()
	}
}

func runCppCode() {
	for {
		id, err := db.BLPop(context.Background(), 0*time.Second, "runcode-cpp").Result()
		if err != nil {
			fmt.Println("BLPop失败", err)
			continue
		}
		userCode, input, err := getCodeAndInput(id[1])
		if err != nil {
			log.Println("getCodeAndInput失败, ", err)
			output, err := json.Marshal(OutputData{
				Status: 1,
				Output: err.Error(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}

		// 将用户代码写入文件中
		if err := os.WriteFile("/home/ganxue/ganxue-runcode/cpp/main.cpp", []byte(userCode), 0644); err != nil {
			log.Println("写入文件失败", err)
			output, err := json.Marshal(OutputData{
				Status: 1,
				Output: err.Error(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
				continue
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}
		timeoutCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// 编译代码
		buildCmd := exec.CommandContext(timeoutCTX, "g++", "/home/ganxue/ganxue-runcode/cpp/main.cpp", "-o", "/home/ganxue/ganxue-runcode/cpp/main")

		// 设置编译的输出捕获
		var buildOutput bytes.Buffer
		buildCmd.Stdout = &buildOutput
		buildCmd.Stderr = &buildOutput

		// 执行编译
		if err := buildCmd.Run(); err != nil {
			// 编译失败，记录错误信息
			log.Printf("编译失败: %v, 输出: %s", err, buildOutput.String())
			// 处理编译失败逻辑
			output, err := json.Marshal(OutputData{
				Status: 2,
				Output: buildOutput.String(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
				continue
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)

				continue
			}
			continue
		}

		// 执行编译后的程序
		cmd := exec.CommandContext(timeoutCTX, "/home/ganxue/ganxue-runcode/cpp/main")
		runCode(cmd, id[1], input, timeoutCTX)
		cancel()
	}
}

func runCCode() {
	for {
		id, err := db.BLPop(context.Background(), 0*time.Second, "runcode-c").Result()
		if err != nil {
			fmt.Println("BLPop失败", err)
			continue
		}
		userCode, input, err := getCodeAndInput(id[1])
		if err != nil {
			log.Println("getCodeAndInput失败, ", err)
			output, err := json.Marshal(OutputData{
				Status: 1,
				Output: err.Error(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}

		// 将用户代码写入文件中
		if err := os.WriteFile("/home/ganxue/ganxue-runcode/c/main.c", []byte(userCode), 0644); err != nil {
			log.Println("写入文件失败", err)
			output, err := json.Marshal(OutputData{
				Status: 1,
				Output: err.Error(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
				continue
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}
		timeoutCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// 编译代码
		buildCmd := exec.CommandContext(timeoutCTX, "gcc", "/home/ganxue/ganxue-runcode/c/main.c", "-o", "/home/ganxue/ganxue-runcode/c/main")

		// 设置编译的输出捕获
		var buildOutput bytes.Buffer
		buildCmd.Stdout = &buildOutput
		buildCmd.Stderr = &buildOutput

		// 执行编译
		if err := buildCmd.Run(); err != nil {
			// 编译失败，记录错误信息
			log.Printf("编译失败: %v, 输出: %s", err, buildOutput.String())
			// 处理编译失败逻辑
			output, err := json.Marshal(OutputData{
				Status: 2,
				Output: buildOutput.String(),
			})
			if err != nil {
				log.Println("marshal output失败", err)
				continue
			}
			err = db.RPush(context.Background(), "runcode-finish"+id[1], output).Err()
			if err != nil {
				log.Println("RPush失败", err)
				continue
			}
			err = db.Expire(context.Background(), "runcode-finish"+id[1], 15*time.Minute).Err()
			if err != nil {
				log.Println("设置过期时间失败", err)
				continue
			}
			continue
		}

		// 执行编译后的程序
		cmd := exec.CommandContext(timeoutCTX, "/home/ganxue/ganxue-runcode/c/main")
		runCode(cmd, id[1], input, timeoutCTX)
		cancel()
	}
}
