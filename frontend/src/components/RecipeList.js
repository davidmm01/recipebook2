import React, { useState, useEffect } from 'react';
import { getRecipes } from '../utils/api';

function RecipeList({ onRecipeClick, filters = {} }) {
  const [recipes, setRecipes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadRecipes();
  }, [filters]);

  const loadRecipes = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await getRecipes(filters);
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
            onClick={() => onRecipeClick && onRecipeClick(recipe.id)}
            style={{
              border: '1px solid #ddd',
              borderRadius: '8px',
              padding: '15px',
              backgroundColor: '#f9f9f9',
              cursor: 'pointer',
              transition: 'all 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#e9e9e9';
              e.currentTarget.style.boxShadow = '0 2px 8px rgba(0,0,0,0.1)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#f9f9f9';
              e.currentTarget.style.boxShadow = 'none';
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '10px' }}>
              {recipe.icon && (
                <img
                  src={recipe.icon.iconUrl}
                  alt="Recipe icon"
                  style={{
                    width: '40px',
                    height: '40px',
                    objectFit: 'contain',
                    flexShrink: 0
                  }}
                />
              )}
              <h3 style={{ margin: 0, flex: 1 }}>{recipe.title}</h3>
            </div>

            <div style={{ display: 'flex', gap: '8px', marginBottom: '10px', flexWrap: 'wrap' }}>
              {recipe.type && (
                <span
                  style={{
                    display: 'inline-block',
                    backgroundColor: '#007bff',
                    color: 'white',
                    padding: '4px 12px',
                    borderRadius: '12px',
                    fontSize: '12px',
                    fontWeight: '600',
                    textTransform: 'uppercase',
                  }}
                >
                  {recipe.type}
                </span>
              )}
              {recipe.cuisine && (
                <span
                  style={{
                    display: 'inline-block',
                    backgroundColor: '#28a745',
                    color: 'white',
                    padding: '4px 12px',
                    borderRadius: '12px',
                    fontSize: '12px',
                    fontWeight: '600',
                    textTransform: 'uppercase',
                  }}
                >
                  {recipe.cuisine}
                </span>
              )}
            </div>

            {recipe.tags && recipe.tags.length > 0 && (
              <div style={{ display: 'flex', gap: '6px', marginBottom: '10px', flexWrap: 'wrap' }}>
                {recipe.tags.map((tag, index) => (
                  <span
                    key={index}
                    style={{
                      padding: '3px 8px',
                      fontSize: '12px',
                      color: '#007bff',
                      backgroundColor: '#e7f3ff',
                      borderRadius: '4px',
                    }}
                  >
                    #{tag}
                  </span>
                ))}
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
