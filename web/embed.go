package web

// @author 肖肖雨歇
// @description 静态资源打包指令，将前端编译物化作 Go 二进制的一部分！

import "embed"

//go:embed dist
var DistFS embed.FS
