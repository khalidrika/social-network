# Stage 1: Build the application
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# ⬇️ Copy assets into public/src/assets/randoms to preserve access path
RUN mkdir -p public/src/assets/randoms && \
    cp -r src/assets/randoms/* public/src/assets/randoms/

RUN npm run build

# Stage 2: Serve with nginx
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
