ssh cloud-vm-45-112 "killall go; export GOPATH=/home/webapps/go; cd ../; cd webapps/go/src/github.com/louiscarteron/WebApps2018/; git pull; screen -d -m go run main.go"
