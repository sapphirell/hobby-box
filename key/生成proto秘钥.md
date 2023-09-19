# 生成私钥
openssl genrsa -des3 -out ca.key 2048 (密码123456)
# 生成证书
openssl req -new -key ca.key -out ca.crt