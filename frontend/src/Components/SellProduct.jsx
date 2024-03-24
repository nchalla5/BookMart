import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import './SellProduct.css';
// import { jwtDecode } from 'jwt-decode';
import Swal from 'sweetalert2';

const SellProduct = () => {
const [product, setProduct] = useState({
title: '',
cost: '',
location: '',
description: '',
image: null,
});
const navigate = useNavigate();
const token = localStorage.getItem('token');
if (!token) {
  navigate('/'); // Redirect to login if token not found
  return;
}
const handleChange = (e) => {
const { name, value } = e.target;
setProduct({
    ...product,
    [name]: value,
});
};

const handleImageChange = (e) => {
setProduct({
    ...product,
    image: e.target.files[0],
});
};

const handleLogout = () => {
    localStorage.removeItem('token'); // Remove the token
    navigate('/'); // Redirect to login page
};

const endpoint = "http://localhost:8080";

const handleSubmit = async (e) => {
e.preventDefault();
const formData = new FormData();
formData.append('image', product.image);
formData.append('title', product.title);
formData.append('cost', product.cost);
formData.append('location', product.location);
formData.append('description', product.description);

try {
    const response = await fetch(`${endpoint}/create-product`, {
    method: 'POST',
    body: formData,
    headers: {
      // 'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
    });
    console.log(response)
    if (!response.ok) {
      if (response.status === 401) {
        Swal.fire({
          title: 'Session Expired',
          text: 'Please login again to continue.',
          icon: 'warning',
          confirmButtonText: 'OK'
        }).then((result) => {
          if (result.isConfirmed) {
            localStorage.removeItem('token'); // Optional: clear token
            navigate('/'); // Redirect to login when OK is clicked
          }
        });
      
      }
      else {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    }
    Swal.fire({
      title: 'Success!',
      text: 'Successfully added the product for sale',
      icon: 'success',
      confirmButtonText: 'OK'
  }).then((result) => {
      if (result.isConfirmed) {
        navigate('/home');// Redirect to login when OK is clicked
      }
  });

    // Assuming you want to do something upon successfully adding a product
     // Redirect to products page
} catch (error) {
    console.error('There was an error selling the product:', error);
}
};

return (
     <div className="sell-product-container">
       <header className="sell-header">
      <h2>Sell Your Product</h2>
      <div className="top-nav">
        <Link to="/home" className="nav-link">Products</Link>
        <button onClick={handleLogout} className="logout-button">Logout</button>
      </div>
    </header>
      <form onSubmit={handleSubmit} className="sell-form">
        <input
          type="text"
          name="title"
          placeholder="Title"
          value={product.title}
          onChange={handleChange}
          required
        />
         <input
    type="text"
    name="cost"
    placeholder="Cost"
    value={product.cost}
    onChange={handleChange}
    required
    />
    <input
    type="text"
    name="location"
    placeholder="Location"
    value={product.location}
    onChange={handleChange}
    required
    />
    <textarea
    name="description"
    placeholder="Description"
    value={product.description}
    onChange={handleChange}
    required
    />
    <input
    type="file"
    name="image"
    onChange={handleImageChange}
    required
    />
    <button type="submit">Add Product</button>
</form>
</div>
  );
};

export default SellProduct;
