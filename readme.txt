//// example command

// 打包成apk：
fyne package -os android -appID top.reminisce.pasecret -icon ./assets/logo.png -appVersion 1.0

// 打包成当前系统可执行文件，如exe：
fyne install -icon ./assets/logo.png -appVersion 1.0

// 捆绑资源文件：
fyne bundle -o bundled.go ./assets/font/STXINWEI.TTF
// 添加捆绑
fyne bundle -o bundled.go -append ./assets/font/STXINWEI.TTF

// 运行示例demo：
go run fyne.io/fyne/v2/cmd/fyne_demo@latest