# --- build stage ---
FROM node:22-alpine AS builder
WORKDIR /app

RUN corepack enable

ENV CI=true

COPY package.json ./
RUN pnpm install
COPY . .
RUN pnpm run build

# --- serve stage ---
FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]