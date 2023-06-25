pid=$(ps -ef|grep -i chat|awk '{print $2}')
sudo kill -9 $pid
sudo rm -f chatgpt-robot
go build
sudo nohup ./chatgpt-robot > myout.log 2>&1 &
tail -f  myout.log