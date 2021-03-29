# export GOOS=linux
# go build -o accountservice-linux-amd64

sudo docker build -t ratnadeepb/accountservice .
sudo docker login
sudo docker push ratnadeepb/accountservice