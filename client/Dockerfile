# This Dockerfile has to be built in the root folder of the project (../)

# Stage 1 - Install dependencies and build the app in a build environment
FROM rust:slim AS build-env

# Install Trunk
RUN cargo install trunk --locked

# Install Rust WASM target
RUN rustup target add wasm32-unknown-unknown

# Install curl and protobuf compiler
RUN apt-get update && \
    apt-get install -y curl protobuf-compiler

# Install Node.js and npm
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs

# Create a new directory for your app
WORKDIR /app

# Copy your project's files into the container
COPY client .
COPY proto ./proto

RUN npm install

# Add node_modules to path
ENV PATH="/app/node_modules/.bin:${PATH}"

ARG GRPC_WEB_ENDPOINT_URL
ARG DEBUG
ARG WEB_APP_URL
ARG TELEGRAM_BOT_NAME
ARG DISCORD_CLIENT_ID
ENV GRPC_WEB_ENDPOINT_URL $GRPC_WEB_ENDPOINT_URL
ENV COSMOS_LOGIN_BUTTON_URL $COSMOS_LOGIN_BUTTON_URL
ENV DEBUG $DEBUG
ENV WEB_APP_URL $WEB_APP_URL
ENV TELEGRAM_BOT_NAME $TELEGRAM_BOT_NAME
ENV DISCORD_CLIENT_ID $DISCORD_CLIENT_ID

# Build your Sycamore project using Trunk
RUN trunk build --release

# Stage 2 - Create the run-time image
FROM caddy:2.5.2-alpine
COPY --from=build-env /app/dist /var/www/star-scope
COPY client/Caddyfile /etc/caddy/Caddyfile