```bash
# 1. 生成私钥
openssl genrsa -out server.key 2048

# 2. 生成证书  可以全部回车
openssl req -new -x509 -key server.key -out server.crt -days 36500

# 3. 生成证书签名请求 
openssl req -new -key server.key -out server.csr

# 4. 修改openssl.cnf
# 打开copy_extensions = copy
# 打开req_extensions = v3_req
# 在[ v3_req ]中添加 subjectAltName = @alt_names
# 添加标签[alt_names]
# DNS.1 = *.kuangstudy.com

# 5.生成证书私钥
openssl genpkey -algorithm RSA -out test.key

# 6.通过证书生成证书请求文件test.csr
openssl req -new -nodes -key test.key -out test.csr -days 36500 -subj "/C=cn/OU=myorg/O=mycomp/CN=myname" -config ./openssl.cnf -extensions v3_req

# 7.生成SAN证书
openssl x509 -req -days 365 -in test.csr -out test.pem -CA server.crt -CAkey server.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req
```


# 简略流程
```bash 
# 1.生成私钥文件  ECC私钥
openssl ecparam -genkey -name secp384r1 -out server.key

# 2.为证书添加SANs信息

# 3.生成自签名证书
openssl req -nodes -new -x509 -sha256 -days 3650 -config server.cnf -extensions 'req_ext' -key server.key -out server.crt

```

# server.cnf文件
```
[ req ]
default_bits       = 4096
default_md		= sha256
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = BEIJING
localityName                = Locality Name (eg, city)
localityName_default        = BEIJING
organizationName            = Organization Name (eg, company)
organizationName_default    = DEV
commonName                  = Common Name (e.g. server FQDN or YOUR name)
commonName_max              = 64
commonName_default          = liwenzhou.com

[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1   = localhost
DNS.2   = liwenzhou.com
IP      = 127.0.0.1
```