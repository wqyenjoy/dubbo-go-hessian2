# PR提交流程指南

以下是向Apache Dubbo项目提交PR的详细步骤：

## 1. Fork仓库

1. 访问 https://github.com/apache/dubbo-go-hessian2
2. 点击右上角的"Fork"按钮，将仓库Fork到您的个人GitHub账户

## 2. 克隆Fork的仓库

```bash
git clone https://github.com/YOUR_USERNAME/dubbo-go-hessian2.git
cd dubbo-go-hessian2
```

## 3. 添加上游仓库

```bash
git remote add upstream https://github.com/apache/dubbo-go-hessian2.git
```

## 4. 创建分支

```bash
git checkout -b perf/add-decoder-encoder-pool
```

## 5. 提交更改

```bash
git add pool.go list.go map.go decode.go encode_benchmark_test.go medium_object_test.go dubbo_benchmark_test.go
git commit -m "perf: add Decoder/Encoder pool & prealloc"
```

## 6. 推送到您的Fork

```bash
git push origin perf/add-decoder-encoder-pool
```

## 7. 创建Pull Request

1. 访问您的Fork仓库 https://github.com/YOUR_USERNAME/dubbo-go-hessian2
2. 点击"Compare & pull request"按钮
3. 填写PR标题：`perf: add Decoder/Encoder pool & prealloc`
4. 在PR描述中粘贴 `PR_DESCRIPTION.md` 的内容
5. 点击"Create pull request"按钮

## 8. CI检查

创建PR后，Apache Dubbo的CI系统会自动运行测试。请确保所有测试都通过：

- 单元测试
- 竞态检测 (`-race`)
- 多次运行 (`-count=50`)
- 不同平台 (linux-amd64, arm64)

## 9. 代码审查

1. 等待维护者的代码审查
2. 根据反馈进行必要的修改
3. 推送新的更改到同一分支

## 10. 合并PR

一旦PR被批准，维护者将合并您的PR。

## 11. 更新文档

PR合并后，请确保更新以下文档：

1. README.md - 添加关于如何启用池化功能的说明
2. CHANGELOG.md - 添加性能优化的说明

## 注意事项

1. 确保代码符合Go代码规范
2. 确保所有测试都通过
3. 确保PR描述清晰、详细
4. 及时响应维护者的反馈

## 附加资源

- [Apache Dubbo贡献指南](https://dubbo.apache.org/zh-cn/docs/contribution-guidelines/)
- [GitHub Pull Request文档](https://docs.github.com/cn/github/collaborating-with-pull-requests) 