import React, { useState, useEffect } from 'react';
import './ProductCheckout.css';
import { useNavigate, useParams } from 'react-router-dom';
import Swal from 'sweetalert2';


const endpoint = "http://localhost:8080";


function ProductCheckout() {
    const navigate = useNavigate();
    const [product, setProduct] = useState(null);
    const { id } = useParams();
    const token = localStorage.getItem('token');
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
            console.log(data);
            setProduct(data);
          } catch (error) {
            console.error('There was a problem with the fetch operation:', error);
          }
        };
        if (!token) {
          navigate('/'); // Redirect to login if token not found
          return;
        }
        fetchProductDetails();
      }, [id]);
    if (!product) {
    return <div>Loading...</div>; // Loading state
    }
    const handleCheckoutClick = (productId) => {
        console.log(productId);
    };

    return (
    <div className="check-page">
        <div className="check-container">
        <div className="check-form">
            <h3>Shipping Address: </h3>
            <input type="text" name="street" placeholder="Street address" required />
            <input type="text" name="city" placeholder="City" required />
            <input type="text" name="state" placeholder="State" required />
            <input type="number" name="postalCode" placeholder="Postal Code" required />
            <input type="text" name="country" placeholder="Country" required />
            <input type="text" name="countryCode" placeholder="Country Code" required />
            <input type="number" name="mobileNumber" placeholder="Mobile Number" required />
            <button onClick={() => handleCheckoutClick(product.productId)}>Checkout</button>
        </div>
        <div className="check-card"><br></br>
          <img src={product.image} alt={product.title} className="check-image" />
          <div className="check-details">
            <h3 className="check-name">{product.title}</h3>
            <p className="check-cost">${product.cost}</p>
            <p className="check-location">{product.location}</p>
            <p className="check-description">{product.description}</p>
          </div>
        </div>
        </div>
    </div>
    )
}

export default ProductCheckout;