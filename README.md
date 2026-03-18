# golang-encdec

Jasypt 호환 암복호화 GUI 도구 (PBEWithMD5AndDES)

## 요구사항

- Go 1.24+
- GCC (CGO 빌드용)
- X11 개발 라이브러리 (Linux)

## 설치

### 1. 의존성 설치 (Linux/Ubuntu)

```bash
sudo apt-get install -y gcc libxxf86vm-dev libgl-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev
```

### 2. 빌드

```bash
go mod tidy
go build -ldflags="-s -w" -o golang-encdec .
```

### 3. 실행

```bash
./golang-encdec
```

## Windows 크로스 컴파일 (Linux에서)

### 1. mingw-w64 설치

```bash
sudo apt-get install -y gcc-mingw-w64-x86-64
```

### 2. 빌드

```bash
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
  go build -ldflags="-s -w -H windowsgui" -o golang-encdec.exe .
```

## 바이너리 크기 줄이기

`-ldflags="-s -w"`로 디버그 정보를 제거하면 약 30~40% 감소합니다 (위 빌드 명령에 이미 포함).

UPX를 사용하면 추가로 50~60% 더 줄일 수 있습니다:

```bash
sudo apt-get install -y upx
upx --best golang-encdec.exe
```

## 사용법

1. **암호화/복호화** 라디오 버튼으로 모드 선택
2. **키** 입력
3. **평문** 또는 **암호문** 입력 (`ENC(...)` 래퍼 자동 제거)
4. **Run** 버튼 클릭
5. 결과가 하단에 표시됨 (복사 가능, 실패 시 빨간 글씨로 사유 출력)
