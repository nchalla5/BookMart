import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const endpoint ="http://localhost:8080"
function SignInForm({ onToggle }) {
  const [error, setError] = useState('');
  const navigate = useNavigate();
    const handleSubmit = async (event) => {
      event.preventDefault();
      setError('');
      const { emailOrPhone, password } = event.target.elements;
      try {
        const response = await fetch(endpoint + '/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            emailOrPhone: emailOrPhone.value,
            password: password.value,
          }),
        });
  
        if (!response.ok) {
          // throw new Error(`HTTP error! status: ${response.status}`);
          if (response.status === 401) {
            setError('Username and/or Password are invalid.');
          } else {
              // Handle other statuses or a general error message
              setError(`An error occurred: ${response.status}. Please try again.`);
          }
          return;
        }
        const data = await response.json();
        console.log(data);
        localStorage.setItem('token', data.token);
        navigate('/home');
        // Handle the response data here. For example, you might save the token in localStorage
      } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
        // Handle errors here
      }
    };

  return (
    <div className="sign-in-form">
      <h2>Login Form</h2>
      {error && <div className="error-message">{error}</div>}
      <form onSubmit={handleSubmit}>
        <input type="text" name="emailOrPhone" placeholder="Email or Phone" required />
        <input type="password" name="password" placeholder="Password" required />
        {/* <div className='forgot-password-container'><a href="/forgot-password" className="linkButton">
        Forgot password?
        </a> </div> */}
        {/* <div className="forgot-password">Forgot password?</div><br /> */}
        <button type="submit" className="btn-login">Login</button>
        <p> Not a member?  </p>
        <div className="toggle-form" onClick={onToggle}>
          Signup now
        </div>
      </form>
    </div>
  );
}

export default SignInForm;
