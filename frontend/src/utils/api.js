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

export async function getRecipes() {
  return await publicFetch('/recipes');
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
