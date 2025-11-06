import React, { useState } from 'react';
import { signInWithPopup, GoogleAuthProvider, signOut } from 'firebase/auth';
import { auth } from '../firebase';
import RoleDisplay from './RoleDisplay';

function Login({ user, onProfileClick }) {
  const [error, setError] = useState('');

  const handleGoogleSignIn = async () => {
    setError('');
    const provider = new GoogleAuthProvider();

    try {
      await signInWithPopup(auth, provider);
      // User is auto-created in SQLite on first authenticated API request
    } catch (err) {
      setError(err.message);
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut(auth);
    } catch (err) {
      setError(err.message);
    }
  };

  if (user) {
    return (
      <div style={{ padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
        <p>
          Logged in as: {user.email}
          <RoleDisplay />
        </p>
        <div style={{ display: 'flex', gap: '10px' }}>
          <button
            onClick={onProfileClick}
            style={{
              padding: '10px 20px',
              cursor: 'pointer',
              backgroundColor: 'white',
              border: '1px solid #007bff',
              color: '#007bff',
              borderRadius: '4px'
            }}
          >
            Profile
          </button>
          <button
            onClick={handleSignOut}
            style={{
              padding: '10px 20px',
              cursor: 'pointer',
              backgroundColor: '#dc3545',
              border: 'none',
              color: 'white',
              borderRadius: '4px'
            }}
          >
            Sign Out
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={{ padding: '20px', border: '1px solid #ccc', borderRadius: '8px', maxWidth: '400px' }}>
      <h2>Welcome to RecipeBook</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <button
        onClick={handleGoogleSignIn}
        style={{
          width: '100%',
          padding: '12px',
          fontSize: '16px',
          cursor: 'pointer',
          backgroundColor: '#4285f4',
          color: 'white',
          border: 'none',
          borderRadius: '4px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: '10px'
        }}
      >
        {/* Google "G" logo icon */}
        <svg width="18" height="18" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
          <path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/>
          <path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/>
          <path fill="#FBBC05" d="M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.27-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z"/>
          <path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z"/>
          <path fill="none" d="M0 0h48v48H0z"/>
        </svg>
        Sign in with Google
      </button>
    </div>
  );
}

export default Login;
