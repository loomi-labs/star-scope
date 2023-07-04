# Stage 1 - Builder stage
FROM node:18-alpine AS build-env

# Install PNPM globally
RUN npm install -g pnpm

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and pnpm-lock.yaml to the working directory
COPY package.json pnpm-lock.yaml ./

# Install dependencies using PNPM
RUN pnpm install --frozen-lockfile

# Copy the entire project to the working directory
COPY . .

# Build the Next.js project
RUN pnpm run build
# Export the Next.js project
RUN pnpm run export

# Stage 2 - Create the run-time image
FROM caddy:2.5.2-alpine
COPY --from=build-env /app/out /var/www/cosmos-login
COPY Caddyfile /etc/caddy/Caddyfile