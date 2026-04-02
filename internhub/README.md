# InternHub — rewritten frontend

## Что добавлено
- swipe + list + liked + history + dashboard
- glassmorphism modal
- confetti, particles background, skeleton loader
- localStorage для лайков, пропусков, фильтров, темы
- match score по techStack
- push-уведомления о дедлайнах
- PWA manifest + service worker

## Запуск локально
```bash
npm install
npm start
```

## Запуск в Docker
```bash
docker-compose up --build
```

По умолчанию фронтенд ходит в backend по:
`http://host.docker.internal:8082/api/internship/`
