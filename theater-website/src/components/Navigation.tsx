import Link from 'next/link';
import { useState } from 'react';
import { useRouter } from 'next/router';

const Navigation = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const router = useRouter();

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  const isActive = (path: string) => {
    return router.pathname === path;
  };

  return (
    <nav className="bg-theater-dark text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex-shrink-0">
            <Link href="/" className="flex items-center">
              <span className="text-2xl font-bold text-theater-primary">
                🎭 Illuminating Windows
              </span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <div className="hidden md:block">
            <div className="ml-10 flex items-baseline space-x-8">
              <Link
                href="/"
                className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  isActive('/') 
                    ? 'bg-theater-primary text-white' 
                    : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
                }`}
              >
                Current Shows
              </Link>
              <Link
                href="/booking"
                className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  isActive('/booking') 
                    ? 'bg-theater-primary text-white' 
                    : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
                }`}
              >
                Book Tickets
              </Link>
              <Link
                href="/my-bookings"
                  className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    isActive('/my-bookings') 
                      ? 'bg-theater-primary text-white' 
                      : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
                  }`}
                >
                  My Bookings
                </Link>
                <Link
                  href="/admin/dashboard"
                  className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    isActive('/admin/dashboard') 
                      ? 'bg-theater-primary text-white' 
                      : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
                  }`}
                >
                  Admin
                </Link>
            </div>
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden">
            <button
              onClick={toggleMenu}
              className="inline-flex items-center justify-center p-2 rounded-md text-gray-300 hover:text-white hover:bg-theater-primary/20 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
              aria-expanded="false"
            >
              <span className="sr-only">Open main menu</span>
              {/* Hamburger icon */}
              <svg
                className={`${isMenuOpen ? 'hidden' : 'block'} h-6 w-6`}
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 6h16M4 12h16M4 18h16"
                />
              </svg>
              {/* Close icon */}
              <svg
                className={`${isMenuOpen ? 'block' : 'hidden'} h-6 w-6`}
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Navigation Menu */}
      <div className={`${isMenuOpen ? 'block' : 'hidden'} md:hidden`}>
        <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3 bg-theater-dark/95">
          <Link
            href="/"
            className={`block px-3 py-2 rounded-md text-base font-medium transition-colors ${
              isActive('/') 
                ? 'bg-theater-primary text-white' 
                : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
            }`}
            onClick={() => setIsMenuOpen(false)}
          >
            Current Shows
          </Link>
          <Link
            href="/booking"
            className={`block px-3 py-2 rounded-md text-base font-medium transition-colors ${
              isActive('/booking') 
                ? 'bg-theater-primary text-white' 
                : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
            }`}
            onClick={() => setIsMenuOpen(false)}
          >
            Book Tickets
          </Link>
          <Link
            href="/my-bookings"
            className={`block px-3 py-2 rounded-md text-base font-medium transition-colors ${
              isActive('/my-bookings') 
                ? 'bg-theater-primary text-white' 
                : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
            }`}
            onClick={() => setIsMenuOpen(false)}
          >
            My Bookings
          </Link>
          <Link
            href="/admin/dashboard"
            className={`block px-3 py-2 rounded-md text-base font-medium transition-colors ${
              isActive('/admin/dashboard') 
                ? 'bg-theater-primary text-white' 
                : 'text-gray-300 hover:bg-theater-primary/20 hover:text-white'
            }`}
            onClick={() => setIsMenuOpen(false)}
          >
            Admin
          </Link>
        </div>
      </div>
    </nav>
  );
};

export default Navigation;
