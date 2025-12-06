# Link in Bio Platform

Modern, fast, and customizable link-in-bio platform built with SvelteKit and Golang.

## Tech Stack

### Frontend
- **SvelteKit** - Fast, modern web framework
- **Shadcn-Svelte** - Beautiful, accessible UI components
- **Tailwind CSS** - Utility-first styling
- **TypeScript** - Type safety

### Backend
- **Golang Fiber** - Fast HTTP framework
- **PostgreSQL** - Reliable database
- **SQLC** - Type-safe SQL queries
- **Redis** - Caching layer (optional)

## Project Structure

```
link-in-bio/
├── frontend/          # SvelteKit application
├── backend/           # Golang API server
└── docker-compose.yml # Local development setup
```

## Getting Started

### Prerequisites
- Node.js 18+
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Development

1. **Start PostgreSQL:**
```bash
docker-compose up -d postgres
```

2. **Backend:**
```bash
cd backend
go mod download
go run main.go
```

3. **Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Features

- ✅ Custom profile URLs
- ✅ Drag & drop link management
- ✅ Real-time preview
- ✅ Analytics & click tracking
- ✅ Theme customization
- ✅ Link scheduling
- ✅ QR code generation
- ✅ SEO optimized
- ✅ Mobile responsive

## License

MIT
