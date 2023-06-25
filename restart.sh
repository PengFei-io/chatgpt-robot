pid=$(ps -ef|grep -i chat|awk '{print $2}')
kill -9 $pid
rm -f chatgpt-robot
go build
sudo nohup ./chatgpt-robot > myout.log 2>&1 &
cat -f myout.log