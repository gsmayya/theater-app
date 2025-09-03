import React from 'react';
import Navigation from './Navigation';
import Head from 'next/head';

interface LayoutProps {
  children: React.ReactNode;
  title?: string;
  description?: string;
}

const Layout: React.FC<LayoutProps> = ({ 
  children, 
  title = 'Illuminating Windows - Your Gateway to Live Entertainment',
  description = 'Book tickets for the best theater shows in town. Discover upcoming performances, view show schedules, and secure your seats today.'
}) => {
  return (
    <>
      <Head>
        <title>{title}</title>
        <meta name="description" content={description} />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      
      <div className="min-h-screen bg-theater-light flex flex-col">
        <Navigation />
        
        <main className="flex-1">
          {children}
        </main>
        
        <footer className="bg-theater-dark text-white py-8">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              {/* Company Info */}
              <div>
                <div className="flex items-center mb-4">
                  <span className="text-2xl font-bold text-theater-primary">
                    🎭 Illuminating Windows
                  </span>
                </div>
                <p className="text-gray-300 mb-4">
                  Your premier destination for live theater experiences. 
                  Discover amazing shows and book your tickets with ease.
                </p>
              </div>
              
              {/* Quick Links */}
              <div>
                <h3 className="text-lg font-semibold mb-4 text-theater-primary">
                  Quick Links
                </h3>
                <ul className="space-y-2 text-gray-300">
                  <li>
                    <a href="/" className="hover:text-theater-primary transition-colors">
                      Current Shows
                    </a>
                  </li>
                  <li>
                    <a href="/booking" className="hover:text-theater-primary transition-colors">
                      Book Tickets
                    </a>
                  </li>
                  <li>
                    <a href="/my-bookings" className="hover:text-theater-primary transition-colors">
                      My Bookings
                    </a>
                  </li>
                </ul>
              </div>
              
              {/* Contact Info */}
              <div>
                <h3 className="text-lg font-semibold mb-4 text-theater-primary">
                  Contact Us
                </h3>
                <div className="space-y-2 text-gray-300">
                  <p>📧 iw@gmail.com</p>
                  <p>📞 +91 99999 99999</p>
                  <p>📍 123 Abc Road, Bangalore, Karnataka</p>
                </div>
              </div>
            </div>
            
            <div className="border-t border-gray-700 mt-8 pt-6 text-center text-gray-400">
              <p>&copy; 2025 IW. All rights reserved.</p>
            </div>
          </div>
        </footer>
      </div>
    </>
  );
};

export default Layout;
