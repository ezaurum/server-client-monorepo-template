# 로컬 개발을 위한 디렉토리

api 서버와 프론트엔드 서버를 proxy 서버를 통해 연결
https환경을 만들어 줌

- docker-compose.yml
- nginx/nginx.conf
- nginx/certs

# 작업 전에 필요한 것

1. https 연결과 nginx에서 backend연결을 위해 `/etc/hosts` 파일에 아래와 같이 추가
    ```
    127.0.0.1 api.dev.local
    127.0.0.1 dev.local
    ```

2. .env.template 파일을 .env로 복사하고, 필요한 환경변수 설정
