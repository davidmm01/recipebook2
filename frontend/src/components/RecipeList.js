import React, { useState, useEffect } from 'react';
import { getRecipes } from '../utils/api';

function RecipeList() {
  const [recipes, setRecipes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadRecipes();
  }, []);

  const loadRecipes = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await getRecipes();
      setRecipes(data || []);
    } catch (err) {
      console.error('Error loading recipes:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div style={{ padding: '20px' }}>Loading recipes...</div>;
  }

  if (error) {
    return (
      <div style={{ padding: '20px', color: 'red' }}>
        <p>Error loading recipes: {error}</p>
        <button onClick={loadRecipes}>Retry</button>
      </div>
    );
  }

  if (recipes.length === 0) {
    return (
      <div style={{ padding: '20px' }}>
        <p>No recipes yet. Create your first recipe!</p>
      </div>
    );
  }

  return (
    <div style={{ padding: '20px' }}>
      <h2>Recipes ({recipes.length})</h2>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
        {recipes.map((recipe) => (
          <div
            key={recipe.id}
            style={{
              border: '1px solid #ddd',
              borderRadius: '8px',
              padding: '15px',
              backgroundColor: '#f9f9f9',
            }}
          >
            <h3 style={{ margin: '0 0 10px 0' }}>{recipe.title}</h3>
            {recipe.type && (
              <span
                style={{
                  display: 'inline-block',
                  backgroundColor: '#4285f4',
                  color: 'white',
                  padding: '4px 8px',
                  borderRadius: '4px',
                  fontSize: '12px',
                  marginBottom: '10px',
                }}
              >
                {recipe.type}
              </span>
            )}
            {recipe.ingredients && (
              <div style={{ marginTop: '10px' }}>
                <strong>Ingredients:</strong>
                <div
                  style={{
                    whiteSpace: 'pre-wrap',
                    fontSize: '14px',
                    marginTop: '5px',
                  }}
                >
                  {recipe.ingredients.substring(0, 150)}
                  {recipe.ingredients.length > 150 && '...'}
                </div>
              </div>
            )}
            <div style={{ marginTop: '10px', fontSize: '12px', color: '#666' }}>
              Updated: {new Date(recipe.updatedAt).toLocaleDateString()}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default RecipeList;
