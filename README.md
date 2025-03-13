WASAText Messenger
WASAText is a web-based messaging application offering both one-on-one and group chat functionalities. It supports features such as file attachments, message reactions, and message forwarding. The project comprises a Go backend (utilizing SQLite as the datastore) and a Vue.js frontend.

Features
Direct & Group Messaging: Initiate private conversations or create group chats.
File Attachments: Share images and GIFs.
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
Project Structure
cmd/: Contains all executables; Go programs here should handle tasks like reading options from the CLI/environment.
cmd/healthcheck: Daemon for checking the health of server daemons; useful when the hypervisor doesn't provide HTTP readiness/liveness probes (e.g., Docker engine).
cmd/webapi: Web API server daemon.
demo/: Demo configuration file.
doc/: Documentation (typically, for APIs, this means an OpenAPI file).
service/: Packages implementing project-specific functionalities.
service/api: API server implementation.
service/globaltime: Wrapper package for time.Time (useful in unit testing).
vendor/: Managed by Go, contains a copy of all dependencies.
webui/: Web frontend in Vue.js; includes:
Bootstrap JavaScript framework.
Customized version of the "Bootstrap dashboard" template.
Feather icons as SVG.
Go code for release embedding.
Other project files include:

open-node.sh: Starts a new (temporary) container using the node:20 image for safe and secure web frontend development.
Go Vendoring
This project uses Go Vendoring. After modifying dependencies (go get or go mod tidy), run go mod vendor and commit all files under the vendor/ directory.

For more information about vendoring:

Go Modules: Vendoring
Go Modules 06: Vendoring
Node/YARN Vendoring
This repository uses yarn and a vendoring technique that exploits the "Offline mirror". Commit the files inside the .yarn directory.

How to Set Up a New Project from This Template
Change the Go module path to your module path in go.mod, go.sum, and in *.go files throughout the project.
Rewrite the API documentation in doc/api.yaml.
If no web frontend is expected, remove webui and cmd/webapi/register-webui.go.
Update the top/package comment inside cmd/webapi/main.go to reflect the actual project usage, goal, and general info.
Update the code in the run() function (cmd/webapi/main.go) to connect to databases or external resources.
Write API code inside service/api, and create any further packages inside service/ (or subdirectories).
How to Build
If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable:

bash
Copy
Edit
go build ./cmd/webapi/
If you're using the WebUI and want to embed it into the final executable:

bash
Copy
Edit
./open-node.sh
# Inside the container:
yarn run build-embed
exit
# Outside the container:
go build -tags webui ./cmd/webapi/
How to Run (in Development Mode)
To launch the backend only:

bash
Copy
Edit
go run ./cmd/webapi/
To launch the WebUI, open a new terminal tab and run:

bash
Copy
Edit
./open-node.sh
# Inside the container:
yarn run dev
How to Build for Production / Homework Delivery
bash
Copy
Edit
./open-node.sh
# Inside the container:
yarn run build-prod
For "Web and Software Architecture" students: before committing and pushing your work for grading, please read the section below named "My build works when I use yarn run dev, however there is a Javascript crash in production/grading".

Known Issues
My Build Works When I Use yarn run dev, However There Is a JavaScript Crash in Production/Grading
Some errors in the code may not appear in vite development mode. To preview the code as it will be in production/grading settings, use the following commands:

bash
Copy
Edit
./open-node.sh
# Inside the container:
yarn run build-prod
yarn run preview
License
See LICENSE.
