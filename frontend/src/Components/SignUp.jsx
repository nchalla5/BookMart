import React from 'react';
import { useState } from 'react';

const endpoint ="http://localhost:8080"
function SignUpForm({ onToggle }) {
    const [error, setError] = useState('');
  
    const handleSubmit = async (event) => {
      event.preventDefault();
      setError(''); // Reset error message
  
      const { name, email, password, confirmPassword } = event.target.elements;
      
      if (password.value !== confirmPassword.value) {
        // Set error message if passwords don't match
        setError('Passwords do not match.');
        return;
      }
  
      try {
        const response = await fetch(endpoint + '/signup', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name: name.value,
            email: email.value,
            password: password.value,
            confirmPassword: confirmPassword.value,
          }),
        });
  
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
  
        const data = await response.json();
        console.log(data);
        // Handle successful sign up here (e.g., redirect to login or dashboard)
      } catch (error) {
        console.error('There was a problem with the signup operation:', error);
        setError(error.message);
      }
};
  
  

  return (
    <div className="sign-up-form">
      <h2>Registration</h2>
      <form onSubmit={handleSubmit}>
        <input type="text" name="name" placeholder="Enter your name" required />
        <input type="email" name="email" placeholder="Enter your email" required />
        <input type="password" name="password" placeholder="Create password" required />
        <input type="password" name="confirmPassword" placeholder="Confirm password" required />
        {error && <div className="error-message">{error}</div>}

        <div className="terms">
            <input type="checkbox" id="terms" className="checkbox" />
            <label htmlFor="terms">I accept all terms & conditions</label>
        </div><br />


        <button type="submit" className="button">Register Now</button>
        <p>Already have an account ?</p>
        <div className="toggle-form" onClick={onToggle}>
        Login now
        </div>
      </form>
    </div>
  );
}

export default SignUpForm;
