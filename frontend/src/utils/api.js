import { auth } from '../firebase';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

/**
 * Get the current user's ID token for authenticated requests
 */
async function getAuthToken() {
  const user = auth.currentUser;
  if (!user) {
    throw new Error('User not authenticated');
  }
  return await user.getIdToken();
}

/**
 * Make an authenticated API request
 */
async function authenticatedFetch(url, options = {}) {
  const token = await getAuthToken();

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
    ...options.headers,
  };

  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(`API Error: ${response.status} - ${error}`);
  }

  // Handle 204 No Content responses
  if (response.status === 204) {
    return null;
  }

  return await response.json();
}

/**
 * Make a public API request (no auth required)
 */
async function publicFetch(url, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(`API Error: ${response.status} - ${error}`);
  }

  return await response.json();
}

// Recipe API functions

export async function getRecipes(filters = {}) {
  // Build query string from filters
  const params = new URLSearchParams();

  if (filters.search) {
    params.append('search', filters.search);
  }

  if (filters.cuisine) {
    params.append('cuisine', filters.cuisine);
  }

  if (filters.type) {
    params.append('type', filters.type);
  }

  if (filters.tags && filters.tags.length > 0) {
    filters.tags.forEach(tag => {
      params.append('tags', tag);
    });
  }

  const queryString = params.toString();
  const url = queryString ? `/recipes?${queryString}` : '/recipes';

  return await publicFetch(url);
}

export async function getRecipeById(id) {
  return await publicFetch(`/recipes/${id}`);
}

export async function createRecipe(recipe) {
  return await authenticatedFetch('/recipes', {
    method: 'POST',
    body: JSON.stringify(recipe),
  });
}

export async function updateRecipe(id, recipe) {
  return await authenticatedFetch(`/recipes/${id}`, {
    method: 'PUT',
    body: JSON.stringify(recipe),
  });
}

export async function deleteRecipe(id) {
  return await authenticatedFetch(`/recipes/${id}`, {
    method: 'DELETE',
  });
}

export async function searchRecipes(query) {
  return await publicFetch(`/recipes/search?q=${encodeURIComponent(query)}`);
}

// Filter metadata functions

export async function getAllTags(recipeType = null) {
  const url = recipeType ? `/tags?type=${encodeURIComponent(recipeType)}` : '/tags';
  return await publicFetch(url);
}

export async function getAllCuisines(recipeType = null) {
  const url = recipeType ? `/cuisines?type=${encodeURIComponent(recipeType)}` : '/cuisines';
  return await publicFetch(url);
}

// Icon API functions

export async function getAllIcons() {
  return await publicFetch('/icons');
}

export async function uploadIcon(file) {
  const token = await getAuthToken();

  const formData = new FormData();
  formData.append('icon', file);

  const response = await fetch(`${API_BASE_URL}/icons`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      // Don't set Content-Type for FormData - browser sets it with boundary
    },
    body: formData,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(`API Error: ${response.status} - ${error}`);
  }

  return await response.json();
}

// User Profile API functions

export async function getUserProfile() {
  return await authenticatedFetch('/user/profile');
}

export async function updateUserProfile(updates) {
  return await authenticatedFetch('/user/profile', {
    method: 'PUT',
    body: JSON.stringify(updates),
  });
}
