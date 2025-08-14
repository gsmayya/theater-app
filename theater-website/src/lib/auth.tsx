// src/lib/auth.ts
import { getAuth, GoogleAuthProvider, signInWithPopup, onAuthStateChanged, User } from 'firebase/auth';
import { initializeApp } from 'firebase/app';

const firebaseConfig = {
  // TODO: Replace with your Firebase config
  apiKey: "YOUR_API_KEY",
  authDomain: "YOUR_AUTH_DOMAIN",
  projectId: "YOUR_PROJECT_ID",
  // ...other config
};

const app = initializeApp(firebaseConfig);
const authInstance = getAuth(app);

export const auth = {
  signInWithPopup: (provider: GoogleAuthProvider) => signInWithPopup(authInstance, provider),
  GoogleAuthProvider: GoogleAuthProvider,
  onAuthStateChanged: (callback: (user: User | null) => void) => onAuthStateChanged(authInstance, callback),
};