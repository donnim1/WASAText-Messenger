WASAText Messenger
WASAText is a web-based messaging application that offers both one-on-one and group chat functionalities. It includes features such as file attachments, message reactions, and message forwarding. The project is built with a Go backend (using SQLite as the datastore) and a Vue.js frontend.

Features
Direct & Group Messaging: Engage in private conversations or create groups.
File Attachments: Send images and GIFs.
Message Reactions: React to messages (e.g., like with a ❤️).
Forwarding & Replying: Forward messages to other chats and reply with context.
Profile Management: Update your username and profile photo.
User Search: Find contacts by username.
Technologies Used
Backend
Go: Core business logic and API server.
SQLite: Embedded database.
httprouter: Lightweight routing.
logrus: Structured logging.
uuid: Unique identifier generation.
Frontend
Vue.js: Reactive UI framework.
Vue Router: SPA navigation.
Axios: HTTP client for API calls.
Bootstrap & Custom CSS: Responsive design.
Setup & Running the Application on Your Local Machine
Backend
Prerequisites:

Go (v1.16 or higher recommended)
SQLite
Build and Run:

Open a terminal in the project root and execute:

bash
Copy
Edit
go run ./cmd/webapi/
By default, the server listens on port 3000. You can adjust settings (such as API host, database file, and timeouts) via command-line flags or by editing the configuration file (default location: /conf/config.yml).

Frontend
Prerequisites:

Node.js (LTS version recommended)
npm or yarn
Build:

To build the distribution, run:

bash
Copy
Edit
yarn run build-prod
Run:

To preview the application, execute:

bash
Copy
Edit
yarn run preview
Docker
The project includes Dockerfiles for both the backend and frontend. Use the following commands to build and run container images:

Build Container Images
Backend:

bash
Copy
Edit
docker build -t wasatext-backend:latest -f Dockerfile.backend .
Frontend:

bash
Copy
Edit
docker build -t wasatext-frontend:latest -f Dockerfile.frontend .
Run Container Images
Backend:

bash
Copy
Edit
docker run -it --rm -p 3000:3000 wasatext-backend:latest
Frontend:

bash
Copy
Edit
docker run -it --rm -p 8080:80 wasatext-frontend:latest
Configuration
The backend configuration is read from command-line flags and an optional YAML file (default: /conf/config.yml). Use these settings to modify the API host, database location, read/write timeouts, and other parameters.
