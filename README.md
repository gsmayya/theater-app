# theater booking system

Testing with docker and trying to write a bare bone ad serving 

Run 

```
docker-compose up --build
```

To test

go 
```
http://localhost:3000/
```

Todo

```bash
# View logs
./docker-manage.sh logs-backend
./docker-manage.sh logs-frontend

# Access database
./docker-manage.sh shell-mysql

# Access Redis
./docker-manage.sh shell-redis

# Complete rebuild
./docker-manage.sh rebuild
```
### Quick start
```bash 
# Make management script executable
chmod +x docker-manage.sh

# Start all services
./docker-manage.sh start

# Check status
./docker-manage.sh status

# Access your applications
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```
