import React, { useState, useEffect } from 'react';
import { getRecipeById, deleteRecipe, getMakeLogs, createMakeLog } from '../utils/api';
import { useUserRole } from '../hooks/useUserRole';
import RecipeForm from './RecipeForm';
import { auth } from '../firebase';

function RecipeDetail({ recipeId, onBack }) {
  const [recipe, setRecipe] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [makeLogs, setMakeLogs] = useState([]);
  const [showMakeLogs, setShowMakeLogs] = useState(false);
  const [showMakeForm, setShowMakeForm] = useState(false);
  const [makeFormData, setMakeFormData] = useState({
    madeAt: new Date().toISOString().split('T')[0], // Today's date in YYYY-MM-DD
    notes: ''
  });
  const [submittingMakeLog, setSubmittingMakeLog] = useState(false);
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

    const fetchMakeLogs = async () => {
      try {
        const logs = await getMakeLogs(recipeId);
        setMakeLogs(logs || []);
      } catch (err) {
        console.error('Failed to load make logs:', err);
      }
    };

    fetchRecipe();
    fetchMakeLogs();
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

  const handleMakeLogSubmit = async (e) => {
    e.preventDefault();
    setSubmittingMakeLog(true);

    try {
      await createMakeLog(recipeId, makeFormData);
      // Refresh make logs and recipe
      const [logs, updatedRecipe] = await Promise.all([
        getMakeLogs(recipeId),
        getRecipeById(recipeId)
      ]);
      setMakeLogs(logs || []);
      setRecipe(updatedRecipe);
      setShowMakeForm(false);
      setMakeFormData({
        madeAt: new Date().toISOString().split('T')[0],
        notes: ''
      });
    } catch (err) {
      setError(err.message || 'Failed to log make');
    } finally {
      setSubmittingMakeLog(false);
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
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '8px' }}>
          {recipe.icon && (
            <img
              src={recipe.icon.iconUrl}
              alt="Recipe icon"
              style={{
                width: '60px',
                height: '60px',
                objectFit: 'contain',
                flexShrink: 0
              }}
            />
          )}
          <h1 style={{ ...styles.title, margin: 0, flex: 1 }}>{recipe.title}</h1>
        </div>
        {recipe.description && (
          <p style={styles.description}>{recipe.description}</p>
        )}

        <div style={styles.metadata}>
          <div style={styles.metadataLeft}>
            {recipe.makeCount > 0 && (
              <span style={styles.makeCountBadge}>
                👨‍🍳 {recipe.makeCount}
              </span>
            )}
          </div>
          <div style={styles.metadataRight}>
            <span style={styles.badge}>{recipe.type}</span>
            {recipe.cuisine && (
              <span style={{...styles.badge, ...styles.cuisineBadge}}>
                {recipe.cuisine}
              </span>
            )}
          </div>
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

      {/* Make Log Button - At top */}
      {user && (
        <div style={styles.makeLogSection}>
          <div style={styles.makeActions}>
            <button
              onClick={() => setShowMakeForm(true)}
              style={styles.madeThisButton}
            >
              ✨ I made this
            </button>

            {makeLogs.length > 0 && (
              <button
                onClick={() => setShowMakeLogs(!showMakeLogs)}
                style={styles.viewHistoryButton}
              >
                {showMakeLogs ? 'Hide History' : '📜 View History'}
              </button>
            )}
          </div>

          {/* Show history inline when toggled */}
          {showMakeLogs && (
            <div style={styles.makeHistory}>
              <h3 style={styles.historyTitle}>Make History</h3>
              {makeLogs.map((log) => (
                <div key={log.id} style={styles.logEntry}>
                  <div style={styles.logDate}>
                    {new Date(log.madeAt).toLocaleDateString('en-US', {
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric'
                    })}
                  </div>
                  {log.notes && (
                    <div style={styles.logNotes}>{log.notes}</div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Make Log Modal */}
      {showMakeForm && (
        <div style={styles.modalOverlay} onClick={() => setShowMakeForm(false)}>
          <div style={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <div style={styles.modalHeader}>
              <h2 style={styles.modalTitle}>Log a Make</h2>
              <button
                onClick={() => setShowMakeForm(false)}
                style={styles.modalCloseButton}
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleMakeLogSubmit} style={styles.makeForm}>
              <div style={styles.formField}>
                <label style={styles.formLabel}>Date</label>
                <input
                  type="date"
                  value={makeFormData.madeAt}
                  onChange={(e) => setMakeFormData({...makeFormData, madeAt: e.target.value})}
                  required
                  style={styles.dateInput}
                />
              </div>
              <div style={styles.formField}>
                <label style={styles.formLabel}>Notes (optional)</label>
                <textarea
                  value={makeFormData.notes}
                  onChange={(e) => setMakeFormData({...makeFormData, notes: e.target.value})}
                  rows="4"
                  placeholder="How did it turn out? Any modifications?"
                  style={styles.notesInput}
                />
              </div>
              <div style={styles.modalActions}>
                <button
                  type="button"
                  onClick={() => setShowMakeForm(false)}
                  style={styles.cancelButton}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={submittingMakeLog}
                  style={{
                    ...styles.submitButton,
                    ...(submittingMakeLog ? styles.submitButtonDisabled : {})
                  }}
                >
                  {submittingMakeLog ? 'Saving...' : 'Save'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

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

      {recipe.sources && (
        <div style={styles.section}>
          <h2 style={styles.sectionTitle}>Sources</h2>
          <pre style={styles.content}>{recipe.sources}</pre>
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
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    boxSizing: 'border-box',
    width: '100%'
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
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '12px'
  },
  metadataLeft: {
    display: 'flex',
    gap: '8px'
  },
  metadataRight: {
    display: 'flex',
    gap: '8px'
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
  makeCountBadge: {
    padding: '4px 12px',
    fontSize: '16px',
    fontWeight: '600',
    color: '#28a745',
    backgroundColor: '#e7f8ed',
    borderRadius: '12px',
    border: '1px solid #28a745'
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
  },
  makeLogSection: {
    marginBottom: '24px'
  },
  makeActions: {
    display: 'flex',
    gap: '10px',
    marginBottom: '16px',
    flexWrap: 'wrap'
  },
  madeThisButton: {
    padding: '10px 20px',
    fontSize: '16px',
    fontWeight: '500',
    color: '#fff',
    backgroundColor: '#28a745',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    transition: 'background-color 0.2s'
  },
  viewHistoryButton: {
    padding: '10px 20px',
    fontSize: '16px',
    fontWeight: '500',
    color: '#007bff',
    backgroundColor: 'transparent',
    border: '1px solid #007bff',
    borderRadius: '4px',
    cursor: 'pointer'
  },
  makeForm: {
    padding: '0'
  },
  formField: {
    marginBottom: '16px'
  },
  formLabel: {
    display: 'block',
    marginBottom: '8px',
    fontSize: '14px',
    fontWeight: '500',
    color: '#333'
  },
  dateInput: {
    width: '100%',
    padding: '10px',
    fontSize: '14px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    boxSizing: 'border-box'
  },
  notesInput: {
    width: '100%',
    padding: '10px',
    fontSize: '14px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    resize: 'vertical',
    fontFamily: 'inherit',
    boxSizing: 'border-box'
  },
  submitButton: {
    padding: '10px 24px',
    fontSize: '16px',
    fontWeight: '500',
    color: '#fff',
    backgroundColor: '#007bff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer'
  },
  submitButtonDisabled: {
    backgroundColor: '#6c757d',
    cursor: 'not-allowed'
  },
  makeHistory: {
    backgroundColor: '#fff',
    padding: '20px',
    borderRadius: '8px',
    marginTop: '16px',
    border: '1px solid #ddd'
  },
  historyTitle: {
    margin: '0 0 16px 0',
    fontSize: '18px',
    color: '#333'
  },
  logEntry: {
    padding: '12px 0',
    borderBottom: '1px solid #eee'
  },
  logDate: {
    fontSize: '15px',
    fontWeight: '600',
    color: '#333',
    marginBottom: '4px'
  },
  logNotes: {
    fontSize: '14px',
    color: '#666',
    fontStyle: 'italic'
  },
  modalOverlay: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000
  },
  modalContent: {
    backgroundColor: '#fff',
    borderRadius: '8px',
    padding: '24px',
    maxWidth: '500px',
    width: '90%',
    maxHeight: '90vh',
    overflow: 'auto',
    boxShadow: '0 4px 20px rgba(0, 0, 0, 0.3)',
    boxSizing: 'border-box'
  },
  modalHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '20px'
  },
  modalTitle: {
    margin: 0,
    fontSize: '24px',
    color: '#333'
  },
  modalCloseButton: {
    background: 'none',
    border: 'none',
    fontSize: '24px',
    color: '#999',
    cursor: 'pointer',
    padding: '0',
    width: '30px',
    height: '30px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
  modalActions: {
    display: 'flex',
    gap: '10px',
    justifyContent: 'flex-end',
    marginTop: '20px'
  },
  cancelButton: {
    padding: '10px 24px',
    fontSize: '16px',
    fontWeight: '500',
    color: '#666',
    backgroundColor: 'transparent',
    border: '1px solid #ddd',
    borderRadius: '4px',
    cursor: 'pointer'
  }
};

export default RecipeDetail;
