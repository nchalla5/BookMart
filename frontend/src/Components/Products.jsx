import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Swal from 'sweetalert2';
import './Products.css';
import  { jwtDecode } from 'jwt-decode'; // Ensure correct import, it might be `jwt-decode` without braces based on the package export type

const endpoint = "http://localhost:8080";

function Products() {
  const [products, setProducts] = useState([]);
  const [searchType, setSearchType] = useState('name');
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/login');
      return;
    }

    const decodedToken = jwtDecode(token);
    console.log(decodedToken);


    const fetchProducts = async () => {
      let queryString = `/products?`;

      // Dynamically add search query based on the searchType
      if (searchTerm) {
        switch (searchType) {
          case 'name':
            queryString += `searchName=${searchTerm}&`;
            break;
          case 'location':
            queryString += `searchLocation=${searchTerm}&`;
            break;
          case 'status':
            queryString += `statusFilter=${searchTerm}&`;
            break;
          default:
            break;
        }
      }

      if (sortField && sortOrder) {
        queryString += `sortField=${sortField}&sortOrder=${sortOrder}`;
      }

      try {
        const response = await fetch(endpoint + queryString, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        setProducts(data);
      } catch (error) {
        console.error("There was an error fetching the products", error);
        Swal.fire("Error", "Failed to fetch products. Please try again later.", "error");
      }
    };

    fetchProducts();
  }, [navigate, searchTerm, sortField, sortOrder, searchType]);

  const handleSearchChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleSearchTypeChange = (e) => {
    setSearchType(e.target.value);
    setSearchTerm(''); // Optionally reset searchTerm if you want the input to clear upon changing type
  };

  const handleSortChange = (e) => {
    const [field, order] = e.target.value.split('-');
    setSortField(field);
    setSortOrder(order);
  };

  const handleSellClick = () => {
    navigate('/sell');
  };

  const handleBuyClick = (productId) => {
    navigate(`/product/${productId}`);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <div className="products-page">
    <div className="products-header">
      <h1>Books</h1>
      <div className="search-sort-bar">
        <div className="search-bar">
          <select value={searchType} onChange={handleSearchTypeChange}>
            <option value="name">Search by Name</option>
            <option value="location">Filter by Location</option>
            <option value="status">Filter by Status</option>
          </select>
          <input
            type="text" className='searchTerm'
            placeholder={`Search by ${searchType}...`}
            value={searchTerm}
            onChange={handleSearchChange}
          />
        </div>
        <div className="sort-bar">
          <select onChange={handleSortChange}>
            <option value="">Default Sort</option>
            <option value="cost-asc">Cost Low to High</option>
            <option value="cost-desc">Cost High to Low</option>
            <option value="name-asc">Name A-Z</option>
            <option value="name-desc">Name Z-A</option>
          </select>
        </div>
        <button className="sell-button" onClick={handleSellClick}>Sell</button>
        <button className="logout-button" onClick={handleLogout}>Logout</button>
      </div>
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
