# Dependencies stage
FROM node:18-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy workspace files
COPY ../../package.json ../../pnpm-workspace.yaml ../../pnpm-lock.yaml ./
COPY ../../packages ./packages
COPY package.json ./apps/dashboard/

# Install dependencies
RUN pnpm install --frozen-lockfile

# Builder stage
FROM node:18-alpine AS builder
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy workspace setup
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/packages ./packages
COPY ../../package.json ../../pnpm-workspace.yaml ./
COPY . ./apps/dashboard

WORKDIR /app/apps/dashboard

# Build-time URL config must be provided by CI/CD for production builds.
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_TELEMETRY_DISABLED=1

# Build the app
RUN pnpm build

# Runner stage
FROM node:18-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy built files
COPY --from=builder /app/apps/dashboard/public ./public
COPY --from=builder /app/apps/dashboard/.next/standalone ./
COPY --from=builder /app/apps/dashboard/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000

CMD ["node", "server.js"]
