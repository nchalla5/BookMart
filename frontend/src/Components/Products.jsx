import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Swal from 'sweetalert2';
import './Products.css';
import  { jwtDecode } from 'jwt-decode'; // Ensure correct import, it might be `jwt-decode` without braces based on the package export type
import { apiUrl, resolveImageUrl } from '../config';

function Products() {
  const [products, setProducts] = useState([]);
  const [searchType, setSearchType] = useState('name');
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [viewerName, setViewerName] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/login');
      return;
    }

    const decodedToken = jwtDecode(token);
    setViewerName(decodedToken.name || 'Reader');


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
        const response = await fetch(apiUrl(queryString), {
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
      <section className="marketplace-hero">
        <div className="marketplace-copy">
          <p className="marketplace-eyebrow">Curated marketplace</p>
          <h1>Book Mart</h1>
          <p className="marketplace-subcopy">
            Welcome back, {viewerName}. Browse affordable second-hand books from readers near you.
          </p>
        </div>
        <div className="marketplace-side">
          <div className="marketplace-stats">
            <div className="stat-card">
              <span>{products.length}</span>
              <p>Books visible</p>
            </div>
            <div className="stat-card">
              <span>{products.filter((product) => product.status !== 'sold').length}</span>
              <p>Available now</p>
            </div>
          </div>
          <div className="action-bar">
            <button className="sell-button" onClick={handleSellClick}>Sell a Book</button>
            <button className="logout-button" onClick={handleLogout}>Logout</button>
          </div>
        </div>
      </section>

      <div className="products-toolbar">
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
        </div>
      </div>
      <div className="products-container">
        {products.length === 0 && <p className="empty-state">No books match the current filters yet.</p>}
        {products.map((product) => (
          <Link to={`/product/${product.productId}`} key={product.productId} className="product-card">
            <img src={resolveImageUrl(product.image)} alt={product.title} className="product-image" />
            <div className="product-details">
              <div className="product-meta-row">
                <span className={`status-pill ${product.status === 'sold' ? 'sold' : 'available'}`}>
                  {product.status || 'available'}
                </span>
                <span className="seller-name">{product.location}</span>
              </div>
              <h3 className="product-name">{product.title}</h3>
              <p className="product-cost">${product.cost}</p>
              <p className="product-description">{product.description}</p>
              <button
                className="buy-button"
                type="button"
                disabled={product.status && product.status !== 'available'}
                onClick={(event) => {
                  event.preventDefault();
                  handleBuyClick(product.productId);
                }}
              >
                {product.status === 'sold' ? 'Sold' : 'Buy'}
              </button>
              </div>
        </Link>
      ))}
    </div>
  </div>
);
}


export default Products;
