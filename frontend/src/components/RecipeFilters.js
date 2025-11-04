import React, { useState, useEffect } from 'react';
import { getAllTags, getAllCuisines } from '../utils/api';

function RecipeFilters({ onFilterChange, recipeType }) {
  const [search, setSearch] = useState('');
  const [selectedTags, setSelectedTags] = useState([]);
  const [selectedCuisine, setSelectedCuisine] = useState('');
  const [availableTags, setAvailableTags] = useState([]);
  const [availableCuisines, setAvailableCuisines] = useState([]);
  const [loading, setLoading] = useState(true);

  // Load available tags and cuisines when recipeType changes
  useEffect(() => {
    const loadFilterOptions = async () => {
      try {
        setLoading(true);
        const [tags, cuisines] = await Promise.all([
          getAllTags(recipeType),
          getAllCuisines(recipeType)
        ]);
        setAvailableTags(tags || []);
        setAvailableCuisines(cuisines || []);
      } catch (err) {
        console.error('Failed to load filter options:', err);
      } finally {
        setLoading(false);
      }
    };

    loadFilterOptions();
    // Reset filters when recipe type changes
    setSearch('');
    setSelectedTags([]);
    setSelectedCuisine('');
  }, [recipeType]);

  // Notify parent component whenever filters change
  useEffect(() => {
    onFilterChange({
      search,
      tags: selectedTags,
      cuisine: selectedCuisine
    });
  }, [search, selectedTags, selectedCuisine, onFilterChange]);

  const handleTagToggle = (tag) => {
    setSelectedTags(prev =>
      prev.includes(tag)
        ? prev.filter(t => t !== tag)
        : [...prev, tag]
    );
  };

  const handleClearFilters = () => {
    setSearch('');
    setSelectedTags([]);
    setSelectedCuisine('');
  };

  const hasActiveFilters = search || selectedTags.length > 0 || selectedCuisine;

  if (loading) {
    return <div style={styles.container}>Loading filters...</div>;
  }

  return (
    <div style={styles.container}>
      <div style={styles.filterSection}>
        <label style={styles.label}>Search</label>
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search recipes..."
          style={styles.searchInput}
        />
      </div>

      <div style={styles.filterSection}>
        <label style={styles.label}>Cuisine</label>
        <select
          value={selectedCuisine}
          onChange={(e) => setSelectedCuisine(e.target.value)}
          style={styles.select}
        >
          <option value="">All Cuisines</option>
          {availableCuisines.map(cuisine => (
            <option key={cuisine} value={cuisine}>
              {cuisine}
            </option>
          ))}
        </select>
      </div>

      {availableTags.length > 0 && (
        <div style={styles.filterSection}>
          <label style={styles.label}>Tags</label>
          <div style={styles.tagContainer}>
            {availableTags.map(tag => (
              <button
                key={tag}
                onClick={() => handleTagToggle(tag)}
                style={{
                  ...styles.tagButton,
                  ...(selectedTags.includes(tag) ? styles.tagButtonActive : {})
                }}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>
      )}

      {hasActiveFilters && (
        <button onClick={handleClearFilters} style={styles.clearButton}>
          Clear All Filters
        </button>
      )}
    </div>
  );
}

const styles = {
  container: {
    padding: '20px',
    backgroundColor: '#f8f9fa',
    borderRadius: '8px',
    marginBottom: '20px'
  },
  filterSection: {
    marginBottom: '16px'
  },
  label: {
    display: 'block',
    marginBottom: '8px',
    fontSize: '14px',
    fontWeight: '600',
    color: '#333'
  },
  searchInput: {
    width: '100%',
    padding: '10px',
    fontSize: '14px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    boxSizing: 'border-box'
  },
  select: {
    width: '100%',
    padding: '10px',
    fontSize: '14px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    backgroundColor: '#fff',
    cursor: 'pointer',
    boxSizing: 'border-box'
  },
  tagContainer: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: '8px'
  },
  tagButton: {
    padding: '6px 12px',
    fontSize: '13px',
    border: '1px solid #007bff',
    borderRadius: '16px',
    backgroundColor: '#fff',
    color: '#007bff',
    cursor: 'pointer',
    transition: 'all 0.2s'
  },
  tagButtonActive: {
    backgroundColor: '#007bff',
    color: '#fff'
  },
  clearButton: {
    marginTop: '8px',
    padding: '8px 16px',
    fontSize: '14px',
    color: '#dc3545',
    backgroundColor: 'transparent',
    border: '1px solid #dc3545',
    borderRadius: '4px',
    cursor: 'pointer'
  }
};

export default RecipeFilters;
