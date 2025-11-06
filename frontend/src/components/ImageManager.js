import React, { useState } from 'react';
import { uploadImage } from '../utils/api';

function ImageManager({ recipeId, existingImages = [], onImageUploaded }) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');

  const handleFileSelect = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    // Validate file type
    const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
    if (!validTypes.includes(file.type)) {
      setError('Invalid file type. Please upload a JPG, PNG, GIF, or WebP file.');
      return;
    }

    // Validate file size (10MB max)
    if (file.size > 10 * 1024 * 1024) {
      setError('File size too large. Maximum size is 10MB.');
      return;
    }

    try {
      setUploading(true);
      setError('');
      const result = await uploadImage(file, recipeId);

      // Notify parent component
      if (onImageUploaded) {
        onImageUploaded(result.imageUrl);
      }
    } catch (err) {
      console.error('Error uploading image:', err);
      setError(err.message);
    } finally {
      setUploading(false);
      // Reset file input
      e.target.value = '';
    }
  };

  if (!recipeId) {
    return (
      <div style={styles.container}>
        <div style={styles.header}>
          <label style={styles.label}>Recipe Images</label>
        </div>
        <div style={styles.infoBox}>
          Images can be added after creating the recipe.
        </div>
      </div>
    );
  }

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <label style={styles.label}>Recipe Images</label>
        <label style={styles.uploadButton}>
          <input
            type="file"
            accept=".jpg,.jpeg,.png,.gif,.webp"
            onChange={handleFileSelect}
            disabled={uploading}
            style={{ display: 'none' }}
          />
          {uploading ? 'Uploading...' : '+ Upload Image'}
        </label>
      </div>

      {error && (
        <div style={styles.error}>{error}</div>
      )}

      {existingImages.length > 0 && (
        <div style={styles.grid}>
          {existingImages.map((image, index) => (
            <div key={image.id || index} style={styles.imageBox}>
              <img
                src={image.imageUrl}
                alt={`Recipe ${index + 1}`}
                style={styles.image}
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
    color: '#28a745',
    backgroundColor: '#e7f8ed',
    border: '1px solid #28a745',
    borderRadius: '4px',
    cursor: 'pointer',
    transition: 'all 0.2s',
  },
  infoBox: {
    padding: '12px',
    backgroundColor: '#e7f3ff',
    color: '#007bff',
    borderRadius: '4px',
    border: '1px solid #007bff',
    fontSize: '13px',
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(150px, 1fr))',
    gap: '12px',
    padding: '12px',
    backgroundColor: '#f8f9fa',
    borderRadius: '4px',
    border: '1px solid #ddd',
  },
  imageBox: {
    width: '100%',
    aspectRatio: '1',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#fff',
    border: '2px solid #ddd',
    borderRadius: '8px',
    overflow: 'hidden',
  },
  image: {
    width: '100%',
    height: '100%',
    objectFit: 'cover',
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

export default ImageManager;
