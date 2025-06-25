# uptrace not installed :

if [ ! -f "/usr/bin/uptrace" ]; then
  wget https://github.com/uptrace/uptrace/releases/download/v2.0.0-beta.1/uptrace_2.0.0-beta.1_amd64.deb
  sudo dpkg -i uptrace_2.0.0-beta.1_amd64.deb
fi


openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
  -keyout uptrace.key -out uptrace.crt -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost"

uptrace --config=../../uptrace.yml serve
