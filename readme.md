# ctf-instance

> Chỉ dành cho ctfd

Để sử dụng, ta cần phải build bằng go:

```bash
go build
```

Khi build xong, ta sẽ được 1 file tên `ctf_instance`, file này sẽ đọc cấu hình trong `config.yml` và sẽ tạo challenge tương ứng. Dưới đây là cây thư mục mẫu với 2 thư mục challenge là `web_eznote` và `web_legacy`:

```
$ pwd
/ctf_instance

$ tree
.
├── config.yml
├── ctf_instance
├── go.mod
├── go.sum
├── instance.sh
├── main.go
├── readme.md
├── sql.db
├── token
├── web_eznote
│   └── docker-compose.yml
└── web_legacy
    └── docker-compose.yml
```

Với cây thư mục đó thì ta có thể cấu hình cho file `config.yml` như sau:

```yaml
title: 'LMAOCTF 2023'

site_token: 'https:/ctf.example.com'

challenges:
  eznote:
    build: ./web_eznote # challenge folder
    timeout: 5
  legacy:
    build: ./web_legacy # challenge folder
    timeout: 5
```

Tên challenge `legacy` và `eznote` sẽ hiện khi người dùng muốn tạo instance và **bắt buộc phải giống với tên service trong `docker-compose.yml`**, còn `build: ./web_eznote` và `build: ./web_legacy` sẽ là đường dẫn tới 2 thư mục challenge tương ứng, không cần thiết phải trùng tên challenge.

Để tạo instance, các challenge đều phải có `docker-compose.yml` và ta sẽ phải chỉnh sửa `container_name` và `ports` bằng cách thêm vào 2 chỉ mục `{{ID}}` và `{{PORT}}`. Dưới dây là file `docker-compose.yml` mẫu

```
services:
  eznote:
    container_name: eznote_{{ID}}
    image: webdevops/php-nginx:8.1-alpine
    ports:
      - {{PORT}}:80
```

Với mỗi service, ta đều phải cấu hình `container_name` để đảm bảo tạo instance không bị trùng tên (nếu trùng sẽ không tạo được). Chẳng hạn với file `docker-compose.yml` khác có 2 service như sau:

```
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: web_legacy_{{ID}}
    networks:
      - no-internet
    restart: always

  nginx:
    image: nginx:latest
    container_name: nginx-proxy_{{ID}}
    ports:
      - {{PORT}}:80
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - app
    networks:
      - no-internet
      - internet
    restart: always

networks:
  internet: {}
  no-internet:
    internal: true
```

Ta thấy có 2 service là `app` và `nginx`, do đó ta đều phải thiết lập `container_name` với chỉ mục `{{ID}}` để mỗi instance tạo ra có một `container_name` duy nhất!

Việc cuối cùng đó là dùng socat dựng server thực thi `ctf_instance` để tạo instance:

```bash
socat tcp-listen:31337,fork,reuseaddr exec:./ctf_instance
```

Hướng dẫn tạo instance: https://gist.github.com/hypnguyen1209/9ee2c1b38c807e86ab3b9045d3a4edca
