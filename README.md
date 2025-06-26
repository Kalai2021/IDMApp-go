# IDM React App

A React-based frontend application for the Identity Management (IDM) system. This app provides a modern, responsive interface for managing users with authentication and CRUD operations.

## Features

- 🔐 **Authentication**: Secure login with token-based authentication
- 👥 **User Management**: Complete CRUD operations for users
- 🎨 **Modern UI**: Clean, responsive design with Material Design principles
- 📱 **Mobile Friendly**: Responsive design that works on all devices
- 🔒 **Protected Routes**: Automatic redirection for unauthenticated users
- 🚀 **Docker Ready**: Containerized deployment with Docker

## Prerequisites

- Node.js 18+ 
- npm or yarn
- Docker (for containerized deployment)
- Access to the IDM Go backend API

## Quick Start

### Development

1. **Clone and install dependencies:**
   ```bash
   npm install
   ```

2. **Set up environment variables:**
   ```bash
   cp env.example .env
   # Edit .env with your API URL
   ```

3. **Start the development server:**
   ```bash
   npm start
   ```

4. **Open your browser:**
   Navigate to `http://localhost:3000`

### Production with Docker

1. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

2. **Or build manually:**
   ```bash
   docker build -t idm-react-app .
   docker run -p 3000:3000 idm-react-app
   ```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `REACT_APP_API_URL` | Backend API URL | `http://localhost:8090` |
| `REACT_APP_AUTH0_DOMAIN` | Auth0 domain (if using Auth0) | - |
| `REACT_APP_AUTH0_CLIENT_ID` | Auth0 client ID | - |
| `REACT_APP_AUTH0_AUDIENCE` | Auth0 API audience | - |

### API Configuration

The app connects to the IDM Go backend API. Make sure the backend is running and accessible at the configured URL.

## Project Structure

```
src/
├── components/          # React components
│   ├── Dashboard.tsx    # Main dashboard
│   ├── Login.tsx        # Login form
│   ├── UserList.tsx     # User management table
│   ├── UserModal.tsx    # User create/edit modal
│   └── ProtectedRoute.tsx # Authentication wrapper
├── contexts/            # React contexts
│   └── AuthContext.tsx  # Authentication state
├── services/            # API services
│   └── api.ts          # API client
├── types/              # TypeScript interfaces
│   └── index.ts        # Type definitions
└── App.tsx             # Main app component
```

## Authentication

The app uses token-based authentication:

1. **Login**: Users authenticate with email/password
2. **Token Storage**: JWT tokens are stored in localStorage
3. **Auto-refresh**: Tokens are automatically included in API requests
4. **Logout**: Tokens are cleared on logout

## User Management

### Features
- **View Users**: Display all users in a responsive table
- **Create User**: Add new users with validation
- **Edit User**: Update existing user information
- **Delete User**: Remove users with confirmation
- **Status Management**: Toggle user active/inactive status

### User Fields
- Email (unique, required)
- First Name (required)
- Last Name (required)
- Display Name (required)
- Status (active/inactive)

## Development

### Available Scripts

- `npm start` - Start development server
- `npm run build` - Build for production
- `npm test` - Run tests
- `npm run eject` - Eject from Create React App

### Code Style

The project uses:
- TypeScript for type safety
- ESLint for code linting
- Prettier for code formatting
- Functional components with hooks

## Docker Deployment

### Development
```bash
docker-compose up --build
```

### Production
```bash
# Build the image
docker build -t idm-react-app .

# Run the container
docker run -d -p 3000:3000 \
  -e REACT_APP_API_URL=http://your-api-url:8090 \
  idm-react-app
```

### Docker Compose
The `docker-compose.yml` file includes:
- Multi-stage build for optimization
- Nginx for serving static files
- Environment variable configuration
- Network configuration

## API Integration

The app integrates with the IDM Go backend API:

### Endpoints Used
- `POST /api/v1/auth/login` - User authentication
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get specific user
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Error Handling
- Automatic token refresh on 401 errors
- User-friendly error messages
- Loading states for better UX

## Security

- HTTPS enforcement in production
- Secure headers via Nginx
- Token-based authentication
- Input validation and sanitization
- XSS protection

## Troubleshooting

### Common Issues

1. **API Connection Failed**
   - Check if the backend is running
   - Verify `REACT_APP_API_URL` in environment
   - Check network connectivity

2. **Authentication Issues**
   - Clear browser localStorage
   - Check token expiration
   - Verify backend authentication

3. **Docker Issues**
   - Ensure Docker is running
   - Check port conflicts
   - Verify Docker Compose version

### Logs
- Development: Check browser console
- Production: Check Docker logs
- API: Check backend logs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 