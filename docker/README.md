# Docker deployment

Konfigurasi aktif berada di root repo:

- `Dockerfile`: multi-stage build Go dengan cache module/build dan runtime Alpine non-root.
- `docker-compose.yml`: build image lokal, menjalankan service, health check, dan network `web`.
- `docker/deploy.sh`: build lebih dulu, lalu menghentikan container lama dan menjalankan image baru.

## Lokal

Siapkan `.env`, lalu jalankan:

```sh
COMPOSE_PROJECT_NAME=ct-core-local \
CONTAINER_NAME=ct-core-local \
IMAGE_NAME=ct-core:local \
sh docker/deploy.sh
```

Port default adalah `8888`. Ubah dengan `HOST_PORT` dan `CONTAINER_PORT` bila diperlukan.

## GitLab CI/CD

Runner harus berada di server tujuan serta memiliki akses ke Docker daemon dan Docker Compose v2.

Buat variable `APP_ENV_FILE` dengan tipe **File** di GitLab CI/CD Variables. Gunakan environment scope:

- `production` untuk branch `main`.
- `staging` untuk branch `staging`.

Isi file mengikuti `.env.example`. Variable ini disalin sementara sebagai `.env.deploy`, tidak dimasukkan ke image, lalu dihapus setelah job selesai. Registry/Harbor tidak digunakan; image dibangun langsung pada runner dengan cache Docker lokal.
