import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react';
import { User } from '@/types/show';
import { getFirebaseApp } from '@/lib/firebaseClient';
import { config, debugLog, errorLog } from '@/lib/config';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: (email: string, password?: string) => Promise<{ success: boolean; error?: string }>;
  loginWithGoogle: () => Promise<{ success: boolean; error?: string }>;
  register: (email: string, password: string, name: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  loading: boolean;
  googleAvailable: boolean;
  updateUser: (userData: Partial<User>) => void;
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
    const initAuth = () => {
      const app = getFirebaseApp();
      if (app) {
        import('firebase/auth').then(({ getAuth, onAuthStateChanged }) => {
          try {
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
                  ...(firebaseUser.phoneNumber && { phone: firebaseUser.phoneNumber })
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
            setLoading(false);
            return undefined;
          }
        }).catch((error) => {
          console.warn('Firebase Auth import failed:', error);
          setGoogleAvailable(false);
          setLoading(false);
          return undefined;
        });
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

  const login = useCallback(async (email: string, password?: string) => {
    try {
      setLoading(true);
      
      if (config.enableMockData) {
        // Mock authentication
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        const mockUser: User = {
          id: `user_${Date.now()}`,
          email,
          name: email.split('@')[0]?.replace(/[^a-zA-Z]/g, '').toLowerCase() || 'User',
          phone: '+1234567890',
          provider: 'local',
        };
        
        setUser(mockUser);
        localStorage.setItem('theater_user', JSON.stringify(mockUser));
        
        debugLog('Mock login successful', mockUser);
        return { success: true };
      }

      // Real authentication would go here
      // const response = await apiService.authenticate({ email, password });
      // setUser(response.user);
      // localStorage.setItem('theater_user', JSON.stringify(response.user));
      
      return { success: false, error: 'Backend authentication not implemented' };
    } catch (error) {
      errorLog('Login failed', error as Error);
      return { success: false, error: 'Login failed. Please try again.' };
    } finally {
      setLoading(false);
    }
  }, []);

  const loginWithGoogle = useCallback(async () => {
    try {
      setLoading(true);
      
      if (!config.googleClientId || !firebaseAuth || !googleAvailable) {
        return { success: false, error: 'Google Sign-In not available' };
      }

      const { GoogleAuthProvider, signInWithPopup } = await import('firebase/auth');
      const provider = new GoogleAuthProvider();
      await signInWithPopup(firebaseAuth, provider);
      // User will be set via onAuthStateChanged listener
      
      debugLog('Google login successful');
      return { success: true };
    } catch (error) {
      errorLog('Google Sign-In failed', error as Error);
      return { success: false, error: 'Google Sign-In failed. Please try again.' };
    } finally {
      setLoading(false);
    }
  }, [firebaseAuth, googleAvailable]);

  const register = useCallback(async (email: string, password: string, name: string) => {
    try {
      setLoading(true);
      
      if (config.enableMockData) {
        // Mock registration
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        const mockUser: User = {
          id: `user_${Date.now()}`,
          email,
          name,
          provider: 'local',
        };
        
        setUser(mockUser);
        localStorage.setItem('theater_user', JSON.stringify(mockUser));
        
        debugLog('Mock registration successful', mockUser);
        return { success: true };
      }

      // Real registration would go here
      // const response = await apiService.register({ email, password, name });
      // setUser(response.user);
      // localStorage.setItem('theater_user', JSON.stringify(response.user));
      
      return { success: false, error: 'Backend registration not implemented' };
    } catch (error) {
      errorLog('Registration failed', error as Error);
      return { success: false, error: 'Registration failed. Please try again.' };
    } finally {
      setLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      if (firebaseAuth && user?.id && !user.id.startsWith('mock_')) {
        // Firebase user - sign out from Firebase
        const { signOut } = await import('firebase/auth');
        await signOut(firebaseAuth);
      }
    } catch (error) {
      errorLog('Firebase sign out failed', error as Error);
    } finally {
      // Always clear local state
      setUser(null);
      localStorage.removeItem('theater_user');
      debugLog('User logged out');
    }
  }, [firebaseAuth, user]);

  const updateUser = useCallback((userData: Partial<User>) => {
    if (user) {
      const updatedUser = { ...user, ...userData };
      setUser(updatedUser);
      localStorage.setItem('theater_user', JSON.stringify(updatedUser));
      debugLog('User data updated', updatedUser);
    }
  }, [user]);

  const value = {
    user,
    isAuthenticated: !!user,
    login,
    loginWithGoogle,
    register,
    logout,
    loading,
    googleAvailable,
    updateUser
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
