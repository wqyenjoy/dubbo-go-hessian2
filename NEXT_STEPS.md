# 提交PR的下一步操作

## 1. 推送到您的个人Fork仓库

```bash
# 添加您的个人Fork仓库作为远程仓库（如果尚未添加）
git remote add fork https://github.com/YOUR_USERNAME/dubbo-go-hessian2.git

# 推送到您的Fork仓库
git push fork perf/add-decoder-encoder-pool
```

## 2. 在GitHub上创建PR

1. 访问您的Fork仓库：https://github.com/YOUR_USERNAME/dubbo-go-hessian2
2. 点击"Compare & pull request"按钮
3. 确保PR的目标是`apache/dubbo-go-hessian2`的`master`分支
4. 在PR标题中填写：`perf: add Decoder/Encoder pool & prealloc`
5. 在PR描述中粘贴`PR_DESCRIPTION.md`的内容
6. 点击"Create pull request"按钮

## 3. 联系维护者进行代码审查

1. 在PR中@提及主要维护者，例如@wongoo和@AlexStocks
2. 或者通过邮件列表联系维护者，使用`MAINTAINER_EMAIL.md`中的内容

## 4. 跟进PR审查

1. 及时响应代码审查中的反馈
2. 根据需要进行修改并推送到同一分支
3. 确保所有CI检查都通过

## 5. PR合并后的工作

1. 更新README.md，添加关于如何启用池化功能的说明
2. 更新CHANGELOG.md，添加性能优化的说明
3. 考虑下一步优化方向，例如：
   - 在Go 1.22+中使用arena.New()替代部分池化功能
   - 进一步优化缓冲区管理策略
   - 添加更细粒度的池化控制选项 