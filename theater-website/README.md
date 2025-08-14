# Ticket Website

This project is a Next.js application that allows users to enter ticket details for events. It includes features for user input, admin login, and ticket management.

## Features

- **User Ticket Form**: Users can enter their mobile number, number of tickets, name, payment mobile number, date and time, and select items from a dropdown list fetched from the backend.
- **Admin Login**: Admins can log in using their Google accounts to manage ticketing.
- **Unique Ticket ID Generation**: Each ticket is assigned a unique ID for tracking purposes.
- **Shareable Ticket Information**: Users can generate shareable text and a QR code containing their ticket information.

## Project Structure

```
ticket-website
├── src
│   ├── components
│   │   ├── AdminLogin.tsx
│   │   ├── TicketForm.tsx
│   │   ├── QRCodeDisplay.tsx
│   │   └── ItemDropdown.tsx
│   ├── pages
│   │   ├── _app.tsx
│   │   ├── index.tsx
│   │   └── admin.tsx
│   ├── lib
│   │   ├── api.ts
│   │   └── auth.ts
│   ├── types
│   │   └── ticket.ts
│   └── utils
│       ├── generateId.ts
│       └── shareText.ts
├── public
│   └── favicon.ico
├── package.json
├── tsconfig.json
└── README.md
```

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   ```
2. Navigate to the project directory:
   ```
   cd ticket-website
   ```
3. Install the dependencies:
   ```
   npm install
   ```

## Usage

To start the development server, run:
```
npm run dev
```
Open your browser and navigate to `http://localhost:3000` to view the application.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.