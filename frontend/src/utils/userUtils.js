import { doc, getDoc, setDoc, serverTimestamp } from 'firebase/firestore';
import { db } from '../firebase';

/**
 * User roles in the system
 * - viewer: Can only view recipes
 * - editor: Can create and edit recipes
 * - admin: Full access to all recipes and user management
 */
export const USER_ROLES = {
  VIEWER: 'viewer',
  EDITOR: 'editor',
  ADMIN: 'admin',
};

/**
 * Creates or updates a user document in Firestore
 * @param {Object} user - Firebase Auth user object
 * @returns {Promise<Object>} The user document data
 */
export async function createOrUpdateUser(user) {
  if (!user) return null;

  const userRef = doc(db, 'users', user.uid);

  try {
    const userSnap = await getDoc(userRef);

    if (userSnap.exists()) {
      // User exists, update last login
      await setDoc(
        userRef,
        {
          lastLoginAt: serverTimestamp(),
        },
        { merge: true }
      );
      return userSnap.data();
    } else {
      // New user, create document with default role
      const newUserData = {
        email: user.email,
        role: USER_ROLES.VIEWER,
        createdAt: serverTimestamp(),
        lastLoginAt: serverTimestamp(),
      };

      await setDoc(userRef, newUserData);
      return newUserData;
    }
  } catch (error) {
    console.error('Error creating/updating user:', error);
    throw error;
  }
}

/**
 * Gets a user's role from Firestore
 * @param {string} userId - The user's UID
 * @returns {Promise<string|null>} The user's role or null
 */
export async function getUserRole(userId) {
  if (!userId) return null;

  try {
    const userRef = doc(db, 'users', userId);
    const userSnap = await getDoc(userRef);

    if (userSnap.exists()) {
      return userSnap.data().role || USER_ROLES.VIEWER;
    }
    return null;
  } catch (error) {
    console.error('Error getting user role:', error);
    return null;
  }
}

/**
 * Checks if a user has a specific role or higher privileges
 * @param {string} userRole - The user's current role
 * @param {string} requiredRole - The required role to check against
 * @returns {boolean} Whether the user has sufficient privileges
 */
export function hasRole(userRole, requiredRole) {
  const roleHierarchy = {
    [USER_ROLES.VIEWER]: 1,
    [USER_ROLES.EDITOR]: 2,
    [USER_ROLES.ADMIN]: 3,
  };

  return roleHierarchy[userRole] >= roleHierarchy[requiredRole];
}
