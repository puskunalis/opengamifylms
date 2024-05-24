# OpenGamifyLMS

OpenGamifyLMS is an open-source system designed to gamify the employee learning management process. It provides a comprehensive set of features to create, manage, and deliver engaging learning experiences for employees in startups and other organizations.

## Features

- Course, module, and submodule creation
- Leaderboards to foster competition and motivation
- Progress tracking to monitor learner advancement
- Badges and challenges to incentivize learner engagement
- Quizzes to assess learner knowledge and understanding

## Technologies Used

- Backend: Go
- Frontend: React
- Database: PostgreSQL
- Object Storage: MinIO

## Prerequisites

To set up and run OpenGamifyLMS locally, ensure you have the following dependencies installed:

- Go
- npm
- Node.js
- Docker
- Docker Compose
- Kind (Kubernetes in Docker)
- Helm

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/puskunalis/opengamifylms.git
   ```

1. Navigate to the project directory:
   ```
   cd opengamifylms
   ```

1. Install dependencies:
   ```
   npm install --prefix ./frontend
   ```

1. Start the application in two separate terminals:
   ```
   make
   ```

   ```
   make frontend
   ```

1. Access OpenGamifyLMS through your web browser at `http://localhost:3001`.

## Deployment

To deploy OpenGamifyLMS to a production environment, follow these steps:

1. Update the `values.yaml` file in the `helm` directory to match your deployment requirements.

1. Run the following command to install the Helm chart in the current cluster:
   ```
   helm install opengamifylms ./helm -n opengamifylms --create-namespace
   ```

1. OpenGamifyLMS will be deployed using the specified configuration.

## Running Tests

End-to-end tests are executed using Cypress, which gets run using Docker Compose. You can run it locally with this command:
   ```
   make test-e2e
   ```

## Contributing

Contributing guidelines are available in [CONTRIBUTING.md](CONTRIBUTING.md).

## Known Limitations / TODO

- The badge system requires the badge award rules to be described in Go code in the Backend, since badges can be awarded for different actions, it is not easy to develop a way to configure the badge award rules without having to write any programming code. There is a need for an easily configurable way of awarding badges, so that badges can be awarded for different actions.
- The challenge system is only partially implemented in the Frontend and does not give experience points for completing them. Similarly to badge awarding, creating new challenges without writing code requires extensive configuration options. The challenge system should also decide which challenges to assign to which users, this way implementing an adaptive experience.
- The new user registration system is incompletely implemented. A deeper analysis of the needs of organisations may lead to the consideration of an authentication system based on OAuth 2.0 that integrates with the centralised identity and access management solutions used by the organisation.
- Some API paths are not protected. Some API paths are partially protected using the Middleware design template, which can be adapted to protect the rest of the API paths.
- The instructor user type is currently no different from a regular user, with the instructor role being performed by users of the administrator type. The instructor user type should be separate from the administrator user type and only have privileges to create courses and course materials, and administrators should have an administrator section where they can perform administrative actions -- e.g. modify the details of other users.
- PostgreSQL and MinIO credentials should be stored as Kubernetes Secret objects and not configured via Helm values.
- Course illustrations are generated randomly for each course using Lorem Picsum, photos are used from Unsplash. There is no option for the course instructor to upload his/her preferred illustration. The function developed to upload a photo of a course could store photos in the MinIO object repository, in a similar way to how videos of submodule elements are stored.

## License

This project is licensed under the [MIT License](LICENSE).
