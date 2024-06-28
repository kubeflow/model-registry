DOMAIN=$1
mkdir -p certs
# create CA cert
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj "/O=modelregistry Inc./CN=$DOMAIN" -keyout certs/domain.key -out certs/domain.crt
# create rest cert and private key
echo "subjectAltName = DNS:modelregistry-sample-rest.$DOMAIN" > certs/modelregistry-sample-rest.domain.ext
openssl req -out certs/modelregistry-sample-rest.domain.csr -newkey rsa:2048 -nodes -keyout certs/modelregistry-sample-rest.domain.key -subj "/CN=modelregistry-sample-rest/O=modelregistry organization" -addext "subjectAltName = DNS:modelregistry-sample-rest.$DOMAIN"
openssl x509 -req -sha256 -days 365 -CA certs/domain.crt -CAkey certs/domain.key -set_serial 0 -in certs/modelregistry-sample-rest.domain.csr -out certs/modelregistry-sample-rest.domain.crt -extfile certs/modelregistry-sample-rest.domain.ext
# create grpc cert and private key
echo "subjectAltName = DNS:modelregistry-sample-grpc.$DOMAIN" > certs/modelregistry-sample-grpc.domain.ext
openssl req -out certs/modelregistry-sample-grpc.domain.csr -newkey rsa:2048 -nodes -keyout certs/modelregistry-sample-grpc.domain.key -subj "/CN=modelregistry-sample-grpc/O=modelregistry organization" -addext "subjectAltName = DNS:modelregistry-sample-grpc.$DOMAIN"
openssl x509 -req -sha256 -days 365 -CA certs/domain.crt -CAkey certs/domain.key -set_serial 0 -in certs/modelregistry-sample-grpc.domain.csr -out certs/modelregistry-sample-grpc.domain.crt -extfile certs/modelregistry-sample-grpc.domain.ext

# create DB service cert and private key
openssl req -out certs/model-registry-db.csr -newkey rsa:2048 -nodes -keyout certs/model-registry-db.key -subj "/CN=model-registry-db/O=modelregistry organization"
openssl x509 -req -sha256 -days 365 -CA certs/domain.crt -CAkey certs/domain.key -set_serial 0 -in certs/model-registry-db.csr -out certs/model-registry-db.crt
