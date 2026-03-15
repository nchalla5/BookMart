import React, { useState, useEffect } from 'react';
import './ProductCheckout.css';
import { useNavigate, useParams } from 'react-router-dom';
import Swal from 'sweetalert2';
import { apiUrl, resolveImageUrl } from '../config';

function ProductCheckout() {
    const navigate = useNavigate();
    const [product, setProduct] = useState(null);
    const [formValues, setFormValues] = useState({
      street: '',
      city: '',
      state: '',
      postalCode: '',
      country: '',
      countryCode: '',
      mobileNumber: '',
    });
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { id } = useParams();
    const token = localStorage.getItem('token');
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
      }, [id, navigate, token]);
    if (!product) {
    return <div>Loading...</div>; // Loading state
    }

    const handleChange = (event) => {
      const { name, value } = event.target;
      setFormValues((current) => ({
        ...current,
        [name]: value,
      }));
    };

    const handleCheckoutClick = async (event) => {
      event.preventDefault();
      setIsSubmitting(true);

      try {
        const response = await fetch(apiUrl(`/product/${product.productId}/purchase`), {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({
            shippingAddress: formValues,
          }),
        });

        if (!response.ok) {
          if (response.status === 409) {
            Swal.fire('Unavailable', 'This book has already been purchased.', 'warning');
            navigate('/home');
            return;
          }

          if (response.status === 401) {
            localStorage.removeItem('token');
            navigate('/');
            return;
          }

          throw new Error(`HTTP error! status: ${response.status}`);
        }

        await response.json();
        Swal.fire('Success', 'Purchase completed successfully.', 'success');
        navigate('/home');
      } catch (error) {
        console.error('There was a problem with checkout:', error);
        Swal.fire('Error', 'Unable to complete checkout right now.', 'error');
      } finally {
        setIsSubmitting(false);
      }
    };

    return (
    <div className="check-page">
        <div className="check-container">
        <div className="check-summary">
          <p className="check-eyebrow">Checkout</p>
          <h2>Finish your order</h2>
          <p className="check-copy">Confirm where the book should go and we&apos;ll mark it as sold in the marketplace.</p>
        </div>
        <form className="check-form" onSubmit={handleCheckoutClick}>
            <h3>Shipping Address: </h3>
            <input type="text" name="street" placeholder="Street address" value={formValues.street} onChange={handleChange} required />
            <input type="text" name="city" placeholder="City" value={formValues.city} onChange={handleChange} required />
            <input type="text" name="state" placeholder="State" value={formValues.state} onChange={handleChange} required />
            <input type="text" name="postalCode" placeholder="Postal Code" value={formValues.postalCode} onChange={handleChange} required />
            <input type="text" name="country" placeholder="Country" value={formValues.country} onChange={handleChange} required />
            <input type="text" name="countryCode" placeholder="Country Code" value={formValues.countryCode} onChange={handleChange} required />
            <input type="text" name="mobileNumber" placeholder="Mobile Number" value={formValues.mobileNumber} onChange={handleChange} required />
            <button type="submit" disabled={isSubmitting || (product.status && product.status !== 'available')}>
              {product.status === 'sold' ? 'Already Sold' : (isSubmitting ? 'Processing...' : 'Checkout')}
            </button>
        </form>
        <div className="check-card">
          <img src={resolveImageUrl(product.image)} alt={product.title} className="check-image" />
          <div className="check-details">
            <h3 className="check-name">{product.title}</h3>
            <p className="check-cost">${product.cost}</p>
            <p className="check-location">{product.location}</p>
            <p className="check-description">{product.description}</p>
            <p className="check-status">Status: {product.status || 'available'}</p>
          </div>
        </div>
        </div>
    </div>
    )
}

export default ProductCheckout;
