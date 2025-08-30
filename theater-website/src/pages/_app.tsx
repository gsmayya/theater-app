import '../styles/globals.css';
import type { AppProps } from 'next/app';
import { AuthProvider } from '../contexts/AuthContext';

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <AuthProvider>
      <div className="min-h-screen bg-theater-light">
        <Component {...pageProps} />
      </div>
    </AuthProvider>
  );
}

export default MyApp;
