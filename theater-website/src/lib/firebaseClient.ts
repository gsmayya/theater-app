import { initializeApp, getApps, FirebaseApp } from 'firebase/app';

// Safely initialize Firebase only in the browser and only if configured
export const getFirebaseApp = (): FirebaseApp | null => {
  if (typeof window === 'undefined') return null; // SSR guard

  const apiKey = process.env.NEXT_PUBLIC_FIREBASE_API_KEY;
  const authDomain = process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN;
  const projectId = process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID;

  if (!apiKey || !authDomain || !projectId) {
    // Not configured; allow app to continue with mock auth
    return null;
  }

  const config = {
    apiKey,
    authDomain,
    projectId,
  };

  if (!getApps().length) {
    return initializeApp(config);
  }
  return getApps()[0] || null;
};

