# Cliq Frontend

Cliq Frontend is the React client for the Cliq short link application. It lets users register, log in, create short links, manage dashboard links, copy short URLs, and access protected pages using the backend's HttpOnly cookie authentication flow.

> This repository contains only the frontend application. The Go backend should live in a separate repository and be configured through environment variables.

---

## Features

- Landing page with short link form
- Register and login pages
- React Router data-mode loaders and actions
- Guest-only auth routes
- Protected dashboard routes
- Cookie-based authentication through backend HttpOnly cookie
- No JWT storage in `localStorage`
- Dashboard link list with pagination data
- Create short link page
- Custom slug validation
- Reserved slug validation
- Delete link action
- Profile page
- SweetAlert2 feedback
- Tailwind CSS and daisyUI styling
- Lucide icons

---

## Tech Stack

- React `19`
- Vite `8`
- React Router `8`
- Tailwind CSS `4`
- daisyUI `5`
- SweetAlert2
- Lucide React
- react-qr-code
- ESLint

---

## Project Structure

```txt
cliq-frontend/
├── public/
│   ├── favicon.svg
│   └── icons.svg
├── src/
│   ├── assets/
│   ├── components/
│   │   ├── LinkCard.jsx
│   │   ├── footer.jsx
│   │   └── header.jsx
│   ├── layouts/
│   │   ├── AppLayout.jsx
│   │   └── AuthLayout.jsx
│   ├── pages/
│   │   ├── CreateLinkPage.jsx
│   │   ├── DashboardPage.jsx
│   │   ├── LandingPage.jsx
│   │   ├── LoginPage.jsx
│   │   ├── ProfilePage.jsx
│   │   └── RegisterPage.jsx
│   ├── utils/
│   │   ├── api.js
│   │   ├── auth.action.js
│   │   ├── auth.js
│   │   ├── link.action.js
│   │   ├── link.loader.js
│   │   └── sweetAlert.js
│   ├── index.css
│   └── main.jsx
├── eslint.config.js
├── index.html
├── package.json
├── package-lock.json
└── vite.config.js
```

---

## Requirements

- Node.js
- npm
- Running Cliq backend API

---

## Environment Variables

Create a `.env` file:

```bash
touch .env
```

Example local development env:

```env
VITE_API_URL=http://localhost:8080
VITE_SHORT_URL_BASE=http://localhost:8080
```

Variables:

| Variable | Description |
| --- | --- |
| `VITE_API_URL` | Base URL of the Cliq backend API |
| `VITE_SHORT_URL_BASE` | Base URL used when displaying generated short links |

If the backend runs at `http://localhost:8080`, keep both values as `http://localhost:8080`.

---

## Installation

Install dependencies:

```bash
npm install
```

Run the development server:

```bash
npm run dev
```

Frontend URL:

```txt
http://localhost:5173
```

Build for production:

```bash
npm run build
```

Preview production build:

```bash
npm run preview
```

Run lint:

```bash
npm run lint
```

---

## Available Scripts

| Command | Description |
| --- | --- |
| `npm run dev` | Start Vite development server |
| `npm run build` | Build production assets |
| `npm run preview` | Preview production build locally |
| `npm run lint` | Run ESLint |

---

## Routes

```txt
/                            Landing page
/auth/register               Register page
/auth/login                  Login page
/dashboard                   User dashboard
/dashboard/create            Create short link page
/dashboard/profile           Profile page
/dashboard/links/:id/delete  Delete link action
/logout                      Logout action
```

Route behavior:

- `/` uses `optionalAuthLoader`.
- `/auth/register` and `/auth/login` use `guestOnlyLoader`.
- `/dashboard/*` uses `requireAuthLoader`.
- `/logout` calls the backend logout endpoint and redirects to login.

---

## Backend API Integration

All API requests go through `src/utils/api.js`.

```js
export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export async function apiRequest(path, options = {}) {
  return fetch(`${API_URL}${path}`, {
    credentials: "include",
    ...options,
  });
}
```

The important part is:

```js
credentials: "include"
```

This allows the browser to send the backend's HttpOnly `access_token` cookie with requests.

---

## Authentication Flow

The frontend does not store JWTs in `localStorage`.

1. User submits the login form.
2. `loginAction` sends credentials to `POST /auth/login`.
3. The backend sets an HttpOnly `access_token` cookie.
4. `getCurrentUser()` calls `GET /auth/me` with `credentials: "include"`.
5. Protected routes use `requireAuthLoader`.
6. Guest-only pages redirect authenticated users away from login/register.
7. Logout calls `POST /auth/logout` and redirects to `/auth/login`.

Because the cookie is HttpOnly, JavaScript cannot read it directly. That is expected and safer than storing tokens manually in browser storage.

---

## API Endpoints Used

| Frontend Feature | Backend Endpoint |
| --- | --- |
| Register | `POST /auth/register` |
| Login | `POST /auth/login` |
| Current user | `GET /auth/me` |
| Logout | `POST /auth/logout` |
| Create short link | `POST /link/create` |
| Dashboard links | `GET /link/dashboard?page=1&limit=10` |
| Delete link | `DELETE /link/:id` |

The public short link redirect is handled by the backend:

```txt
GET /:slug
```

---

## Slug Validation

The frontend validates custom slugs before sending them to the backend.

Rules:

- Optional field
- Must be 3-50 characters when provided
- Can contain only letters, numbers, and hyphens
- Spaces are converted to hyphens
- Reserved slugs are blocked

Reserved slugs:

```txt
api, login, register, dashboard, auth, link, profile, swagger, img
```

The backend should still validate slugs too. Frontend validation is for user experience, not security.

---

## CORS and Cookie Requirements

Because authentication uses cookies, the backend must allow the frontend origin.

Backend `.env` example:

```env
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
COOKIE_SECURE=false
COOKIE_SAMESITE=lax
```

For local development, this is usually enough.

For production:

- Use HTTPS.
- Set `COOKIE_SECURE=true`.
- Use a correct production frontend URL in `ALLOWED_ORIGINS`.
- Use `COOKIE_SAMESITE=none` only if the frontend and backend are on different sites and cross-site cookies are required.

---

## Development Notes

- Do not save the backend access token in `localStorage` or `sessionStorage`.
- Keep `VITE_API_URL` pointed to the backend repo deployment.
- Keep `VITE_SHORT_URL_BASE` pointed to the public backend domain used for redirects.
- React Router actions are used for form submission instead of manually handling all submits with component state.
- Loaders protect routes by calling `/auth/me`, not by checking a locally saved token.

---

## License

This project is licensed under the MIT License.
