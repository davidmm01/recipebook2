import React, { useState, useEffect, useCallback } from 'react';
import { onAuthStateChanged } from 'firebase/auth';
import { auth } from './firebase';
import Login from './components/Login';
import RecipeList from './components/RecipeList';
import RecipeForm from './components/RecipeForm';
import RecipeDetail from './components/RecipeDetail';
import RecipeFilters from './components/RecipeFilters';
import UserProfile from './components/UserProfile';

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showForm, setShowForm] = useState(false);
  const [selectedRecipeId, setSelectedRecipeId] = useState(null);
  const [filters, setFilters] = useState({});
  const [showProfile, setShowProfile] = useState(false);
  const [recipeType, setRecipeType] = useState('food');

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (currentUser) => {
      // User creation now happens automatically on backend during first authenticated request
      setUser(currentUser);
      setLoading(false);
    });

    return () => unsubscribe();
  }, []);

  const handleRecipeCreated = () => {
    // Trigger recipe list refresh and hide form
    setRefreshTrigger(prev => prev + 1);
    setShowForm(false);
  };

  const handleFilterChange = useCallback((newFilters) => {
    setFilters(newFilters);
  }, []);

  const handleRecipeTypeChange = (type) => {
    setRecipeType(type);
    setFilters({}); // Reset filters when changing type
  };

  if (loading) {
    return <div style={{ padding: '20px' }}>Loading...</div>;
  }

  return (
    <div style={{ padding: '20px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h1 style={{ margin: 0 }}>RecipeBook</h1>
        <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
          {user && (
            <button
              onClick={() => setShowProfile(true)}
              style={{
                padding: '8px 16px',
                fontSize: '14px',
                color: '#007bff',
                backgroundColor: 'white',
                border: '1px solid #007bff',
                borderRadius: '4px',
                cursor: 'pointer',
              }}
            >
              Profile
            </button>
          )}
          <Login user={user} />
        </div>
      </div>

      {showProfile && (
        <UserProfile onClose={() => setShowProfile(false)} />
      )}

      {user && (
        <div style={{ marginTop: '20px' }}>
          {selectedRecipeId ? (
            <RecipeDetail
              recipeId={selectedRecipeId}
              onBack={() => setSelectedRecipeId(null)}
            />
          ) : (
            <>
              {/* Recipe Type Toggle */}
              <div style={{ marginBottom: '20px', display: 'flex', gap: '10px', alignItems: 'center' }}>
                <button
                  onClick={() => handleRecipeTypeChange('food')}
                  style={{
                    padding: '12px 24px',
                    fontSize: '16px',
                    fontWeight: '600',
                    color: recipeType === 'food' ? '#fff' : '#007bff',
                    backgroundColor: recipeType === 'food' ? '#007bff' : '#fff',
                    border: '2px solid #007bff',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    transition: 'all 0.2s'
                  }}
                >
                  Food
                </button>
                <button
                  onClick={() => handleRecipeTypeChange('drink')}
                  style={{
                    padding: '12px 24px',
                    fontSize: '16px',
                    fontWeight: '600',
                    color: recipeType === 'drink' ? '#fff' : '#007bff',
                    backgroundColor: recipeType === 'drink' ? '#007bff' : '#fff',
                    border: '2px solid #007bff',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    transition: 'all 0.2s'
                  }}
                >
                  Drinks
                </button>
                <button
                  onClick={() => setShowForm(!showForm)}
                  style={{
                    padding: '12px 24px',
                    fontSize: '16px',
                    fontWeight: '500',
                    color: '#fff',
                    backgroundColor: showForm ? '#6c757d' : '#28a745',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    marginLeft: 'auto'
                  }}
                >
                  {showForm ? 'Cancel' : '+ New Recipe'}
                </button>
              </div>

              {showForm && (
                <div style={{ marginBottom: '40px' }}>
                  <RecipeForm
                    onRecipeCreated={handleRecipeCreated}
                    defaultRecipeType={recipeType}
                  />
                </div>
              )}

              <RecipeFilters
                onFilterChange={handleFilterChange}
                recipeType={recipeType}
              />

              <RecipeList
                key={`${refreshTrigger}-${recipeType}`}
                onRecipeClick={(recipeId) => setSelectedRecipeId(recipeId)}
                filters={{ ...filters, type: recipeType }}
              />
            </>
          )}
        </div>
      )}
    </div>
  );
}

export default App;
