import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom'; 
import './Products.css'; // Ensure this file exists and has the content from below

const endpoint = "http://localhost:8080";

function Products() {
  const [products, setProducts] = useState([]);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await fetch(endpoint + '/products', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        setProducts(data);
      } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
      }
    };

    fetchProducts();
  }, []);

  const handleSellClick = () => {
    console.log('Sell button clicked');
    // Placeholder for future functionality
  };

  const handleBuyClick = (productId) => {
    // Placeholder for future buy functionality
    console.log(`Buy button clicked for product ID: ${productId}`);
    // Here you would typically make an API call to your buy endpoint, passing the productId
  };

//   return (
//     <div className="products-page">
//       <div className="products-header">
//         <h1>Products</h1>
//         <button className="sell-button" onClick={handleSellClick}>Sell</button>
//       </div>
//       <div className="products-container">
//         {products.map((product) => (
//           <div key={product.productId} className="product-card">
//             <img src={product.image} alt={product.title} className="product-image" />
//             <div className="product-details">
//               <h3 className="product-name">{product.title}</h3>
//               <p className="product-cost">${product.cost}</p>
//               <p className="product-location">{product.location}</p>
//               <p className="product-description">{product.description}</p>
//               <button className="buy-button" onClick={() => handleBuyClick(product.productId)}>Buy</button>
//             </div>
//           </div>
//         ))}
//       </div>
//     </div>
//   );
// }

// export default Products;

return (
    <div className="products-page">
    <div className="products-header">
    <h1>Products</h1>
    <button className="sell-button" onClick={handleSellClick}>Sell</button>
    </div>
    <div className="products-container">
      {products.map((product) => (
        <Link to={`/product/${product.productId}`} key={product.productId} className="product-card">
          <img src={product.image} alt={product.title} className="product-image" />
          <div className="product-details">
            <h3 className="product-name">{product.title}</h3>
            <p className="product-cost">${product.cost}</p>
            <p className="product-location">{product.location}</p>
            <p className="product-description">{product.description}</p>
            <button className="buy-button" onClick={() => handleBuyClick(product.productId)}>Buy</button>
          </div>
        </Link>
      ))}
    </div>
  </div>
);
}

export default Products;