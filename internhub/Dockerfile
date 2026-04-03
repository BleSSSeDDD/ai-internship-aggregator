FROM node:20-alpine AS build
WORKDIR /app
COPY package.json ./
COPY vite.config.js ./
RUN npm install
COPY . .
RUN npm run build

FROM alpine:latest
RUN apk add --no-cache python3
COPY --from=build /app/dist /usr/share/html
WORKDIR /usr/share/html
EXPOSE 80
CMD ["python3", "-m", "http.server", "80"]