#!/bin/sh

echo "正在格式化 Go 文件..."

# shellcheck disable=SC2044
for item in $(find . -type f -name '*.go' -not -path './.idea/*'); do
  goimports -l -w "$item"
done

echo "格式化完成"
