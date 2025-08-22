#!/bin/bash
mkdir -p ./test_cache

mkdir -p ./test_cache/db

NOWPWD=$(pwd)

echo "PWD: ${NOWPWD}"

cat <<EOF > "${NOWPWD}/test_cache/db/config.toml"
[mysql]
host = "127.0.0.1"           # Mysql地址
port = 3306                  # Mysql端口
maxconn = 50                 # 最大连接数
minconn = 10                 # 最小连接数
user = "root"                # 用户名，首次建议使用root，之后可以创建新用户
password = "wMTs5aXwfjndimtT"  # 密码

[redis]
host = "127.0.0.1"           # Redis地址
port = 6379                  # Redis端口
dbname = 0                   # 数据库编号，通常无需改动
password = "wMTs5aXwfjndimtT"  # 密码
cachetime = 3600             # 缓存时间，单位秒

[grpc]
host = "127.0.0.1" # GRPC地址
port = 50051       # GRPC监听端口
log = true         # 启用日志，调试功能，上线建议关闭
EOF

cp ../bin/StealthIMDB ./test_cache/db/StealthIMDB
chmod +x ./test_cache/db/StealthIMDB

echo "Start DB"
cd ${NOWPWD}/test_cache/db && ./StealthIMDB --config=${NOWPWD}/test_cache/db/config.toml > ${NOWPWD}/test_cache/db.log 2>&1 &

sleep 20s
echo "Start Test"
echo "::group::Test Log"
(pytest test_db_gateway.py -v; echo "$?">${NOWPWD}/test_cache/.ret) | tee ${NOWPWD}/test_cache/test.log
RETVAL=$(cat ${NOWPWD}/test_cache/.ret)
echo "::endgroup::"

sleep 3s
echo "Clean"
ps -aux | grep '[S]tealthIM' | awk '{print $2}' | xargs kill -9

sleep 2s
echo "::group::DB Log"
cat ${NOWPWD}/test_cache/db.log
echo "::endgroup::"

if [ "$RETVAL" -ne 0 ]; then
    echo "::error title=Test failed::Test Log: ${NOWPWD}/test_cache/test.log"
    while IFS= read -r line
    do
        echo "::error::$line"
    done < "${NOWPWD}/test_cache/test.log"
fi

rm ${NOWPWD}/test_cache -r

exit ${RETVAL}
