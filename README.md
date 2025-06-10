# Go+Vue Monorepo Template

- 모노리포
- [mise](https://mise.jdx.dev/)를 사용하여 관리

# 실행

1. [mise 설치](https://mise.jdx.dev/installing-mise.html)
2. [mise 활성화](https://mise.jdx.dev/getting-started.html#activate-mise) -
   자동으로 현재 디렉토리의 .mise.yaml 파일을 찾아 활성화할 수 있도록 합니다.
3. 프로젝트 설치
   ```bash
   mise install
   ```
4. 개발 환경 실행 -
   로그는 필요한 - `mise dev` 명령어를 사용하여 개발 환경 전체를 한 번에 실행할 수 있습니다.
   [hivemind](https://github.com/DarthSim/hivemind) 를 사용하여 로그를 통합합니다.
   `bash
 mise dev
 ` - 클라이언트 실행
   `bash
     mise dev:client
     ` - 서버 실행
   `bash
     mise dev:server
     ` - 환경을 위한 docker compose 실행
   `bash
     mise infra:up 
     `
