# Builder stage
FROM node:18-alpine AS builder

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

# Running stage
FROM node:18-alpine

# Install PNPM globally
RUN npm install -g pnpm

# Set the working directory inside the container
WORKDIR /app

# Copy built files from the builder stage
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./
COPY --from=builder /app/pnpm-lock.yaml ./

# Install only production dependencies using PNPM
RUN pnpm install --frozen-lockfile --production

# Set the environment variable
ENV NODE_ENV production

# Expose the desired port (adjust if necessary)
EXPOSE 3000

# Start the Next.js application
CMD ["pnpm", "start"]