import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import './ProductDetails.css'; // Make sure to create this CSS file for styling

const endpoint = "http://localhost:8080";

function ProductDetails() {
  const [productDetails, setProductDetails] = useState(null);
  const { id } = useParams();

  useEffect(() => {
    const fetchProductDetails = async () => {
      try {
        const response = await fetch(`${endpoint}/product/${id}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        setProductDetails(data);
      } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
      }
    };

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
      </div>
    </div>
  </div>
  );
}

export default ProductDetails;
