if [ ! -d /root/twenty_game ];then
        mkdir /root/twenty_game
        chmod -R 777 /root/twenty_game
        cd /root/twenty_game
        git init
        git remote add origin https://github.com/ppeew/twenty_game.git
else
        cd /root/twenty_game
fi

cat ~/.bash_profile
go version

#echo 清理程序
#killall game_web store_web file_web process_web user_web hall_web
#echo 等待srv资源释放
#sleep 2s
#killall game_srv user_srv


echo 删除所有未提交的工作区内容
sleep 2s
git clean -df
git fetch --all
git reset --hard origin/master

git pull origin master 
echo 更新代码成功


cd build
#sh jenkins-build.sh

#重新构建容器
docker-compsoe up --build -d

sleep 4s 
echo 执行远程脚本成功
