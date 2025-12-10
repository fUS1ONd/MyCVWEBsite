# Build Stage
FROM node:20-alpine as builder
WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

ARG VITE_BACKEND_URL
ARG VITE_GTM_ID
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL
ENV VITE_GTM_ID=$VITE_GTM_ID

RUN npm run build

# Serving Stage (Internal Nginx)
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx-internal.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
