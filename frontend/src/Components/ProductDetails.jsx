import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import './ProductDetails.css';
import { useNavigate, Link } from 'react-router-dom';
import Swal from 'sweetalert2';
import { apiUrl, resolveImageUrl } from '../config';

function ProductDetails() {
  const [productDetails, setProductDetails] = useState(null);
  const { id } = useParams();
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const handleBuyClick = (productId) => {
    // Placeholder for future buy functionality
    navigate(`/checkout/${productId}`);
    console.log(`Buy button clicked for product ID: ${productId}`);
    // Here you would typically make an API call to your buy endpoint, passing the productId
  };

  useEffect(() => {
    const fetchProductDetails = async () => {
      try {
        const response = await fetch(apiUrl(`/product/${id}`), {
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
  }, [id, navigate, token]); // Only re-run the effect if the ID changes

  if (!productDetails) {
    return <div>Loading...</div>; // Loading state
  }

  // Render the product details
  return (
    <div className="product-details-page">
      <div className="detail-topbar">
        <Link to="/home" className="nav-link">Back to marketplace</Link>
      </div>
      <div className="product-content">
        <div className="product-visual">
          <img src={resolveImageUrl(productDetails.image)} alt={productDetails.title} className="product-image" />
        </div>
        <div className="product-info">
          <span className={`detail-status-pill ${productDetails.status === 'sold' ? 'sold' : 'available'}`}>
            {productDetails.status || 'available'}
          </span>
          <h1 className="product-title">{productDetails.title}</h1>
          <p className="product-cost">${productDetails.cost}</p>
          <p className="product-location">Pickup in {productDetails.location}</p>
          <p className="product-description">{productDetails.description}</p>
          <div className="detail-meta">
            <span>Seller: {productDetails.seller || 'Book Mart member'}</span>
            {productDetails.buyer && <span>Purchased by: {productDetails.buyer}</span>}
          </div>
          <button
            className="buy-button"
            onClick={() => handleBuyClick(productDetails.productId)}
            disabled={productDetails.status && productDetails.status !== 'available'}
          >
            {productDetails.status === 'sold' ? 'Already sold' : 'Continue to checkout'}
          </button>
        </div>
      </div>
    </div>
  );
}

export default ProductDetails;
