import React, { useState, useEffect } from 'react';
import { createRecipe, updateRecipe } from '../utils/api';
import IconManager from './IconManager';

function RecipeForm({ initialRecipe, onRecipeCreated, onRecipeUpdated, defaultRecipeType = 'food' }) {
  const isEditing = !!initialRecipe;

  const [formData, setFormData] = useState({
    title: '',
    description: '',
    type: defaultRecipeType,
    cuisine: '',
    tags: '',
    ingredients: '',
    method: '',
    notes: '',
    sources: '',
    iconId: null
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Load initial data when editing or update type when defaultRecipeType changes
  useEffect(() => {
    if (initialRecipe) {
      setFormData({
        title: initialRecipe.title || '',
        description: initialRecipe.description || '',
        type: initialRecipe.type || defaultRecipeType,
        cuisine: initialRecipe.cuisine || '',
        tags: initialRecipe.tags ? initialRecipe.tags.join(', ') : '',
        ingredients: initialRecipe.ingredients || '',
        method: initialRecipe.method || '',
        notes: initialRecipe.notes || '',
        sources: initialRecipe.sources || '',
        iconId: initialRecipe.iconId || null
      });
    } else {
      // Update type when defaultRecipeType changes for new recipes
      setFormData(prev => ({ ...prev, type: defaultRecipeType }));
    }
  }, [initialRecipe, defaultRecipeType]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      // Convert tags string to array
      const tagsArray = formData.tags
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag.length > 0);

      const recipeData = {
        ...formData,
        tags: tagsArray
      };

      if (isEditing) {
        await updateRecipe(initialRecipe.id, recipeData);
        if (onRecipeUpdated) {
          onRecipeUpdated();
        }
      } else {
        await createRecipe(recipeData);
        // Reset form only when creating
        setFormData({
          title: '',
          description: '',
          type: defaultRecipeType,
          cuisine: '',
          tags: '',
          ingredients: '',
          method: '',
          notes: '',
          sources: '',
          iconId: null
        });
        if (onRecipeCreated) {
          onRecipeCreated();
        }
      }
    } catch (err) {
      setError(err.message || `Failed to ${isEditing ? 'update' : 'create'} recipe`);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  return (
    <div style={styles.container}>
      <h2 style={styles.heading}>{isEditing ? 'Edit Recipe' : 'Create New Recipe'}</h2>

      {error && (
        <div style={styles.error}>
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} style={styles.form}>
        <div style={styles.row}>
          <div style={styles.field}>
            <label style={styles.label}>Title *</label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
              required
              style={styles.input}
              placeholder="e.g., Spaghetti Carbonara"
            />
          </div>

          <div style={styles.field}>
            <label style={styles.label}>Type *</label>
            <select
              name="type"
              value={formData.type}
              onChange={handleChange}
              required
              style={styles.input}
            >
              <option value="food">Food</option>
              <option value="drink">Drink</option>
            </select>
          </div>
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Description</label>
          <input
            type="text"
            name="description"
            value={formData.description}
            onChange={handleChange}
            style={styles.input}
            placeholder="Brief description of the recipe"
          />
        </div>

        <div style={styles.row}>
          <div style={styles.field}>
            <label style={styles.label}>Cuisine</label>
            <input
              type="text"
              name="cuisine"
              value={formData.cuisine}
              onChange={handleChange}
              style={styles.input}
              placeholder="e.g., italian, mexican, chinese"
            />
          </div>

          <div style={styles.field}>
            <label style={styles.label}>Tags</label>
            <input
              type="text"
              name="tags"
              value={formData.tags}
              onChange={handleChange}
              style={styles.input}
              placeholder="pasta, quick, vegetarian (comma-separated)"
            />
          </div>
        </div>

        <IconManager
          selectedIconId={formData.iconId}
          onIconSelect={(iconId) => setFormData(prev => ({ ...prev, iconId }))}
        />

        <div style={styles.field}>
          <label style={styles.label}>Ingredients *</label>
          <textarea
            name="ingredients"
            value={formData.ingredients}
            onChange={handleChange}
            required
            rows="6"
            style={{...styles.input, ...styles.textarea}}
            placeholder="- 400g spaghetti&#10;- 200g pancetta&#10;- 4 eggs"
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Method *</label>
          <textarea
            name="method"
            value={formData.method}
            onChange={handleChange}
            required
            rows="8"
            style={{...styles.input, ...styles.textarea}}
            placeholder="1. Boil pasta&#10;2. Fry pancetta&#10;3. Mix eggs and cheese"
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Notes</label>
          <textarea
            name="notes"
            value={formData.notes}
            onChange={handleChange}
            rows="3"
            style={{...styles.input, ...styles.textarea}}
            placeholder="Optional notes, tips, or variations"
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Sources</label>
          <textarea
            name="sources"
            value={formData.sources}
            onChange={handleChange}
            rows="3"
            style={{...styles.input, ...styles.textarea}}
            placeholder="Recipe sources, references, or attribution (e.g., cookbook name, website URL)"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          style={{
            ...styles.button,
            ...(loading ? styles.buttonDisabled : {})
          }}
        >
          {loading ? (isEditing ? 'Updating...' : 'Creating...') : (isEditing ? 'Update Recipe' : 'Create Recipe')}
        </button>
      </form>
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
  heading: {
    marginTop: 0,
    marginBottom: '24px',
    color: '#333'
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px'
  },
  row: {
    display: 'grid',
    gridTemplateColumns: '1fr 1fr',
    gap: '16px'
  },
  field: {
    display: 'flex',
    flexDirection: 'column',
    gap: '4px'
  },
  label: {
    fontSize: '14px',
    fontWeight: '500',
    color: '#555'
  },
  input: {
    padding: '8px 12px',
    fontSize: '14px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    fontFamily: 'inherit'
  },
  textarea: {
    resize: 'vertical',
    fontFamily: 'monospace'
  },
  button: {
    padding: '12px 24px',
    fontSize: '16px',
    fontWeight: '500',
    color: '#fff',
    backgroundColor: '#007bff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    marginTop: '8px'
  },
  buttonDisabled: {
    backgroundColor: '#6c757d',
    cursor: 'not-allowed'
  },
  error: {
    padding: '12px',
    backgroundColor: '#f8d7da',
    color: '#721c24',
    borderRadius: '4px',
    marginBottom: '16px'
  }
};

export default RecipeForm;
