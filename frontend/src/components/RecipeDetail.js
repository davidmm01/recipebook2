import React, { useState, useEffect } from 'react';
import { getRecipeById, deleteRecipe } from '../utils/api';
import { useUserRole } from '../hooks/useUserRole';
import RecipeForm from './RecipeForm';
import { auth } from '../firebase';

function RecipeDetail({ recipeId, onBack }) {
  const [recipe, setRecipe] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const user = auth.currentUser;
  const { canEdit, isAdmin } = useUserRole(user);

  useEffect(() => {
    const fetchRecipe = async () => {
      try {
        const data = await getRecipeById(recipeId);
        setRecipe(data);
      } catch (err) {
        setError(err.message || 'Failed to load recipe');
      } finally {
        setLoading(false);
      }
    };

    fetchRecipe();
  }, [recipeId]);

  if (loading) {
    return <div style={styles.container}>Loading recipe...</div>;
  }

  if (error) {
    return (
      <div style={styles.container}>
        <div style={styles.error}>{error}</div>
        <button onClick={onBack} style={styles.backButton}>
          ← Back to Recipes
        </button>
      </div>
    );
  }

  if (!recipe) {
    return <div style={styles.container}>Recipe not found</div>;
  }

  const handleRecipeUpdated = async () => {
    // Refresh recipe data after update
    try {
      const data = await getRecipeById(recipeId);
      setRecipe(data);
      setIsEditing(false);
    } catch (err) {
      setError(err.message || 'Failed to reload recipe');
    }
  };

  const handleDelete = async () => {
    if (!window.confirm(`Are you sure you want to delete "${recipe.title}"? This cannot be undone.`)) {
      return;
    }

    try {
      await deleteRecipe(recipeId);
      onBack(); // Navigate back to list after deletion
    } catch (err) {
      setError(err.message || 'Failed to delete recipe');
    }
  };

  // Show edit form if in editing mode
  if (isEditing) {
    return (
      <div style={styles.container}>
        <button onClick={() => setIsEditing(false)} style={styles.backButton}>
          ← Cancel Editing
        </button>
        <RecipeForm
          initialRecipe={recipe}
          onRecipeUpdated={handleRecipeUpdated}
        />
      </div>
    );
  }

  return (
    <div style={styles.container}>
      <div style={{display: 'flex', gap: '10px', marginBottom: '20px', flexWrap: 'wrap'}}>
        <button onClick={onBack} style={styles.backButton}>
          ← Back to Recipes
        </button>
        {canEdit && (
          <button
            onClick={() => setIsEditing(true)}
            style={{...styles.backButton, ...styles.editButton}}
          >
            ✏️ Edit Recipe
          </button>
        )}
        {isAdmin && (
          <button
            onClick={handleDelete}
            style={{...styles.backButton, ...styles.deleteButton}}
          >
            🗑️ Delete Recipe
          </button>
        )}
      </div>

      <div style={styles.header}>
        <h1 style={styles.title}>{recipe.title}</h1>
        {recipe.description && (
          <p style={styles.description}>{recipe.description}</p>
        )}

        <div style={styles.metadata}>
          <span style={styles.badge}>{recipe.type}</span>
          {recipe.cuisine && (
            <span style={{...styles.badge, ...styles.cuisineBadge}}>
              {recipe.cuisine}
            </span>
          )}
        </div>

        {recipe.tags && recipe.tags.length > 0 && (
          <div style={styles.tags}>
            {recipe.tags.map((tag, index) => (
              <span key={index} style={styles.tag}>
                #{tag}
              </span>
            ))}
          </div>
        )}
      </div>

      {recipe.images && recipe.images.length > 0 && (
        <div style={styles.images}>
          {recipe.images.map((image) => (
            <img
              key={image.id}
              src={image.imageUrl}
              alt={`${recipe.title}`}
              style={styles.image}
            />
          ))}
        </div>
      )}

      <div style={styles.section}>
        <h2 style={styles.sectionTitle}>Ingredients</h2>
        <pre style={styles.content}>{recipe.ingredients}</pre>
      </div>

      <div style={styles.section}>
        <h2 style={styles.sectionTitle}>Method</h2>
        <pre style={styles.content}>{recipe.method}</pre>
      </div>

      {recipe.notes && (
        <div style={styles.section}>
          <h2 style={styles.sectionTitle}>Notes</h2>
          <pre style={styles.content}>{recipe.notes}</pre>
        </div>
      )}

      <div style={styles.footer}>
        <small style={styles.timestamp}>
          Created: {new Date(recipe.createdAt).toLocaleDateString()}
        </small>
        {recipe.updatedAt !== recipe.createdAt && (
          <small style={styles.timestamp}>
            Updated: {new Date(recipe.updatedAt).toLocaleDateString()}
          </small>
        )}
      </div>
    </div>
  );
}

const styles = {
  container: {
    maxWidth: '800px',
    margin: '0 auto',
    padding: '20px',
    backgroundColor: '#fff',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
  },
  backButton: {
    padding: '8px 16px',
    fontSize: '14px',
    color: '#007bff',
    backgroundColor: 'transparent',
    border: '1px solid #007bff',
    borderRadius: '4px',
    cursor: 'pointer',
    marginBottom: '20px'
  },
  editButton: {
    color: '#28a745',
    borderColor: '#28a745'
  },
  deleteButton: {
    color: '#dc3545',
    borderColor: '#dc3545'
  },
  header: {
    marginBottom: '24px',
    paddingBottom: '16px',
    borderBottom: '2px solid #eee'
  },
  title: {
    margin: '0 0 8px 0',
    fontSize: '32px',
    color: '#333'
  },
  description: {
    margin: '8px 0 16px 0',
    fontSize: '18px',
    color: '#666',
    fontStyle: 'italic'
  },
  metadata: {
    display: 'flex',
    gap: '8px',
    marginBottom: '12px'
  },
  badge: {
    padding: '4px 12px',
    fontSize: '12px',
    fontWeight: '600',
    textTransform: 'uppercase',
    backgroundColor: '#007bff',
    color: '#fff',
    borderRadius: '12px'
  },
  cuisineBadge: {
    backgroundColor: '#28a745'
  },
  tags: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: '8px'
  },
  tag: {
    padding: '4px 8px',
    fontSize: '13px',
    color: '#007bff',
    backgroundColor: '#e7f3ff',
    borderRadius: '4px'
  },
  images: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
    gap: '16px',
    marginBottom: '24px'
  },
  image: {
    width: '100%',
    height: 'auto',
    borderRadius: '8px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
  },
  section: {
    marginBottom: '24px'
  },
  sectionTitle: {
    margin: '0 0 12px 0',
    fontSize: '20px',
    color: '#333',
    borderLeft: '4px solid #007bff',
    paddingLeft: '12px'
  },
  content: {
    margin: 0,
    padding: '16px',
    fontSize: '15px',
    lineHeight: '1.6',
    backgroundColor: '#f8f9fa',
    borderRadius: '4px',
    whiteSpace: 'pre-wrap',
    fontFamily: 'inherit'
  },
  footer: {
    display: 'flex',
    gap: '16px',
    marginTop: '32px',
    paddingTop: '16px',
    borderTop: '1px solid #eee'
  },
  timestamp: {
    color: '#999',
    fontSize: '13px'
  },
  error: {
    padding: '12px',
    backgroundColor: '#f8d7da',
    color: '#721c24',
    borderRadius: '4px',
    marginBottom: '16px'
  }
};

export default RecipeDetail;
