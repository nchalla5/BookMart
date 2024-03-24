import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import './ProductDetails.css';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';
import Swal from 'sweetalert2';

const endpoint = "http://localhost:8080";

const handleBuyClick = (productId) => {
    // Placeholder for future buy functionality
    console.log(`Buy button clicked for product ID: ${productId}`);
    // Here you would typically make an API call to your buy endpoint, passing the productId
  };

function ProductDetails() {
  const [productDetails, setProductDetails] = useState(null);
  const { id } = useParams();
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  

  useEffect(() => {
    const fetchProductDetails = async () => {
      try {
        const response = await fetch(`${endpoint}/product/${id}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          }
        });

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

        const data = await response.json();
        setProductDetails(data);
      } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
      }
    };
    if (!token) {
      navigate('/'); // Redirect to login if token not found
      return;
    }

    fetchProductDetails();
  }, [id]); // Only re-run the effect if the ID changes

  if (!productDetails) {
    return <div>Loading...</div>; // Loading state
  }

  // Render the product details
  return (
    <div className="product-details-page">
    <h1 className="product-title">{productDetails.title}</h1> {/* Title at the top */}
    <div className="product-content"> {/* Flex container for image and details */}
      <img src={productDetails.image} alt={productDetails.title} className="product-image" />
      <div className="product-info">
        <p className="product-cost">${productDetails.cost}</p>
        <p className="product-location">{productDetails.location}</p>
        <p className="product-description">{productDetails.description}</p>
        {/* Include Buy button or other details as needed */}
        <button className="buy-button" onClick={() => handleBuyClick(productDetails.productId)}>Buy</button>
      </div>
    </div>
  </div>
  );
}

export default ProductDetails;
