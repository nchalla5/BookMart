import React, { useState } from 'react';
// import { useNavigate } from 'react-router-dom';
import Swal from 'sweetalert2';

const endpoint ="http://localhost:8080"
function SignUpForm({ onToggle }) {
    const [error, setError] = useState('');
    // const navigate = useNavigate();
  
    const handleSubmit = async (event) => {
      event.preventDefault();
      setError(''); // Reset error message
  
      const { name, email, password, confirmPassword, terms } = event.target.elements;
      
      if (password.value !== confirmPassword.value) {
        // Set error message if passwords don't match
        setError('Passwords do not match.');
        return;
      }
      if (!terms.checked) {
        // Set error message if terms are not accepted
        setError('You must accept the terms and conditions.');
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
          // throw new Error(`HTTP error! status: ${response.status}`);
          if (response.status === 409) {
            setError(`Account already exists for ${email.value}`);
          }
          else {
          // Handle other statuses or a general error message
          setError(`An error occurred: ${response.status}. Please try again.`);
      }
      return;

        }
  
        const data = await response.json();
        console.log(data);
        // window.alert('User creation is successful'); // Display an alert
        // onToggle();
        Swal.fire({
          title: 'Success!',
          text: 'User creation is successful',
          icon: 'success',
          confirmButtonText: 'OK'
      }).then((result) => {
          if (result.isConfirmed) {
              onToggle(); // Redirect to login when OK is clicked
          }
      });

        // Handle successful sign up here (e.g., redirect to login or dashboard)
      } catch (error) {
        console.error('There was a problem with the signup operation:', error);
        setError(error.message);
      }
};
  
  

  return (
    <div className="sign-up-form">
      <h2>Registration</h2>
      
      {error && <div className="error-message">{error}</div>}
      <form onSubmit={handleSubmit}>
        <input type="text" name="name" placeholder="Enter your name" required />
        <input type="email" name="email" placeholder="Enter your email" required />
        <input type="password" name="password" placeholder="Create password" required />
        <input type="password" name="confirmPassword" placeholder="Confirm password" required />

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
