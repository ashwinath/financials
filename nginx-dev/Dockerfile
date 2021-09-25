FROM nginx:alpine

WORKDIR /app
CMD ["nginx", "-c", "/app/nginx.conf", "-g", "daemon off;"]
COPY . .
