# --- build stage ---
FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable

ENV CI=true

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm run build


# --- serve stage ---
FROM node:22-alpine

WORKDIR /app

RUN corepack enable

COPY --from=builder /app/package.json ./
COPY --from=builder /app/pnpm-lock.yaml ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist

EXPOSE 5173

CMD ["pnpm", "vite", "preview", "--host", "0.0.0.0", "--port", "5173"]