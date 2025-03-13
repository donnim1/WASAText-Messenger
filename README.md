# WASAText Messenger

WASAText Messenger is a web-based, real-time messaging application offering both one-on-one and group chat functionalities. With advanced features like file attachments, message reactions, message forwarding, and user search, WASAText Messenger is designed for both personal and collaborative communication.

## Features

- **Direct & Group Messaging:** Easily start private chats or create group conversations for teams, friends, or communities.
- **File Attachments:** Effortlessly share images, GIFs, and files during chats.
- **Message Reactions:** Express yourself with reactions (e.g., a ❤️) on messages.
- **Forwarding & Replying:** Quickly forward messages to other conversations and reply with message context.
- **Profile Management:** Update your username and profile photo at any time.
- **User Search:** Find contacts by username for faster connectivity.

## Technologies Used

### Backend
- **Go:** Implements the core business logic and API server.
- **SQLite:** Serves as the embedded datastore, ensuring lightweight and efficient data management.
- **[httprouter](https://github.com/julienschmidt/httprouter):** Provides fast and minimal routing for API endpoints.
- **[logrus](https://github.com/sirupsen/logrus):** Enables structured, leveled logging for debugging and monitoring.
- **[uuid](https://github.com/gofrs/uuid):** Generates unique identifiers for users, messages, and conversations.

### Frontend
- **Vue.js:** Provides a reactive and component-based UI for a smooth user experience.
- **Vue Router:** Manages single-page application (SPA) navigation.
- **Axios:** Serves as the HTTP client for seamless API integration.
- **Bootstrap & Custom CSS:** Ensure a responsive, clean, and modern design.

## Setup & Running the Application Locally

### Backend

1. **Prerequisites:**
   - [Go](https://golang.org/dl/) (v1.16 or higher recommended)
   - SQLite (included with most OS installations)

2. **Build and Run:**

   Open a terminal in the project root and execute:

   ```bash
   go run ./cmd/webapi/
   ```

   By default, the API server listens on port `3000`. Settings such as the API host, database file, and timeouts can be adjusted via command-line flags or by editing the configuration file located at `/conf/config.yml`.

### Frontend

1. **Prerequisites:**
   - [Node.js](https://nodejs.org/) (LTS version recommended)
   - npm or yarn (package manager)

2. **Build:**

   Create the production build with:

   ```bash
   yarn run build-prod
   ```

3. **Run:**

   Launch the development preview using:

   ```bash
   yarn run preview
   ```

## Docker

Dockerfiles are provided for both the backend and frontend, allowing you to build container images for deployment.

### Build Container Images

- **Backend:**

  ```bash
  docker build -t wasatext-backend:latest -f Dockerfile.backend .
  ```

- **Frontend:**

  ```bash
  docker build -t wasatext-frontend:latest -f Dockerfile.frontend .
  ```

### Run Container Images

- **Backend:**

  ```bash
  docker run -it --rm -p 3000:3000 wasatext-backend:latest
  ```

- **Frontend:**

  ```bash
  docker run -it --rm -p 8080:80 wasatext-frontend:latest
  ```

## Configuration

The application reads backend configuration settings from command-line flags and an optional YAML configuration file (default location: `/conf/config.yml`). Adjust these settings to change the API host, database file location, timeout values, and other parameters.

---

Happy chatting with WASAText Messenger!