# SpecForge Frontend

The SpecForge Frontend is a modern, high-performance web application built with React, TypeScript, and Vite. It serves as the primary interface for managing projects, roadmap items, and technical contracts within the SpecForge ecosystem.

## 🚀 Tech Stack

- **Framework**: [React 19](https://react.dev/)
- **Build Tool**: [Vite 7](https://vitejs.dev/)
- **Language**: [TypeScript](https://www.typescriptlang.org/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **UI Components**: [Radix UI](https://www.radix-ui.com/)
- **State Management**: [TanStack Query (React Query)](https://tanstack.com/query/latest)
- **Icons**: [Lucide React](https://lucide.dev/)
- **Visualization**: [React Flow](https://reactflow.dev/) (for dependency graphs)

## 🛠️ Getting Started

### Prerequisites

- **Node.js**: 24+
- **npm**: 10+

### Installation

1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```

### Environment Variables

Create or update the `.env` file in the project root (or `frontend/.env`) with the following:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### Available Scripts

- `npm run dev`: Starts the development server with Hot Module Replacement (HMR).
- `npm run build`: Compiles the application for production, including type checking and OpenApi type generation.
- `npm run generate-types`: Generates TypeScript types from the `specforge-openapi.yaml` definition.
- `npm run type-check`: Runs TypeScript compiler in no-emit mode to check for type errors.
- `npm run lint`: Runs ESLint to identify and report on patterns found in ECMAScript/JavaScript code.
- `npm run preview`: Locally previews the production build.

## 📁 Project Structure

- `src/api`: API client and generated types.
- `src/components`: Reusable UI components (buttons, dialogs, etc.).
- `src/features`: Feature-based modules (projects, roadmap, etc.).
- `src/hooks`: Custom React hooks.
- `src/lib`: Utility functions and shared library configurations.
- `src/store`: Global state management.

## 🎨 Styling and Components

SpecForge uses a custom design system built on top of Tailwind CSS and Radix UI. Components are designed to be accessible, responsive, and aesthetically premium.

---
© 2026 SpecForge Platform. All Rights Reserved.

