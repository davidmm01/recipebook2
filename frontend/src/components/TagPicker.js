import React, { useState, useEffect, useRef, useCallback } from 'react';
import { getAllTags } from '../utils/api';

function TagPicker({ selectedTags, onChange, recipeType }) {
  const [allTags, setAllTags] = useState([]);
  const [searchText, setSearchText] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(0);
  const containerRef = useRef(null);
  const inputRef = useRef(null);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const loadTags = async () => {
      try {
        const tags = await getAllTags(recipeType);
        setAllTags(tags || []);
      } catch (err) {
        console.error('Failed to load tags:', err);
      }
    };
    loadTags();
  }, [recipeType]);

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

  const filteredTags = allTags.filter(tag =>
    tag.name.toLowerCase().includes(searchText.toLowerCase())
  );

  const exactMatch = allTags.some(
    tag => tag.name.toLowerCase() === searchText.trim().toLowerCase()
  );
  const showCreateOption = searchText.trim() && !exactMatch;

  const totalOptions = filteredTags.length + (showCreateOption ? 1 : 0);

  // Scroll highlighted item into view
  useEffect(() => {
    if (!isOpen || !dropdownRef.current) return;
    const item = dropdownRef.current.children[highlightedIndex];
    if (item) {
      item.scrollIntoView({ block: 'nearest' });
    }
  }, [highlightedIndex, isOpen]);

  const toggleTag = useCallback((tagName) => {
    const normalized = tagName.trim().toLowerCase();
    if (!normalized) return;

    if (selectedTags.includes(normalized)) {
      onChange(selectedTags.filter(t => t !== normalized));
    } else {
      onChange([...selectedTags, normalized]);
      // Add locally if it's a brand new tag
      if (!allTags.some(t => t.name.toLowerCase() === normalized)) {
        setAllTags(prev => [...prev, { name: normalized, count: 0 }]);
      }
    }
    setSearchText('');
  }, [selectedTags, onChange, allTags]);

  const removeTag = useCallback((tagName) => {
    onChange(selectedTags.filter(t => t !== tagName));
  }, [selectedTags, onChange]);

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
        if (highlightedIndex < filteredTags.length) {
          toggleTag(filteredTags[highlightedIndex].name);
        } else if (showCreateOption) {
          toggleTag(searchText.trim());
        }
      } else if (searchText.trim()) {
        toggleTag(searchText.trim());
      }
    } else if (e.key === 'Escape') {
      setIsOpen(false);
    } else if (e.key === 'Backspace' && searchText === '' && selectedTags.length > 0) {
      removeTag(selectedTags[selectedTags.length - 1]);
    }
  };

  const handleInputChange = (e) => {
    setSearchText(e.target.value);
    setHighlightedIndex(0);
    if (!isOpen) setIsOpen(true);
  };

  const handleInputFocus = () => {
    setIsOpen(true);
  };

  return (
    <div ref={containerRef} style={styles.wrapper}>
      <div
        style={styles.inputArea}
        onClick={() => inputRef.current?.focus()}
      >
        {selectedTags.map(tag => (
          <span key={tag} style={styles.pill}>
            {tag}
            <span
              style={styles.pillRemove}
              onClick={(e) => { e.stopPropagation(); removeTag(tag); }}
            >
              &times;
            </span>
          </span>
        ))}
        <input
          ref={inputRef}
          type="text"
          value={searchText}
          onChange={handleInputChange}
          onFocus={handleInputFocus}
          onKeyDown={handleKeyDown}
          placeholder={selectedTags.length === 0 ? 'Search or create tags...' : ''}
          style={styles.textInput}
        />
      </div>

      {isOpen && (
        <div ref={dropdownRef} style={styles.dropdown}>
          {filteredTags.length === 0 && !showCreateOption && (
            <div style={styles.noResults}>No tags found</div>
          )}
          {filteredTags.map((tag, index) => {
            const isSelected = selectedTags.includes(tag.name);
            const isHighlighted = index === highlightedIndex;
            return (
              <div
                key={tag.name}
                style={{
                  ...styles.dropdownItem,
                  ...(isHighlighted ? styles.dropdownItemHighlighted : {}),
                }}
                onMouseEnter={() => setHighlightedIndex(index)}
                onMouseDown={(e) => {
                  e.preventDefault();
                  toggleTag(tag.name);
                }}
              >
                <span style={styles.checkmark}>
                  {isSelected ? '\u2713' : ''}
                </span>
                <span style={styles.tagName}>{tag.name}</span>
                <span style={styles.tagCount}>{tag.count}</span>
              </div>
            );
          })}
          {showCreateOption && (
            <div
              style={{
                ...styles.dropdownItem,
                ...styles.createOption,
                ...(highlightedIndex === filteredTags.length ? styles.dropdownItemHighlighted : {}),
              }}
              onMouseEnter={() => setHighlightedIndex(filteredTags.length)}
              onMouseDown={(e) => {
                e.preventDefault();
                toggleTag(searchText.trim());
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
  inputArea: {
    display: 'flex',
    flexWrap: 'wrap',
    alignItems: 'center',
    gap: '4px',
    padding: '4px 8px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    backgroundColor: '#fff',
    cursor: 'text',
    minHeight: '38px',
    boxSizing: 'border-box',
  },
  pill: {
    display: 'inline-flex',
    alignItems: 'center',
    gap: '4px',
    padding: '2px 8px',
    backgroundColor: '#e7f1ff',
    color: '#007bff',
    borderRadius: '12px',
    fontSize: '13px',
    lineHeight: '20px',
    whiteSpace: 'nowrap',
  },
  pillRemove: {
    cursor: 'pointer',
    fontWeight: 'bold',
    fontSize: '14px',
    lineHeight: '1',
    color: '#007bff',
    marginLeft: '2px',
  },
  textInput: {
    flex: '1',
    minWidth: '80px',
    border: 'none',
    outline: 'none',
    fontSize: '14px',
    padding: '4px 0',
    fontFamily: 'inherit',
    backgroundColor: 'transparent',
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
  tagName: {
    flex: '1',
  },
  tagCount: {
    color: '#888',
    fontSize: '12px',
    flexShrink: 0,
  },
  createOption: {
    color: '#28a745',
    borderTop: '1px solid #eee',
  },
  noResults: {
    padding: '8px 12px',
    color: '#888',
    fontSize: '14px',
    fontStyle: 'italic',
  },
};

export default TagPicker;
