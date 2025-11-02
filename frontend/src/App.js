import React, { useState, useEffect } from 'react';
import { onAuthStateChanged } from 'firebase/auth';
import { auth } from './firebase';
import Login from './components/Login';
import RecipeList from './components/RecipeList';
import { createOrUpdateUser } from './utils/userUtils';

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (currentUser) => {
      if (currentUser) {
        // Create or update user document when auth state changes
        await createOrUpdateUser(currentUser);
      }
      setUser(currentUser);
      setLoading(false);
    });

    return () => unsubscribe();
  }, []);

  if (loading) {
    return <div style={{ padding: '20px' }}>Loading...</div>;
  }

  return (
    <div style={{ padding: '20px' }}>
      <h1>RecipeBook</h1>
      <Login user={user} />
      {user && (
        <div style={{ marginTop: '20px' }}>
          <RecipeList />
        </div>
      )}
    </div>
  );
}

export default App;
