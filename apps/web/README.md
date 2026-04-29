## Cryplio Web

Frontend for the Cryplio platform.

### Development

Install dependencies and start the app:

```bash
npm install
npm run dev
```

The app runs on [http://localhost:3000](http://localhost:3000).

### Environment

`apps/web/.env.local` should point `API_GATEWAY_URL` at the backend, typically:

```bash
API_GATEWAY_URL=http://localhost:8080
```

### Lint

```bash
npm run lint
```
