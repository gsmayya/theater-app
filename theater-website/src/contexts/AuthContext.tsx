import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User } from '../types/show';
import { getFirebaseApp } from '../lib/firebaseClient';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: (email: string, password?: string) => Promise<void>;
  loginWithGoogle: () => Promise<void>;
  logout: () => void;
  loading: boolean;
  googleAvailable: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [googleAvailable, setGoogleAvailable] = useState(false);
  const [firebaseAuth, setFirebaseAuth] = useState<any>(null);

  useEffect(() => {
    // Check Firebase availability and initialize auth listener
    const initAuth = async () => {
      const app = getFirebaseApp();
      if (app) {
        try {
          const { getAuth, onAuthStateChanged } = await import('firebase/auth');
          const auth = getAuth(app);
          setFirebaseAuth(auth);
          setGoogleAvailable(true);

          // Listen to Firebase auth state changes
          const unsubscribe = onAuthStateChanged(auth, (firebaseUser) => {
            if (firebaseUser) {
              const user: User = {
                id: firebaseUser.uid,
                email: firebaseUser.email || '',
                name: firebaseUser.displayName || firebaseUser.email?.split('@')[0] || 'User',
                phone: firebaseUser.phoneNumber || undefined
              };
              setUser(user);
              localStorage.setItem('theater_user', JSON.stringify(user));
            } else if (!user) {
              // Only clear if we don't have a mock user logged in
              const savedUser = localStorage.getItem('theater_user');
              if (savedUser) {
                setUser(JSON.parse(savedUser));
              }
            }
            setLoading(false);
          });

          return unsubscribe;
        } catch (error) {
          console.warn('Firebase Auth initialization failed:', error);
          setGoogleAvailable(false);
        }
      } else {
        // Firebase not configured, fall back to localStorage check
        const savedUser = localStorage.getItem('theater_user');
        if (savedUser) {
          setUser(JSON.parse(savedUser));
        }
        setLoading(false);
      }
    };

    initAuth();
  }, []);

  const login = async (email: string, password?: string) => {
    // Mock authentication - in a real app, this would call your API
    setLoading(true);
    
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Mock user data
    const mockUser: User = {
      id: 'mock_' + Math.random().toString(36).substr(2, 9),
      email: email,
      name: email.split('@')[0].replace(/[^a-zA-Z]/g, '').toLowerCase(),
      phone: '+1234567890'
    };
    
    setUser(mockUser);
    localStorage.setItem('theater_user', JSON.stringify(mockUser));
    setLoading(false);
  };

  const loginWithGoogle = async () => {
    if (!firebaseAuth || !googleAvailable) {
      throw new Error('Google Sign-In not available');
    }

    setLoading(true);
    try {
      const { GoogleAuthProvider, signInWithPopup } = await import('firebase/auth');
      const provider = new GoogleAuthProvider();
      await signInWithPopup(firebaseAuth, provider);
      // User will be set via onAuthStateChanged listener
    } catch (error) {
      console.error('Google Sign-In failed:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    if (firebaseAuth && user?.id && !user.id.startsWith('mock_')) {
      // Firebase user - sign out from Firebase
      try {
        const { signOut } = await import('firebase/auth');
        await signOut(firebaseAuth);
      } catch (error) {
        console.error('Firebase sign out failed:', error);
      }
    }
    
    // Always clear local state
    setUser(null);
    localStorage.removeItem('theater_user');
  };

  const value = {
    user,
    isAuthenticated: !!user,
    login,
    loginWithGoogle,
    logout,
    loading,
    googleAvailable
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
