# Listenarr Frontend

React frontend for Listenarr audiobook collection manager.

## Technology Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Material-UI (MUI)** - Component library
- **React Router** - Navigation
- **React Query** - Data fetching and caching
- **Zustand** - State management (when needed)
- **Axios** - HTTP client

## Development

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
npm install
```

### Development Server

```bash
npm run dev
```

Runs on `http://localhost:3000` with proxy to backend API at `http://localhost:8686`

### Build

```bash
npm run build
```

Outputs to `dist/` directory.

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
src/
├── components/        # Reusable components
│   └── Layout/        # Main layout component
├── pages/             # Page components
│   ├── Dashboard/
│   ├── Library/
│   ├── Downloads/
│   ├── Processing/
│   ├── Search/
│   └── Settings/
├── services/          # API services
│   ├── api.ts         # API client
│   ├── library.ts     # Library API
│   └── download.ts    # Download API
├── store/             # State management (Zustand)
├── types/             # TypeScript types
├── utils/             # Utility functions
├── theme/             # Material-UI theme
├── App.tsx            # Main app component
└── main.tsx           # Entry point
```

## Features

- ✅ Responsive layout with sidebar navigation
- ✅ Dark theme (Material-UI)
- ✅ React Router for navigation
- ✅ React Query for data fetching
- ✅ TypeScript for type safety
- ✅ API client with interceptors
- 🚧 Dashboard with statistics
- 🚧 Library management
- 🚧 Download queue
- 🚧 Processing queue
- 🚧 Search functionality
- 🚧 Settings page

## API Integration

The frontend communicates with the backend API at `/api/v1/`. The API client is configured in `src/services/api.ts` and uses Axios with interceptors for authentication and error handling.

## Environment Variables

Create a `.env` file based on `.env.example`:

```env
VITE_API_URL=http://localhost:8686
```

## Code Style

- ESLint for linting
- Prettier for formatting
- TypeScript strict mode enabled

Run linting:
```bash
npm run lint
```

Format code:
```bash
npm run format
```

