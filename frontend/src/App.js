import React, { useState, useEffect, useCallback } from 'react';
import { onAuthStateChanged } from 'firebase/auth';
import { auth } from './firebase';
import Login from './components/Login';
import RecipeList from './components/RecipeList';
import RecipeForm from './components/RecipeForm';
import RecipeDetail from './components/RecipeDetail';
import RecipeFilters from './components/RecipeFilters';
import UserProfile from './components/UserProfile';
import { useUserRole } from './hooks/useUserRole';

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showForm, setShowForm] = useState(false);
  const [selectedRecipeId, setSelectedRecipeId] = useState(null);
  const [filters, setFilters] = useState({});
  const [showProfile, setShowProfile] = useState(false);
  const [recipeType, setRecipeType] = useState('food');
  const { canEdit } = useUserRole(user);

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
    <div style={{ padding: '20px', boxSizing: 'border-box', maxWidth: '1400px', margin: '0 auto', overflowX: 'hidden' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', gap: '10px' }}>
        <div style={{ flex: 1 }}></div>
        <h1 style={{ margin: 0, flex: 1, textAlign: 'center' }}>recipebook2</h1>
        <div style={{ flex: 1, display: 'flex', gap: '10px', alignItems: 'center', justifyContent: 'flex-end' }}>
          <Login user={user} onProfileClick={() => setShowProfile(true)} />
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
              <div style={{ marginBottom: '20px', display: 'flex', gap: '10px', alignItems: 'center', flexWrap: 'wrap' }}>
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

              {canEdit && (
                <div style={{ marginBottom: '20px' }}>
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
                      cursor: 'pointer'
                    }}
                  >
                    {showForm ? 'Cancel' : '+ New Recipe'}
                  </button>
                </div>
              )}

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
