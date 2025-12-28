# Quotio Frontend

## Project Structure

```
frontend/
├── src/
│   ├── main/              # Electron main process
│   ├── preload/           # Preload scripts (IPC bridge)
│   ├── renderer/          # React UI
│   ├── shared/            # Shared code between processes
│   └── assets/           # Static assets (icons, images)
├── build/                # Build scripts and icons
├── dist/                 # Renderer build output
├── dist-electron/        # Electron (main/preload) build output
├── release/              # Electron installers (.dmg, .exe, etc.)
└── public/               # Public static files
```

## Development

```bash
# Install dependencies
npm install

# Start Vite dev server (port 5173)
npm run dev

# Start Electron with dev server
npm run electron:dev
```

## Build

```bash
# Build renderer only
npm run build:renderer

# Build Electron main/preload only
npm run build:electron

# Build everything and create installers
npm run build

# Build for specific platform
npm run electron:pack
```

## Tech Stack

- **Framework**: React 18
- **Bundler**: Vite 5
- **Styling**: Tailwind CSS 3.4
- **Desktop**: Electron 28
- **State Management**: Zustand
- **Routing**: React Router 6
- **HTTP Client**: Axios
- **Icons**: Lucide React

## Best Practices Followed

✅ Separation of main/preload/renderer processes
✅ Reusable UI components
✅ Centralized API services
✅ State management with Zustand
✅ Cross-platform build configuration
✅ Shared code between processes
✅ Proper folder structure for scalability

