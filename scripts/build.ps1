cd ../

$root = "github.com/ELQASASystem/server/"
$mainapp = "cmd/class/main.go"
$commitid = git rev-parse --short HEAD

"开始检查代码错误..."

go vet cmd/class/main.go
scripts/golangci-lint.exe run

"请确认，按任意键开始编译"
[Console]::ReadKey() | Out-Null

"开始编译..."

$env:GOOS="linux"
go build -ldflags "-w -X ${root}configs.CommitID=${commitid}" -o build/main ${mainapp}

"编译完成，按任意键退出..."
[Console]::ReadKey() | Out-Null
