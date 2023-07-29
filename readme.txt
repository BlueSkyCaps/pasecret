//// example command

// 打包成apk：
D:\GoEnv\gopath\bin\fyne.exe package -os android -appID top.example.pasecret -icon ./assets/logo.png -appVersion 1.0


// 打包成当前系统可执行文件，如exe：
D:\GoEnv\gopath\bin\fyne.exe install -icon ./assets/logo.png -appVersion 1.0

// 捆绑资源文件：
D:\GoEnv\gopath\bin\fyne.exe bundle -o bundled.go ./assets/font/STXINWEI.TTF
// 添加捆绑
D:\GoEnv\gopath\bin\fyne.exe bundle -o bundled.go -append ./assets/font/STXINWEI.TTF

// 运行示例demo：
go run fyne.io/fyne/v2/cmd/fyne_demo@latest