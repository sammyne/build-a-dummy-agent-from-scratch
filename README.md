# 从 0 开始构建一个普通智能体

演示构建一个能够对话、调用工具查看/编辑文件（简单编写代码等）的智能体。

## 快速开始

```bash
# 根据个人情况调整即可。本详细白嫖硅基流动的免费额度。
export OPENAI_ENDPOINT="https://api.siliconflow.cn/v1/chat/completions"
# 硅基流动提供了多种模型，其中 DeepSeek-V3、Qwen/Qwen3-32B 等的速度可用，DeepSeek-R1 的不大行。
export OPENAI_MODEL=deepseek-ai/DeepSeek-V3
export OPENAI_API_KEY=你的API密钥

go run main.go
# 后续参照 https://mp.weixin.qq.com/s/bxKJxHXq0AYbFnXJXHJRsg 和这个智能体交互即可。
```

## 温馨提示
- 谷歌的 Gemini 声称兼容 OpenAI 的接口经测试不兼容，特别是涉及到工具的的使用
  - 相关 issue：https://github.com/googleapis/google-api-python-client/issues/2570

## 参考文献
- [代码Agent没有护城河？我用Go标准库和DeepSeek证明给你看！](https://mp.weixin.qq.com/s/bxKJxHXq0AYbFnXJXHJRsg)
