import React, { useState, useEffect } from 'react';
import { getAllIcons, uploadIcon } from '../utils/api';

function IconManager({ selectedIconId, onIconSelect }) {
  const [icons, setIcons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadIcons();
  }, []);

  const loadIcons = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await getAllIcons();
      setIcons(data || []);
    } catch (err) {
      console.error('Error loading icons:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleFileSelect = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    // Validate file type
    const validTypes = ['image/jpeg', 'image/png', 'image/svg+xml', 'image/webp'];
    if (!validTypes.includes(file.type)) {
      setError('Invalid file type. Please upload a JPG, PNG, SVG, or WebP file.');
      return;
    }

    // Validate file size (2MB max)
    if (file.size > 2 * 1024 * 1024) {
      setError('File size too large. Maximum size is 2MB.');
      return;
    }

    try {
      setUploading(true);
      setError('');
      const newIcon = await uploadIcon(file);
      setIcons([newIcon, ...icons]);
      // Auto-select the newly uploaded icon
      if (onIconSelect) {
        onIconSelect(newIcon.id);
      }
    } catch (err) {
      console.error('Error uploading icon:', err);
      setError(err.message);
    } finally {
      setUploading(false);
      // Reset file input
      e.target.value = '';
    }
  };

  if (loading) {
    return <div style={styles.container}>Loading icons...</div>;
  }

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <label style={styles.label}>Recipe Icon</label>
        <label style={styles.uploadButton}>
          <input
            type="file"
            accept=".jpg,.jpeg,.png,.svg,.webp"
            onChange={handleFileSelect}
            disabled={uploading}
            style={{ display: 'none' }}
          />
          {uploading ? 'Uploading...' : '+ Upload New Icon'}
        </label>
      </div>

      {error && (
        <div style={styles.error}>{error}</div>
      )}

      {icons.length === 0 ? (
        <div style={styles.empty}>
          No icons available. Upload your first icon!
        </div>
      ) : (
        <div style={styles.grid}>
          <div
            onClick={() => onIconSelect && onIconSelect(null)}
            style={{
              ...styles.iconBox,
              ...(selectedIconId === null ? styles.iconBoxSelected : {}),
            }}
          >
            <div style={styles.noIconText}>No Icon</div>
          </div>
          {icons.map((icon) => (
            <div
              key={icon.id}
              onClick={() => onIconSelect && onIconSelect(icon.id)}
              style={{
                ...styles.iconBox,
                ...(selectedIconId === icon.id ? styles.iconBoxSelected : {}),
              }}
            >
              <img
                src={icon.iconUrl}
                alt={icon.filename}
                style={styles.icon}
              />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

const styles = {
  container: {
    marginBottom: '16px',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '8px',
  },
  label: {
    fontSize: '14px',
    fontWeight: '500',
    color: '#555',
  },
  uploadButton: {
    padding: '6px 12px',
    fontSize: '13px',
    color: '#007bff',
    backgroundColor: '#e7f3ff',
    border: '1px solid #007bff',
    borderRadius: '4px',
    cursor: 'pointer',
    transition: 'all 0.2s',
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(80px, 1fr))',
    gap: '12px',
    padding: '12px',
    backgroundColor: '#f8f9fa',
    borderRadius: '4px',
    border: '1px solid #ddd',
  },
  iconBox: {
    width: '80px',
    height: '80px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#fff',
    border: '2px solid #ddd',
    borderRadius: '8px',
    cursor: 'pointer',
    transition: 'all 0.2s',
    overflow: 'hidden',
  },
  iconBoxSelected: {
    borderColor: '#007bff',
    boxShadow: '0 0 0 2px #e7f3ff',
  },
  icon: {
    maxWidth: '100%',
    maxHeight: '100%',
    objectFit: 'contain',
  },
  noIconText: {
    fontSize: '11px',
    color: '#999',
    textAlign: 'center',
    padding: '4px',
  },
  empty: {
    padding: '20px',
    textAlign: 'center',
    color: '#999',
    backgroundColor: '#f8f9fa',
    borderRadius: '4px',
    border: '1px solid #ddd',
  },
  error: {
    padding: '8px 12px',
    backgroundColor: '#f8d7da',
    color: '#721c24',
    borderRadius: '4px',
    marginBottom: '12px',
    fontSize: '13px',
  },
};

export default IconManager;
