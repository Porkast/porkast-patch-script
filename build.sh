gf pack manifest internal/packed/data.go -n packed -y
gf build main.go -n guoshao-fm-patch-script -trimpath -a amd64 -s linux,darwin -p ./bin
rm -f internal/packed/data.go 