# Etapa 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY . .
RUN npm install -g pnpm && pnpm install && pnpm run build

# Etapa 2: Servir arquivos estáticos
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY --from=builder /app/public /usr/share/nginx/html

# Copia configuração customizada do nginx para habilitar HTTPS
COPY filebrowser.conf /etc/nginx/nginx.conf

# Copia certificados SSL (ajuste os nomes conforme necessário)
COPY ssl/livoltek.com.br.crt /etc/ssl/certs/cert.pem
COPY ssl/livoltek.com.br.key /etc/ssl/private/key.pem

EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]
