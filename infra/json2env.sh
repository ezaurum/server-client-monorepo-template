#!/usr/bin/env bash

# 사용법 안내 함수
usage() {
  echo "Usage: $0 <input_json_file> [output_env_file]"
  echo "  <input_json_file>  : 입력 JSON 파일 경로"
  echo "  [output_env_file] : 출력 .env 파일 경로 (기본값: ./.env)"
  exit 1
}

# 인자 개수 확인
if [[ $# -lt 1 ]] || [[ $# -gt 2 ]]; then
  usage
fi

input_json_file="$1"
output_env_file="${2:-./.env}" # 두 번째 인자가 없으면 ./.env 사용

# 입력 JSON 파일 존재 확인
if [[ ! -f "$input_json_file" ]]; then
  echo -u2 "Error: 입력 JSON 파일 '$input_json_file'을(를) 찾을 수 없습니다."
  exit 1
fi

# jq 설치 확인
if ! command -v jq &> /dev/null; then
  echo -u2 "Error: 'jq'가 설치되어 있지 않습니다. 이 스크립트를 실행하려면 jq를 설치해주세요."
  echo -u2 "Installation: https://stedolan.github.io/jq/download/"
  exit 1
fi

# 입력 파일이 유효한 JSON 객체인지 확인
if ! jq -e 'type == "object"' "$input_json_file" > /dev/null 2>&1; then
  echo -u2 "Error: 입력 파일 '$input_json_file'은(는) 유효한 JSON 객체가 아닙니다."
  exit 1
fi

echo "Creating '$output_env_file' from '$input_json_file' using jq..."

# jq를 사용하여 JSON을 .env 형식으로 변환하고 파일에 저장
# 1. to_entries[]: JSON 객체를 {"key": "K", "value": V} 형태의 스트림으로 변환
# 2. select(.value != null and (.value | type | IN("string", "number", "boolean"))):
#    값이 null이 아니고, 타입이 string, number, boolean 중 하나인 항목만 선택 (배열, 객체 등 복잡한 타입은 제외)
# 3. .key + "=" + (.value | tostring): "KEY=VALUE" 형태로 문자열 조합
#    (.value | tostring)은 숫자와 불리언을 문자열로 변환 (문자열은 그대로 둠)
# -r 옵션: raw output, 즉 결과 문자열에서 따옴표를 제거하고 순수 텍스트로 출력
# 대문자 환경 변수 이름을 원한다면 .key | ascii_upcase를 사용하여 대문자로 변환할 수 있습니다.
if jq -r 'to_entries[] | select(.value != null and (.value | type | IN("string", "number", "boolean"))) | (.key | ascii_upcase) + "=" + (.value | tostring)' "$input_json_file" > "$output_env_file"; then
  echo "Success: '$output_env_file' 파일이 성공적으로 생성/업데이트되었습니다."
else
  # jq 실행 중 오류가 발생하면 $?는 0이 아닌 값을 가짐
  # (jq는 유효하지 않은 입력에 대해 오류를 stderr로 출력하고 non-zero로 종료될 수 있음)
  echo -u2 "Error: .env 파일을 생성하는 중 오류가 발생했습니다 (jq 실행 오류 또는 파일 쓰기 오류)."
  # 실패 시 생성된 빈 파일이나 부분적으로 작성된 파일을 삭제할 수 있습니다.
  # rm -f "$output_env_file"
  exit 1
fi

exit 0
