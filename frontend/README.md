# Frontend - Personal Web Platform

React application for the personal blog and portfolio.

## Tech Stack

- **Core:** [React 18](https://react.dev/), [TypeScript](https://www.typescriptlang.org/), [Vite](https://vitejs.dev/)
- **UI & Styling:** [Tailwind CSS](https://tailwindcss.com/), [shadcn/ui](https://ui.shadcn.com/) (Radix UI)
- **Data Fetching & State:** [TanStack Query v5](https://tanstack.com/query/latest), [Axios](https://axios-http.com/)
- **Routing:** [React Router v6](https://reactrouter.com/)
- **Editor:** [TipTap](https://tiptap.dev/) (Rich Text Editor)
- **Forms:** React Hook Form, Zod

## Features

### Public

- Home page with profile and bio.
- Blog feed with infinite scroll.
- Article view with markdown support.
- Comments and theming.

### Admin

- Dashboard.
- Profile editor.
- Post management.
- Rich text editor.

## Installation

1.  Install dependencies:

    ```bash
    npm i
    ```

2.  Setup environment:

    Copy `.env.example` to `.env` and set `VITE_BACKEND_URL`.

3.  Run dev server:

    ```bash
    npm run dev
    ```

## Project Structure

- `src/components` - UI components (shadcn), blog blocks, comments, editor.
- `src/contexts` - Context providers (Auth, Theme).
- `src/hooks` - Custom hooks.
- `src/lib` - Utils, API client, types.
- `src/pages` - Pages (Home, Blog, Article, Login, Admin).
