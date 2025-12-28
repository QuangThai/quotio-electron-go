# Shared Resources (Quotio Electron Go)

This directory serves as the bridge between the **Go Backend** and the **Electron Frontend**. It is designed to house any logic, assets, or definitions that are common to both environments, ensuring a "Single Source of Truth" across the application.

## 🎯 Purpose

While the backend and frontend are built using different technologies (Go and TypeScript/React), they share many core concepts that need to stay in sync:

- **Provider Metadata**: Shared definitions for AI providers (Claude, OpenAI, Gemini, etc.), including their identifiers and default configurations.
- **Data Models**: The structure of objects transmitted over the REST API (Backend ↔ Frontend) and Electron IPC (Main ↔ Renderer).
- **Error Handling**: Unified error codes and message structures.
- **Global Constants**: Shared values such as default ports, API versioning, and configuration keys.

## 📂 Structure

To maintain order as the project grows, we recommend following this structure:

| Folder       | Description                                                     |
| :----------- | :-------------------------------------------------------------- |
| `types/`     | Shared JSON schemas or interface definitions (matches Go/TS).   |
| `assets/`    | Shared visual assets like logos, icons, and branding materials. |
| `constants/` | Global constants used across the entire stack.                  |
| `scripts/`   | Utility scripts for deployment, testing, or environment setup.  |

## 🛠 Development Workflow

1.  **Define Once**: When introducing a new feature that affects both sides (e.g., a new Provider type), define it here first.
2.  **Synchronization**: Reference these definitions to update:
    - `backend/internal/providers/` (Go structs)
    - `frontend/src/renderer/types/` (TypeScript interfaces)
3.  **Documentation**: Use this folder to store any architectural diagrams or technical specs that apply to the project as a whole.

---

_Note: This directory ensures consistency and reduces code duplication between the Go and Node.js environments._
