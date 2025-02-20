FROM node:23-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

ARG MODE=production
ENV MODE=${MODE}

COPY . .
RUN npm run build -- --mode=${MODE}

# Production image
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
