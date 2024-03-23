
// import Home from './Home';
import BuyItem from './BuyItem';
import NotFound from './NotFound';
import Products from './Products';
import ProductDetails from './ProductDetails';
import SellProduct from './SellProduct';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Login from './Login';

function Body() {
    return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />}></Route>
        <Route path="/home" element={<Products />}></Route>
        <Route path="/product/:id" element={<ProductDetails />} />
        <Route path="/buy" element={<BuyItem />}></Route>
        <Route path="/sell" element={<SellProduct />}></Route>
        <Route path="*" element={<NotFound />}></Route>
      </Routes>
    </Router>
    );
}
export default Body