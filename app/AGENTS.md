# AGENTS.md

## Project Overview

This is a React 19 + TypeScript + Vite + Tailwind CSS v4 application using shadcn/ui components. **The app currently uses mock server data since the backend is not yet implemented.** Features are semi-mock - API responses come from JSON files but the frontend follows real patterns.

## Commands

```bash
# Development
pnpm run dev           # Start Vite dev server

# Build & Preview
pnpm run build         # TypeScript compile + Vite build
pnpm run preview       # Preview production build

# Linting & Formatting
pnpm run lint          # Run ESLint
```

No test framework is currently configured.

## Path Aliases

Configured in `tsconfig.app.json` and `vite.config.ts`:

- `@/` → `src/`
- `@components` → `src/components`
- `@ui` → `src/components/ui`
- `@lib` → `src/lib`
- `@hooks` → `src/hooks`

## Code Style

### Formatting (Prettier)

Config in `.prettierrc`:
- 2 space indent
- No semicolons
- Double quotes
- Trailing commas (es5)
- Print width 80
- Tailwind CSS plugin enabled (`cn`, `cva` as Tailwind functions)

### Linting (ESLint)

`eslint.config.js` extends:
- `@eslint/js` (recommended)
- `typescript-eslint` (recommended)
- `eslint-plugin-react-hooks` (recommended)
- `eslint-plugin-react-refresh` (recommended)

### TypeScript

- Strict mode enabled via `tsconfig.app.json`
- Use explicit types for function parameters and return types
- Avoid `any` - use `unknown` or proper generics

### Naming Conventions

- **Files**: kebab-case (`my-component.tsx`, `app-error.ts`)
- **Components**: PascalCase (`AppLayout.tsx`, `SiteHeader.tsx`)
- **Hooks**: camelCase with `use` prefix (`useDebounce.ts`, `useMobile.ts`)
- **Utilities**: camelCase (`cn.ts`, `utils.ts`)
- **Types/Interfaces**: PascalCase (`Project`, `GetProjectsParams`)

### Imports

Order (follow existing patterns):
1. External packages (react, react-router, @tanstack/react-query)
2. Internal aliases (@/lib, @/components, @/api)
3. Relative imports (../../)

```typescript
import { useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { Link } from "react-router"
import { cn } from "@/lib/utils"
import { AppError } from "@/lib/app-error"
import { getProjects, PROJECTS_KEY } from "@/api/resources/projects/projects"
import { Button } from "@/components/ui/button"
import type { Project } from "@/api/resources"
import "./styles.css"
```

### Components

- Use shadcn/ui components from `@/components/ui`
- Use `cva` (class-variance-authority) for component variants
- Use `cn()` for conditional class merging
- Follow existing component patterns in `src/components/`

### Error Handling

Use the `AppError` class from `@/lib/app-error`:

```typescript
throw new AppError({
  status: 400,
  message: "Bad Request",
  description?: "Optional detailed description"
})
```

### API Layer

- API functions return data or throw `AppError`
- Response types defined in `@/api/response.ts`
- Resource-specific types in `@/api/resources/*/types.ts`
- Mock data in public folder (`/dummies/*.json`)

## Architecture

```
src/
├── api/
│   ├── resources/       # API endpoints by resource
│   │   └── projects/   # Projects API
│   └── response.ts     # Shared response types
├── components/
│   ├── ui/             # shadcn/ui components
│   ├── *.tsx           # App-specific components
│   └── app-layout.tsx  # Main layout wrapper
├── hooks/              # Custom React hooks
├── lib/
│   ├── utils.ts        # cn() utility
│   └── app-error.ts    # AppError class
├── pages/              # Route pages
├── App.tsx             # Root component
└── main.tsx            # Entry point
```

## Backend Integration Notes

When implementing the backend:
- Update API URLs in `@/api/resources/*/projects.ts` from mock JSON to real endpoints
- Current: `new URL('/dummies/projects.json', window.location.origin)`
- Expected: `new URL('/api/projects',import.meta.env.VITE_API_URL)`
- Keep the error handling pattern with `AppError`
