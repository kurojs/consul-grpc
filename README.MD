# Consul

## Cài đặt Consul

- Tải xuống Consul tại đây [Consul Download page](https://www.consul.io/downloads.html)

- Sau đó đảm bảo rằng file binary nằm trong một thư mục của đường dẫn $PATH (mv consul /usr/bin/)

- Chạy Consul Agent:

    -Consul Agent có thể được chạy dưới 2 mode: client và server. Mỗi consul datacenter cần phải có ít nhất một server agent để đảm bảo duy trì trạng thái (state) của Consult. Trạng thái này bao gồm thông tin về các server và client khác, các service sẵng sàng cho discovery, và những service nào được phép giao tiếp với nhau. Để đảm bảo tránh mất mát dữ liệu, nên chạy 3 hoặc 5 server Consul.

    - Lệnh chạy consul dưới chế độ server `consul agent -dev -node machine`. (Ở production phải bỏ cờ -dev đi).

    - Các Consul agent client dùng để đăng ký service, health check, forward queries đến server. [Client UI](http://localhost:8500)

## Đăng ký service và Health check, Service Discovery

Để đăng ký service có hai cách: sử dụng config file (.json hoặc .hcl) service definition hoặc gọi đến HTTP API để đăng ký service (Register Service)

### Service definition



### Register service

## 
