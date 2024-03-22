import React from 'react';

const endpoint ="http://localhost:8080"
function SignInForm({ onToggle }) {
    const handleSubmit = async (event) => {
      event.preventDefault();
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
          throw new Error(`HTTP error! status: ${response.status}`);
        }
  
        const data = await response.json();
        console.log(data);
        // Handle the response data here. For example, you might save the token in localStorage
      } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
        // Handle errors here
      }
    };

  return (
    <div className="sign-in-form">
      <h2>Login Form</h2>
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
