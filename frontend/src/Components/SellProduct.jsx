import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const SellProduct = () => {
const [product, setProduct] = useState({
title: '',
cost: '',
location: '',
description: '',
image: null,
});
const navigate = useNavigate();

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
    });
    console.log(response)
    if (!response.ok) {
    throw new Error('Failed to sell product');
    }

    // Assuming you want to do something upon successfully adding a product
    navigate('/home'); // Redirect to products page
} catch (error) {
    console.error('There was an error selling the product:', error);
}
};

return (
<form onSubmit={handleSubmit}>
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
);
};

export default SellProduct;
