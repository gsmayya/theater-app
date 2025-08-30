// src/lib/auth.ts (mock)
// Replaces Firebase auth entirely with a minimal mock to avoid runtime errors.
// The app now uses AuthContext (src/contexts/AuthContext.tsx) instead.

export type User = {
  uid: string;
  email?: string;
  displayName?: string;
};

class GoogleAuthProviderMock {}

type Unsubscribe = () => void;

type Callback = (user: User | null) => void;

let currentUser: User | null = null;
let subscribers: Callback[] = [];

export const auth = {
  GoogleAuthProvider: GoogleAuthProviderMock as unknown as new () => any,
  signInWithPopup: async (_provider: any) => {
    // Simulate a sign-in by creating a mock user
    currentUser = { uid: Math.random().toString(36).slice(2), email: 'mock@example.com', displayName: 'Mock User' };
    subscribers.forEach(cb => cb(currentUser));
    return { user: currentUser } as any;
  },
  onAuthStateChanged: (callback: Callback): Unsubscribe => {
    subscribers.push(callback);
    // Immediately call with current state
    callback(currentUser);
    return () => {
      subscribers = subscribers.filter(cb => cb !== callback);
    };
  },
};
