export const validatePassword = (password) => {
  const minLength = 8;
  const hasUpperCase = /[A-Z]/.test(password);
  const hasLowerCase = /[a-z]/.test(password);
  const hasNumbers = /\d/.test(password);
  const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password);

  const errors = [];
  if (password.length < minLength) errors.push(`Password must be at least ${minLength} characters`);
  if (!hasUpperCase) errors.push('Password must contain at least one uppercase letter');
  if (!hasLowerCase) errors.push('Password must contain at least one lowercase letter');
  if (!hasNumbers) errors.push('Password must contain at least one number');
  if (!hasSpecialChar) errors.push('Password must contain at least one special character');

  return errors;
};

export const validateEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email) ? [] : ['Please enter a valid email address'];
};

export const validateDob = (dob) => {
  const dobDate = new Date(dob);
  const today = new Date();
  const minAge = 16; 
  
  const age = today.getFullYear() - dobDate.getFullYear();
  const errors = [];
  
  if (age < minAge) errors.push(`You must be at least ${minAge} years old`);
  if (dobDate > today) errors.push('Date of birth cannot be in the future');
  
  return errors;
};
