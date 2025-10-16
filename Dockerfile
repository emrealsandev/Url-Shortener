# --- 1. Aşama: Build Aşaması ---
# Go'nun kurulu olduğu bir imajı temel alarak başlıyoruz
FROM golang:1.25-alpine AS builder

# Çalışma dizinini ayarlıyoruz
WORKDIR /app

# go.mod ve go.sum dosyalarını kopyalıyoruz. Bu dosyalar değişmediği sürece
# Docker bu katmanı cache'ler ve her seferinde bağımlılıkları indirmez.
COPY go.mod go.sum ./
RUN go mod download

# Tüm proje kaynak kodunu kopyalıyoruz
COPY . .

# Uygulamayı derliyoruz. CGO_ENABLED=0 statik bir binary oluşturur.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# --- 2. Aşama: Final Aşaması ---
# Sadece uygulamayı çalıştırmak için KÜÇÜCÜK bir imaj kullanıyoruz.
# Scratch veya Alpine en iyi seçeneklerdir. Alpine debug için daha kolaydır.
FROM alpine:latest

# Çalışma dizini
WORKDIR /app

# Sadece ve sadece bir önceki aşamada derlediğimiz "main" dosyasını
# builder'dan bu yeni imajın içine kopyalıyoruz.
COPY --from=builder /app/main .

# Uygulama 8080 portunu dinleyecek
EXPOSE 8080

CMD ["/app/main"]