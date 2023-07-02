dcp down -v
dc rmi -f robot:1.0
sudo docker build -t robot:1.0 -f Dockerfile .
sudo docker-compose up -d