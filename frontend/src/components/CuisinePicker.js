import React, { useState, useEffect, useRef } from 'react';
import { getAllCuisines } from '../utils/api';

function CuisinePicker({ value, onChange, recipeType }) {
  const [allCuisines, setAllCuisines] = useState([]);
  const [searchText, setSearchText] = useState(value || '');
  const [isOpen, setIsOpen] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(0);
  const containerRef = useRef(null);
  const inputRef = useRef(null);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const loadCuisines = async () => {
      try {
        const cuisines = await getAllCuisines(recipeType);
        setAllCuisines(cuisines || []);
      } catch (err) {
        console.error('Failed to load cuisines:', err);
      }
    };
    loadCuisines();
  }, [recipeType]);

  // Sync searchText when value changes externally (e.g. form reset or edit load)
  useEffect(() => {
    setSearchText(value || '');
  }, [value]);

  // Close dropdown on outside click
  useEffect(() => {
    const handleClickOutside = (e) => {
      if (containerRef.current && !containerRef.current.contains(e.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const filteredCuisines = allCuisines.filter(c =>
    c.toLowerCase().includes(searchText.toLowerCase())
  );

  const exactMatch = allCuisines.some(
    c => c.toLowerCase() === searchText.trim().toLowerCase()
  );
  const showCreateOption = searchText.trim() && !exactMatch;

  const totalOptions = filteredCuisines.length + (showCreateOption ? 1 : 0);

  // Scroll highlighted item into view
  useEffect(() => {
    if (!isOpen || !dropdownRef.current) return;
    const item = dropdownRef.current.children[highlightedIndex];
    if (item) {
      item.scrollIntoView({ block: 'nearest' });
    }
  }, [highlightedIndex, isOpen]);

  const selectCuisine = (cuisine) => {
    const normalized = cuisine.trim().toLowerCase();
    onChange(normalized);
    setSearchText(normalized);
    setIsOpen(false);
    // Add locally if brand new
    if (!allCuisines.some(c => c.toLowerCase() === normalized)) {
      setAllCuisines(prev => [...prev, normalized]);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (!isOpen) {
        setIsOpen(true);
      } else {
        setHighlightedIndex(prev => (prev + 1) % totalOptions);
      }
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (isOpen) {
        setHighlightedIndex(prev => (prev - 1 + totalOptions) % totalOptions);
      }
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (isOpen && totalOptions > 0) {
        if (highlightedIndex < filteredCuisines.length) {
          selectCuisine(filteredCuisines[highlightedIndex]);
        } else if (showCreateOption) {
          selectCuisine(searchText.trim());
        }
      } else if (searchText.trim()) {
        selectCuisine(searchText.trim());
      }
    } else if (e.key === 'Escape') {
      setIsOpen(false);
    }
  };

  const handleInputChange = (e) => {
    setSearchText(e.target.value);
    onChange(e.target.value);
    setHighlightedIndex(0);
    if (!isOpen) setIsOpen(true);
  };

  const handleInputFocus = () => {
    setIsOpen(true);
  };

  return (
    <div ref={containerRef} style={styles.wrapper}>
      <input
        ref={inputRef}
        type="text"
        value={searchText}
        onChange={handleInputChange}
        onFocus={handleInputFocus}
        onKeyDown={handleKeyDown}
        placeholder="e.g., italian, mexican, chinese"
        style={styles.input}
      />

      {isOpen && totalOptions > 0 && (
        <div ref={dropdownRef} style={styles.dropdown}>
          {filteredCuisines.map((cuisine, index) => {
            const isSelected = value && value.toLowerCase() === cuisine.toLowerCase();
            const isHighlighted = index === highlightedIndex;
            return (
              <div
                key={cuisine}
                style={{
                  ...styles.dropdownItem,
                  ...(isHighlighted ? styles.dropdownItemHighlighted : {}),
                }}
                onMouseEnter={() => setHighlightedIndex(index)}
                onMouseDown={(e) => {
                  e.preventDefault();
                  selectCuisine(cuisine);
                }}
              >
                <span style={styles.checkmark}>
                  {isSelected ? '\u2713' : ''}
                </span>
                <span style={styles.cuisineName}>{cuisine}</span>
              </div>
            );
          })}
          {showCreateOption && (
            <div
              style={{
                ...styles.dropdownItem,
                ...styles.createOption,
                ...(highlightedIndex === filteredCuisines.length ? styles.dropdownItemHighlighted : {}),
              }}
              onMouseEnter={() => setHighlightedIndex(filteredCuisines.length)}
              onMouseDown={(e) => {
                e.preventDefault();
                selectCuisine(searchText.trim());
              }}
            >
              Create "<strong>{searchText.trim()}</strong>"
            </div>
          )}
        </div>
      )}
    </div>
  );
}

const styles = {
  wrapper: {
    position: 'relative',
  },
  input: {
    width: '100%',
    padding: '8px 12px',
    fontSize: '14px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    fontFamily: 'inherit',
    boxSizing: 'border-box',
  },
  dropdown: {
    position: 'absolute',
    top: '100%',
    left: 0,
    right: 0,
    marginTop: '4px',
    backgroundColor: '#fff',
    border: '1px solid #ddd',
    borderRadius: '4px',
    boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
    maxHeight: '240px',
    overflowY: 'auto',
    zIndex: 1000,
  },
  dropdownItem: {
    display: 'flex',
    alignItems: 'center',
    padding: '8px 12px',
    cursor: 'pointer',
    fontSize: '14px',
    gap: '8px',
  },
  dropdownItemHighlighted: {
    backgroundColor: '#f0f0f0',
  },
  checkmark: {
    width: '16px',
    textAlign: 'center',
    color: '#007bff',
    fontWeight: 'bold',
    flexShrink: 0,
  },
  cuisineName: {
    flex: '1',
  },
  createOption: {
    color: '#28a745',
    borderTop: '1px solid #eee',
  },
};

export default CuisinePicker;
