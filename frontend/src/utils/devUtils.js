import { doc, updateDoc } from 'firebase/firestore';
import { db } from '../firebase';
import { auth } from '../firebase';

// Development-only utility to upgrade current user to admin
export const upgradeToAdmin = async () => {
  const user = auth.currentUser;
  if (!user) {
    throw new Error('No user logged in');
  }

  const userRef = doc(db, 'users', user.uid);
  await updateDoc(userRef, {
    role: 'admin'
  });

  console.log('✅ Upgraded to admin! Refresh the page.');
  alert('You are now an admin! Refresh the page to see changes.');
};

// Utility to upgrade to editor
export const upgradeToEditor = async () => {
  const user = auth.currentUser;
  if (!user) {
    throw new Error('No user logged in');
  }

  const userRef = doc(db, 'users', user.uid);
  await updateDoc(userRef, {
    role: 'editor'
  });

  console.log('✅ Upgraded to editor! Refresh the page.');
  alert('You are now an editor! Refresh the page to see changes.');
};
