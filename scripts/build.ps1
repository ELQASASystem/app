cd ../

$root = "github.com/ELQASASystem/app"
$mainapp = "cmd/class/main.go"
$commitid = git rev-parse --short master

"��ʼ���������..."

go vet cmd/class/main.go
scripts/golangci-lint.exe run

"��ȷ�ϣ����������ʼ����"
[Console]::ReadKey() | Out-Null

"��ʼ����..."

$env:GOOS="linux"
go build -ldflags "-w -X ${root}configs.CommitId=${commitid}" -o build/main ${mainapp}

"������ɣ���������˳�..."
[Console]::ReadKey() | Out-Null
