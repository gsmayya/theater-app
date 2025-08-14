import React from 'react';
import { useEffect } from 'react';
import { auth } from '../lib/auth'; // Import your auth functions
import { useRouter } from 'next/router';

const AdminLogin: React.FC = () => {
  const router = useRouter();

  useEffect(() => {
    const unsubscribe = auth.onAuthStateChanged((user) => {
      if (user) {
        // Redirect to admin dashboard or another page upon successful login
        router.push('/admin/dashboard');
      }
    });

    return () => unsubscribe();
  }, [router]);

  const handleLogin = async () => {
    try {
      await auth.signInWithPopup(new auth.GoogleAuthProvider());
    } catch (error) {
      console.error("Login failed:", error);
    }
  };

  return (
    <div>
      <h1>Admin Login</h1>
      <button onClick={handleLogin}>Login with Google</button>
    </div>
  );
};

export default AdminLogin;