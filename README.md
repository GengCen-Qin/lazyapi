# Lazyapi

从 [lazygit](https://github.com/jesseduffield/lazygit) 中受到启发，并且无法忍受现有的Api管理工具的启动速度慢，操作繁琐等问题，希望能在命令行中快速管理请求并进行发送查看

## Feature

1. Api的增删改查

<img src="./pic/1743236384511.png" alt="增加" style="zoom:30%;" />

2. Api请求发送，暂支持Get/Post请求，可配置请求体信息（统一使用Json），并展示响应内容

<img src="./pic/1743236340919.png" alt="增加" style="zoom:30%;" />

3. 数据本地存储，使用SQLite

Mac用户会自动在 ~/Library/Application Support/lazyapi 下生成lazyapi.db文件

## 快捷键

### API 列表视图 (API_LIST)
- `n` - 创建新API
- `e` - 编辑API
- `d` - 删除API
- `r` - 发起请求
- `空格` - 跳转到详情
- `tab` - 切换视图(在API列表和请求记录列表切换)
- `g` - 发起快速GET请求
- `p` - 发起快速POST请求

### API 详情视图 (API_INFO)
- `↑` - 向上翻页
- `↓` - 向下翻页
- `esc` - 返回列表

### 响应信息视图 (RESPOND_INFO)
- `↑` - 向上翻页
- `↓` - 向下翻页
- `esc` - 返回列表

### 请求确认视图 (REQUEST_CONFIRM_VIEW)
- `ctrl-r` - 确认
- `ctrl-q` - 取消

### 记录列表视图 (RECORD_LIST)
- `d` - 删除
- `空格` - 跳转到详情
- `tab` - 切换视图
- `g` - 快速GET请求
- `p` - 快速POST请求

### API新建和编辑视图
- `esc` - 取消

编写json信息时支持 `ctrl-f` 进行信息格式化

## 待改进

1. 底层使用[gocui](https://github.com/jroimartin/gocui) 对中文支持不好，目前我fork了一个版本，并临时解决了中文输入与展示问题。但是对于光标移动和字符删除仍有问题。比如删除上一个字符，则需要按两下删除键，才能把字符完整删除掉，光标移动也是类似的。英文使用无误。后期待修改
2. 暂时仅支持Get和Post请求，后续会加上其他请求方式，并增加Header的设置
3. 需要能通过SwaggerUrl快速导入Api数据，可快速上手
4. 支持可配置

## 使用方式

下载代码后执行：`go build` 得到对应的可执行文件
