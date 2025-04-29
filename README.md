# eShop - Modern Full-Stack E-Commerce Platform

A comprehensive e-commerce platform built with Go (backend) and Next.js (frontend), offering a complete solution for online retail with customer-facing storefront and admin management capabilities.

![eShop Platform](https://place-holder-for-your-screenshot.com)

## 🌟 Features

### Customer Features
- **User Authentication** - Secure login/registration with JWT
- **Product Browsing** - Filter, search, and discover products
- **Shopping Cart** - Add, update, remove items
- **Checkout Process** - Integrated with Stripe for payments
- **Order Management** - Track and manage orders
- **User Profiles** - Manage addresses and preferences

### Admin Features
- **Dashboard** - Sales analytics and business insights
- **Product Management** - Create, edit, delete products
- **Order Processing** - Update order status and track fulfillment
- **User Management** - Handle customer accounts
- **Content Management** - Update site content
- **Inventory Control** - Manage stock levels

## 🏗️ Architecture

### Backend (Go/Golang)
- **RESTful API** built with Gin web framework
- **PostgreSQL** for primary database
- **Redis** for caching and session management
- **PASETO/JWT** for authentication
- **Stripe** integration for payment processing
- **Cloudinary** for image hosting
- **Swagger** for API documentation
- **Asynq** for background job processing
- **Zerolog** for structured logging

### Frontend (Next.js)
- **React 19** with TypeScript
- **Next.js 15** for server-side rendering and routing
- **Tailwind CSS** for styling
- **React Hook Form** for form handling
- **SWR** for data fetching
- **Headless UI & Heroicons** for UI components
- **Framer Motion** for animations
- **React Toastify** for notifications
- **TipTap** for rich text editing
- **React Dropzone** for file uploads

## 📋 Prerequisites

- Go 1.24 or higher
- Node.js 18.x or higher
- PostgreSQL 14+
- Redis 6+
- Docker & Docker Compose (for local development)

## 🚀 Getting Started

### Clone the Repository

```bash
git clone https://github.com/yourusername/eshop.git
cd eshop
```

### Backend Setup

1. Navigate to the server directory:
   ```bash
   cd server
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Configure environment variables:
   ```bash
   cp app.env.example app.env
   # Edit app.env with your configuration
   ```

4. Start the database and Redis with Docker:
   ```bash
   docker-compose up -d postgres redis
   ```

5. Run database migrations:
   ```bash
   make migrate-up
   ```

6. Seed the database (optional):
   ```bash
   make seed
   ```

7. Run the server:
   ```bash
   make run
   ```

### Frontend Setup

1. Navigate to the client directory:
   ```bash
   cd ../client
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Create and configure environment variables:
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your configuration
   ```

4. Start the development server:
   ```bash
   npm run dev
   ```

5. Access the application at http://localhost:3001

## 📚 API Documentation

Once the server is running, Swagger API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

## 🧪 Testing

### Backend Tests

```bash
cd server
make test
```

### Frontend Tests

```bash
cd client
npm run test
```

## 📁 Project Structure

```
eshop/
├── client/                 # Next.js frontend
│   ├── app/                # App router components
│   │   ├── (shop)/         # Customer-facing pages
│   │   └── admin/          # Admin dashboard pages
│   ├── components/         # Reusable React components
│   ├── lib/                # Utility functions and API clients
│   ├── public/             # Static assets
│   └── types/              # TypeScript type definitions
│
├── server/                 # Go backend
│   ├── assets/             # Static assets for backend
│   ├── cmd/                # Application entry points
│   ├── config/             # Configuration management
│   ├── docs/               # Swagger API documentation
│   ├── internal/           # Internal packages
│   │   ├── api/            # API handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Data models
│   │   ├── repository/     # Database access
│   │   └── service/        # Business logic
│   ├── migrations/         # Database migrations
│   ├── pkg/                # Public packages
│   └── seeds/              # Database seed data
│
└── tmp/                    # Temporary files (images, etc.)
```

## 🔧 Available Commands

### Backend (Makefile)

- `make run` - Start the server
- `make build` - Build the application
- `make migrate-up` - Apply database migrations
- `make migrate-down` - Revert database migrations
- `make seed` - Seed the database with initial data
- `make test` - Run tests
- `make swagger` - Generate Swagger documentation

### Frontend (npm/yarn)

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run linter

## 💳 Stripe Integration

### Setting Up Stripe Webhooks Locally

To properly test the payment flow in your development environment, you need to forward Stripe webhook events to your local server. This requires the Stripe CLI.

1. Install the Stripe CLI from [Stripe CLI Installation Guide](https://stripe.com/docs/stripe-cli)

2. Login to your Stripe account:
   ```bash
   stripe login
   ```

3. Start webhook forwarding with all events:
   ```bash
   stripe listen --forward-to localhost:4000/webhook/v1/stripe
   ```

4. Or forward only specific payment-related events:
   ```bash
   stripe listen --events payment_intent.canceled,payment_intent.partially_funded,payment_intent.payment_failed,payment_intent.processing,payment_intent.succeeded,refund.created,refund.failed,refund.updated --forward-to localhost:4000/webhook/v1/stripe
   ```

5. Copy the webhook signing secret provided in the CLI output and add it to your environment variables:
   ```
   STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxx
   ```

The Stripe CLI will forward all webhook events to your local server, allowing you to test the complete payment flow without deploying your application.

## 🔐 Authentication

The application uses JWT (JSON Web Tokens) and PASETO for secure authentication. Access tokens are short-lived with refresh token functionality for enhanced security.

## 🌐 Deployment

### Backend Deployment

1. Build the Docker image:
   ```bash
   cd server
   docker build -t eshop-backend .
   ```

2. Deploy to your server or cloud platform of choice, ensuring environment variables are properly configured.

### Frontend Deployment

1. Build the Next.js application:
   ```bash
   cd client
   npm run build
   ```

2. Deploy to Vercel, Netlify, or your preferred hosting solution:
   ```bash
   npx vercel
   ```

## 📈 Roadmap

- [ ] Mobile application using React Native
- [ ] Internationalization support
- [ ] Advanced analytics dashboard
- [ ] AI-powered product recommendations
- [ ] Multi-vendor marketplace capabilities

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Next.js](https://nextjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Stripe](https://stripe.com/)
- [Cloudinary](https://cloudinary.com/)

---

Developed with ❤️ by [Your Name]