#!/bin/sh

# 检查 goimports 是否安装
if ! command -v goimports >/dev/null 2>&1; then
  echo "goimports 未安装"
  printf "是否安装 goimports? (y/n): "
  read -r REPLY
  if [ "$REPLY" = "y" ] || [ "$REPLY" = "Y" ]; then
    # 检查 go 是否可用
    if ! command -v go >/dev/null 2>&1; then
      echo "未找到 go 命令，请先安装 Go 开发环境后再重试。"
      exit 1
    fi
    echo "正在安装 goimports..."
    if ! go install golang.org/x/tools/cmd/goimports@latest; then
      echo "goimports 安装失败"
      exit 1
    fi
    # 确保当前会话能找到 goimports
    if ! command -v goimports >/dev/null 2>&1; then
      GOBIN_DIR=$(go env GOBIN)
      GOPATH_DIR=$(go env GOPATH)
      # 尝试将典型安装目录加入 PATH
      [ -n "$GOBIN_DIR" ] && PATH="$PATH:$GOBIN_DIR"
      [ -n "$GOPATH_DIR" ] && PATH="$PATH:$GOPATH_DIR/bin"
      export PATH
      if ! command -v goimports >/dev/null 2>&1; then
        echo "goimports 已安装，但当前 PATH 中未找到可执行文件。"
        [ -n "$GOBIN_DIR" ] && echo "请将 $GOBIN_DIR 加入 PATH 后重试。"
        [ -n "$GOPATH_DIR" ] && echo "或将 $GOPATH_DIR/bin 加入 PATH 后重试。"
        exit 1
      fi
    fi
    echo "goimports 安装成功"
  else
    echo "取消安装，脚本退出"
    exit 1
  fi
fi

echo "正在格式化 Go 文件..."

# shellcheck disable=SC2044
for item in $(find . -type f -name '*.go' -not -path './.idea/*'); do
  goimports -l -w "$item"
done

echo "格式化完成"
