cd ../

$root = "github.com/ELQASASystem/server/"
$mainapp = "cmd/class/main.go"
$commitid = git rev-parse --short HEAD

"��ʼ���������..."

go vet cmd/class/main.go
scripts/golangci-lint.exe run

"��ȷ�ϣ����������ʼ����"
[Console]::ReadKey() | Out-Null

"��ʼ����..."

$env:GOOS="linux"
go build -ldflags "-w -X ${root}configs.CommitID=${commitid}" -o build/main ${mainapp}

"������ɣ���������˳�..."
[Console]::ReadKey() | Out-Null
