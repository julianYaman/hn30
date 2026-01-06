# Contributing to hn30

Thank you for your interest in contributing to hn30! This guide will help you get started.

## 🚀 Getting Started

### Prerequisites

- Node.js 18+
- Go 1.23+ (if working on backend)
- Git

### Frontend Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/julianYaman/hn-news-page.git
   cd hn-news-page/frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and configure:
   - `PRIVATE_API_BASE_URL`: URL to your backend server (default: `http://localhost:8080`)

4. **Start development server**
   ```bash
   npm run dev
   ```

   The frontend will be available at `http://localhost:5173`

### Backend Setup

See the main [README.md](../README.md) for backend setup instructions.

## 📝 Code Style

This project uses:
- **ESLint** for JavaScript/Svelte linting
- **Prettier** for code formatting

Before submitting a PR, run:
```bash
npm run lint      # Check for linting issues
npm run format    # Format code
```

## 🔧 Project Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── components/     # Svelte components
│   │   ├── stores/         # Svelte stores for state management
│   │   └── utils/          # Utility functions
│   └── routes/             # SvelteKit routes & API endpoints
├── static/                 # Static assets
└── build/                  # Production build (gitignored)
```

## 🐛 Reporting Bugs

When reporting bugs, please include:
- Description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Browser/environment details
- Screenshots (if applicable)

## ✨ Suggesting Features

Feature suggestions are welcome! Please:
- Search existing issues first
- Describe the use case
- Explain why this would benefit users

## 🔀 Pull Request Process

1. **Fork the repository** and create your branch from `main`
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes** and test thoroughly
   - Write clear, descriptive commit messages
   - Keep commits focused and atomic
   - Add tests if applicable

3. **Run quality checks**
   ```bash
   npm run lint
   npm run format
   npm run check
   ```

4. **Push to your fork** and submit a pull request
   ```bash
   git push origin feature/my-new-feature
   ```

5. **PR Description** should include:
   - What changes were made
   - Why these changes were necessary
   - Any related issues (use `Fixes #123` or `Closes #456`)

## 🔒 Security

**Never commit:**
- API keys or secrets
- `.env` files (use `.env.example` instead)
- Personal information
- Database files

The `.gitignore` is configured to prevent this, but always double-check.

## 📄 License

By contributing, you agree that your contributions will be licensed under the same license as the project.

## 💬 Questions?

Feel free to open an issue for questions or reach out to the maintainers.

---

Thank you for contributing to hn30! 🎉
